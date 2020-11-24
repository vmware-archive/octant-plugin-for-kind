package main // import "github.com/vmware-tanzu/octant-plugin-for-kind/cmd/octant-plugin-for-kind"

import (
	"github.com/vmware-tanzu/octant/pkg/plugin/service"

	"github.com/vmware-tanzu/octant-plugin-for-kind/pkg/plugin/settings"
)

func main() {
	name := settings.GetName()
	description := settings.GetDescription()
	capabilities := settings.GetCapabilities()
	options := settings.GetOptions()
	plugin, err := service.Register(name, description, capabilities, options...)
	if err != nil {
		panic(err)
	}
	plugin.Serve()
}