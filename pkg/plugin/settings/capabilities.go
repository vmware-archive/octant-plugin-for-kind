package settings

import (
	"github.com/vmware-tanzu/octant-plugin-for-kind/pkg/plugin/actions"
	"github.com/vmware-tanzu/octant/pkg/plugin"
)

func GetCapabilities() *plugin.Capabilities {
	return &plugin.Capabilities{
		ActionNames: []string{actions.CreateKindClusterAction, actions.DeleteKindClusterAction},
		IsModule:    true,
	}
}
