package actions

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/vmware-tanzu/octant-plugin-for-kind/pkg/docker"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"sigs.k8s.io/kind/pkg/apis/config/v1alpha4"
	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cluster/nodes"
	"sigs.k8s.io/kind/pkg/cluster/nodeutils"
	"sigs.k8s.io/kind/pkg/fs"
)

const (
	// CreateKindClusterAction is the action name for creating a cluster
	CreateKindClusterAction = "octant-plugin-for-kind.dev/create"
	// DeleteKindClusterAction is the action name for deleting a cluster
	DeleteKindClusterAction = "octant-plugin-for-kind.dev/delete"
	// LoadImageAction is the action name for loading a kind image
	LoadImageAction = "octant-plugin-for-kind.dev/loadImage"
	// DeleteImageAction is the action name for deleting a kind image
	DeleteImageAction = "octant-plugin-for-kind.dev/deleteImage"
)

// ActionHandler is a handler for actions
func ActionHandler(request *service.ActionRequest) error {
	actionName, err := request.Payload.String("action")
	if err != nil {
		return err
	}

	provider := cluster.NewProvider()

	switch actionName {
	case CreateKindClusterAction:
		var details ClusterDetails
		clusterConfigData, err := request.Payload.Raw("clusterConfiguration")
		if err != nil {
			return err
		}
		if err := json.Unmarshal(clusterConfigData, &details); err != nil {
			return err
		}

		var features map[string][]string
		featureGatesData, err := request.Payload.Raw("featureGates")
		if err != nil {
			return err
		}
		if err := json.Unmarshal(featureGatesData, &features); err != nil {
			return err
		}

		var networking NetworkingDetails
		networkingConfigData, err := request.Payload.Raw("networking")
		if err != nil {
			return err
		}
		if err := json.Unmarshal(networkingConfigData, &networking); err != nil {
			return err
		}

		formData := ClusterConfig{
			Details:    details,
			Features:   features["__featureGate"],
			Networking: convertTov1alpha4Networking(networking),
		}
		return createCluster(request, formData, provider)
	case DeleteKindClusterAction:
		return deleteCluster(request, provider)
	case LoadImageAction:
		imageName, err := request.Payload.String("imageName")
		if err != nil {
			return err
		}
		clusterName, err := request.Payload.String("clusterName")
		if err != nil {
			return err
		}
		return loadImage(request, provider, clusterName, imageName)
	case DeleteImageAction:
		imageID, err := request.Payload.String("imageID")
		if err != nil {
			return err
		}
		clusterName, err := request.Payload.String("clusterName")
		if err != nil {
			return err
		}
		return deleteImage(request, clusterName, imageID)
	default:
		return fmt.Errorf("unable to find handler for plugin: %s", "kind")
	}
}

func convertTov1alpha4Networking(networking NetworkingDetails) *v1alpha4.Networking {
	var kindNetworking v1alpha4.Networking
	if networking.IPFamily != "" {
		kindNetworking.IPFamily = v1alpha4.ClusterIPFamily(networking.IPFamily)
	}
	if networking.APIServerAddress != "" {
		kindNetworking.APIServerAddress = networking.APIServerAddress
	}
	if networking.APIServerPort != 0 {
		kindNetworking.APIServerPort = networking.APIServerPort
	}
	if networking.PodSubnet != "" {
		kindNetworking.PodSubnet = networking.PodSubnet
	}
	if networking.ServiceSubnet != "" {
		kindNetworking.ServiceSubnet = networking.ServiceSubnet
	}
	if len(networking.DisableDefaultCNI) == 1 && networking.DisableDefaultCNI[0] == "true" {
		kindNetworking.DisableDefaultCNI = true
	}
	return &kindNetworking
}

func createCluster(request *service.ActionRequest, clusterConfig ClusterConfig, provider *cluster.Provider) error {
	var nodes []v1alpha4.Node

	// Check if cluster name already exists
	n, err := provider.ListNodes("cluster-name")
	if err != nil {
		return err
	}

	if len(n) != 0 {
		return errors.Errorf("Name already in use %q", "cluster-name")
	}

	clusterName := clusterConfig.Details.ClusterName
	if clusterName == "" {
		return errors.Errorf("Cluster name cannot be empty")
	}

	if clusterConfig.Details.ControlPlaneNodes >= 0 {
		var version string
		if len(clusterConfig.Details.Version) > 0 {
			version = clusterConfig.Details.Version[0]
		}

		for i := 0; i < clusterConfig.Details.ControlPlaneNodes; i++ {
			node := v1alpha4.Node{
				Role:  v1alpha4.ControlPlaneRole,
				Image: version,
			}
			nodes = append(nodes, node)
		}
	}

	if clusterConfig.Details.Workers >= 0 {
		var version string
		if len(clusterConfig.Details.Version) > 0 {
			version = clusterConfig.Details.Version[0]
		}
		for i := 0; i < clusterConfig.Details.Workers; i++ {
			worker := v1alpha4.Node{
				Role:  v1alpha4.WorkerRole,
				Image: version,
			}
			nodes = append(nodes, worker)
		}
	}

	featureGates := make(map[string]bool)

	for _, feature := range clusterConfig.Features {
		featureGates[feature] = true
	}

	kindCluster := &v1alpha4.Cluster{
		Nodes:        nodes,
		FeatureGates: featureGates,
		Networking:   *clusterConfig.Networking,
	}
	alert := action.CreateAlert(action.AlertTypeInfo, "Creating cluster: "+clusterName, action.DefaultAlertExpiration)
	request.DashboardClient.SendAlert(request.Context(), request.ClientState.ClientID, alert)

	// TODO: Show status when creating cluster
	// TODO: Kind switches to new context once ready. Bad UX?
	if err := provider.Create(clusterName, cluster.CreateWithV1Alpha4Config(kindCluster)); err != nil {
		return err
	}

	return nil
}

func deleteCluster(request *service.ActionRequest, provider *cluster.Provider) error {
	payload := request.Payload

	var kubeConfigPath string
	kubeConfig := os.Getenv("KUBECONFIG")

	if kubeConfig == "" {
		kubeConfigPath = clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
	}

	clusterName, err := payload.String("name")
	if err != nil {
		return err
	}

	if err := provider.Delete(clusterName, kubeConfigPath); err != nil {
		return err
	}

	alert := action.CreateAlert(action.AlertTypeInfo, "Deleted kind cluster: "+clusterName, action.DefaultAlertExpiration)
	request.DashboardClient.SendAlert(request.Context(), request.ClientState.ClientID, alert)
	return nil
}

// ClusterConfig contains input from stepper
type ClusterConfig struct {
	Details    ClusterDetails       `json:"clusterConfiguration"`
	Features   []string             `json:"featureGates"`
	Networking *v1alpha4.Networking `json:"networking"`
}

// ClusterDetails are used to build v1alpha4Config
type ClusterDetails struct {
	ClusterName       string   `json:"clusterName"`
	ControlPlaneNodes int      `json:"controlPlaneNodes"`
	Workers           int      `json:"workers"`
	Version           []string `json:"version"`
}

type NetworkingDetails struct {
	IPFamily          string   `json:"ipFamily,omitempty"`
	APIServerPort     int32    `json:"apiServerPort,omitempty"`
	APIServerAddress  string   `json:"apiServerAddress,omitempty"`
	PodSubnet         string   `json:"podSubnet,omitempty"`
	ServiceSubnet     string   `json:"serviceSubnet,omitempty"`
	DisableDefaultCNI []string `json:"disableDefaultCNI,omitempty"`
}

func loadImage(request *service.ActionRequest, provider *cluster.Provider, clusterName, imageName string) error {
	nodeList, err := provider.ListInternalNodes(clusterName)
	if err != nil {
		return err
	}
	if len(nodeList) == 0 {
		return fmt.Errorf("no nodes for cluster %q", clusterName)
	}

	var selectedNodes []nodes.Node
	for _, node := range nodeList {
		id, err := nodeutils.ImageID(node, imageName)
		if err != nil || imageName != id {
			selectedNodes = append(selectedNodes, node)
		}
	}

	if len(selectedNodes) == 0 {
		return nil
	}

	dir, err := fs.TempDir("", "image-tar")
	if err != nil {
		return fmt.Errorf("failed to create tempdir: %+v", err)
	}
	defer os.RemoveAll(dir)
	imageTarPath := filepath.Join(dir, "image.tar")

	client := docker.NewDockerClient()
	if err := client.Save(request.Context(), imageTarPath, imageName); err != nil {
		return err
	}

	f, err := os.Open(imageTarPath)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, selectedNode := range selectedNodes {
		nodeutils.LoadImageArchive(selectedNode, f)
	}

	alert := action.CreateAlert(action.AlertTypeInfo, "Loading image: "+imageName, action.DefaultAlertExpiration)
	request.DashboardClient.SendAlert(request.Context(), request.ClientState.ClientID, alert)
	return nil
}

func deleteImage(request *service.ActionRequest, clusterName string, imageID string) error {
	client := docker.NewDockerClient()

	if err := client.DeleteKindImage(request.Context(), clusterName, imageID); err != nil {
		alert := action.CreateAlert(action.AlertTypeError, "Failed to delete kind image: "+err.Error(), action.DefaultAlertExpiration)
		request.DashboardClient.SendAlert(request.Context(), request.ClientState.ClientID, alert)
		return err
	}
	alert := action.CreateAlert(action.AlertTypeInfo, "Deleted kind image: "+imageID, action.DefaultAlertExpiration)
	request.DashboardClient.SendAlert(request.Context(), request.ClientState.ClientID, alert)
	return nil
}
