package settings

import (
	"github.com/vmware-tanzu/octant-plugin-for-kind/pkg/plugin/actions"
	"github.com/vmware-tanzu/octant-plugin-for-kind/pkg/plugin/router"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"strings"
)

func GetOptions() []service.PluginOption {
	return []service.PluginOption{
		service.WithActionHandler(actions.ActionHandler),
		service.WithNavigation(
			func(_ *service.NavigationRequest) (navigation.Navigation, error) {
				return navigation.Navigation{
					Title:    strings.Title(name),
					Path:     name,
					IconName: icon,
				}, nil
			},
			router.InitRoutes,
		),
	}
}
