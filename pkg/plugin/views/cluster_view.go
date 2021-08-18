package views

import (
	"github.com/vmware-tanzu/octant-plugin-for-kind/pkg/docker"
	"sigs.k8s.io/kind/pkg/apis/config/v1alpha4"
	"strconv"
	"time"

	"github.com/vmware-tanzu/octant-plugin-for-kind/pkg/plugin/actions"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"

	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cluster/nodeutils"
)

// BuildKindClusterView builds the layout of kind
func BuildKindClusterView(request service.Request) (component.Component, error) {
	ctx := request.Context()
	client := docker.NewDockerClient()

	table := component.NewTableWithRows(
		"Clusters",
		"There are no kind clusters",
		component.NewTableCols("Name", "Status", "Nodes", "Version", "Age"),
		[]component.TableRow{})

	provider := cluster.NewProvider()
	clusterNames, err := provider.List()
	if err != nil {
		return nil, err
	}

	for _, name := range clusterNames {
		container, err := client.KindControlPlaneContainer(ctx, name)
		if err != nil {
			return nil, err
		}

		nodes, err := provider.ListNodes(name)
		if err != nil {
			return nil, err
		}

		// Exclude external load balancer to match `kind get nodes`
		internalNodes, err := nodeutils.InternalNodes(nodes)
		if err != nil {
			return nil, err
		}

		tableRow := component.TableRow{
			"Name":    component.NewLink(name, name, name),
			"Status":  component.NewText(container.State),
			"Nodes":   component.NewText(strconv.Itoa(len(internalNodes))),
			"Version": component.NewText(container.Version),
			"Age":     component.NewTimestamp(time.Unix(container.Created, 0)),
		}

		tableRow.AddAction(component.GridAction{
			Name:       "Delete",
			ActionPath: actions.DeleteKindClusterAction,
			Payload: action.Payload{
				"name": name,
			},
			Confirmation: &component.Confirmation{
				Title: "Delete kind cluster - " + name,
				Body:  "Are you sure?",
			},
			Type: "",
		})

		table.Add(tableRow)
	}

	clusterNameFormField := component.NewFormFieldText("Cluster Name", "clusterName", "")
	clusterNameFormField.AddValidator("This cannot be empty", map[component.FormValidator]interface{}{
		component.FormValidatorRequired: true,
	})
	clusterConfigurationForm := component.Form{
		Fields: []component.FormField{
			clusterNameFormField,
			component.NewFormFieldNumber("Control Plane Nodes", "controlPlaneNodes", "1"),
			component.NewFormFieldNumber("Workers", "workers", "0"),
			component.NewFormFieldSelect(
				"Version",
				"version",
				generateInputChoices(NewImageMap()),
				true),
		},
	}

	featureGates, err := getFeatureGateList()
	if err != nil {
		return nil, err
	}
	var choices []component.InputChoice
	for _, fg := range Unique(featureGates) {
		choices = append(choices, component.InputChoice{
			Label:   fg.Feature,
			Value:   fg.Feature,
			Checked: false,
		})
	}
	featureGatesForm := component.Form{
		Fields: []component.FormField{
			component.NewFormFieldCheckBox("", "__featureGate", choices),
		},
	}

	networkingForm := component.Form{
		Fields: []component.FormField{
			component.NewFormFieldRadio("IP Family", "ipFamily", []component.InputChoice{
				{
					"IPv4",
					string(v1alpha4.IPv4Family),
					false,
				},
				{
					"IPv6",
					string(v1alpha4.IPv6Family),
					false,
				},
			}),
			component.NewFormFieldText("API Server Address", "apiServerAddress", ""),
			component.NewFormFieldNumber("API Server Port", "apiServerPort", ""),
			component.NewFormFieldCheckBox("Disable Default CNI", "disableDefaultCNI", []component.InputChoice{
				{
					"Yes",
					"true",
					false,
				},
				{
					"No",
					"false",
					false,
				},
			}),
			component.NewFormFieldText("Pod Subnet", "podSubnet", ""),
			component.NewFormFieldText("Service Subnet", "serviceSubnet", ""),
			component.NewFormFieldRadio("", "", []component.InputChoice{
				{
					"iptables",
					"iptables",
					false,
				},
				{
					"ipvs",
					"ipvs",
					false,
				},
			}),
		},
	}

	stepper := component.Stepper{
		Base: component.Base{},
		Config: component.StepperConfig{
			Action: actions.CreateKindClusterAction,
			Steps: []component.StepConfig{
				{
					Name:        "clusterConfiguration",
					Form:        clusterConfigurationForm,
					Title:       "Cluster configuration",
					Description: "Build a cluster config for kind.x-k8s.io/v1alpha4",
				},
				{
					Name:        "networking",
					Form:        networkingForm,
					Title:       "Networking",
					Description: "Customize networking",
				},
				{
					Name:        "featureGates",
					Form:        featureGatesForm,
					Title:       "Feature Gates",
					Description: "Select feature gates to be enabled",
				},
			},
		},
	}

	modal := component.NewModal(component.TitleFromString("New Kind Cluster"))
	modal.SetBody(&stepper)
	modal.SetSize(component.ModalSizeExtraLarge)

	button := component.NewButton("Create a cluster", action.Payload{}, component.WithModal(modal))
	table.Config.ButtonGroup = component.NewButtonGroup()
	table.Config.ButtonGroup.AddButton(button)

	flexLayout := component.NewFlexLayout("")
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: table},
	})

	return flexLayout, nil
}

func generateInputChoices(orderedMap *OrderedMap) []component.InputChoice {
	var inputChoices []component.InputChoice
	for index, key := range orderedMap.Keys() {
		choice := component.InputChoice{
			Label: key,
			Value: orderedMap.Map()[key],
		}
		if index == 0 {
			choice.Checked = true
		}
		inputChoices = append(inputChoices, choice)
	}
	return inputChoices
}
