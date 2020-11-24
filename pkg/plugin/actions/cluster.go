package actions

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"k8s.io/client-go/tools/clientcmd"
	_ "k8s.io/client-go/tools/clientcmd"
	"os"
	"sigs.k8s.io/kind/pkg/apis/config/v1alpha4"
	"sigs.k8s.io/kind/pkg/cluster"
	"strconv"
)

const (
	CreateKindClusterAction = "octant-plugin-for-kind/create"
	DeleteKindClusterAction = "octant-plugin-for-kind/delete"
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

		var features map[string]interface{}
		featureGatesData, err := request.Payload.Raw("featureGates")
		if err != nil {
			return err
		}
		if err := json.Unmarshal(featureGatesData, &features); err != nil {
			return err
		}

		formData := ClusterConfig{
			Details:  details,
			Features: features,
		}
		return createCluster(request, formData, provider)
	case DeleteKindClusterAction:
		return deleteCluster(request, provider)
	default:
		return fmt.Errorf("unable to find handler for plugin: %s", "kind")
	}
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

	// TODO: Fix upstream bug where NewFormFieldNumber default returns a string
	numCPNodes := int(clusterConfig.Details.ControlPlaneNodes)
	numWorkers := int(clusterConfig.Details.Workers)

	if numCPNodes >= 0 {
		var version string
		if len(clusterConfig.Details.Version) > 0 {
			version = clusterConfig.Details.Version[0]
		}

		for i := 0; i < numCPNodes; i++ {
			node := v1alpha4.Node{
				Role:  v1alpha4.ControlPlaneRole,
				Image: version,
			}
			nodes = append(nodes, node)
		}
	}

	if numWorkers >= 0 {
		var version string
		if len(clusterConfig.Details.Version) > 0 {
			version = clusterConfig.Details.Version[0]
		}
		for i := 0; i < numWorkers; i++ {
			worker := v1alpha4.Node{
				Role:  v1alpha4.WorkerRole,
				Image: version,
			}
			nodes = append(nodes, worker)
		}
	}

	featureGates := make(map[string]bool)

	for key, value := range clusterConfig.Features {
		if value != nil {
			featureGates[key] = value.(bool)
		}
	}

	kindCluster := &v1alpha4.Cluster{
		Nodes:        nodes,
		FeatureGates: featureGates,
	}

	alert := action.CreateAlert(action.AlertTypeInfo, "Creating cluster: " + clusterName, action.DefaultAlertExpiration)
	request.DashboardClient.SendAlert(request.Context(), request.ClientID, alert)

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

	alert := action.CreateAlert(action.AlertTypeInfo, "Deleted kind cluster: " + clusterName, action.DefaultAlertExpiration)
	request.DashboardClient.SendAlert(request.Context(), request.ClientID, alert)
	return nil
}

// ClusterConfig contains input from stepper
type ClusterConfig struct {
	Details  ClusterDetails         `json:"clusterConfiguration"`
	Features map[string]interface{} `json:"featureGates"`
}

// ClusterDetails are used to build v1alpha4Config
type ClusterDetails struct {
	ClusterName       string   `json:"clusterName"`
	ControlPlaneNodes FlexInt  `json:"controlPlaneNodes"`
	Workers           FlexInt  `json:"workers"`
	Version           []string `json:"version"`
}

type FlexInt int

func (fi *FlexInt) UnmarshalJSON(b []byte) error {
	if b[0] != '"' {
		return json.Unmarshal(b, (*int)(fi))
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*fi = FlexInt(i)
	return nil
}
