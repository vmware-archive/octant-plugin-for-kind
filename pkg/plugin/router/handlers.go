package router

import (
	"github.com/vmware-tanzu/octant-plugin-for-kind/pkg/plugin/views"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func InitRoutes(router *service.Router) {
	router.HandleFunc("", kindListHandler)
}

func kindListHandler(request service.Request) (component.ContentResponse, error) {
	view, err := views.BuildKindClusterView(request)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	response := component.NewContentResponse(nil)
	response.Add(view)
	return *response, nil
}
