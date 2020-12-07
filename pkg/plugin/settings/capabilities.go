package settings

import (
	"github.com/vmware-tanzu/octant-plugin-for-kind/pkg/plugin/actions"
	"github.com/vmware-tanzu/octant/pkg/plugin"
)

// GetCapabilities returns the list of plugin capabilities
func GetCapabilities() *plugin.Capabilities {
	return &plugin.Capabilities{
		ActionNames: []string{
			actions.CreateKindClusterAction,
			actions.DeleteKindClusterAction,
			actions.LoadImageAction,
			actions.DeleteImageAction},
		IsModule: true,
	}
}
