package router

import (
	"github.com/vmware-tanzu/octant-plugin-for-kind/pkg/plugin/views"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// InitRoutes sets up plugin routes
func InitRoutes(router *service.Router) {
	router.HandleFunc("", kindListHandler)
	router.HandleFunc("/*", registryHandler)
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

func registryHandler(request service.Request) (component.ContentResponse, error) {
	view, err := views.BuildRegistryView(request)
	if err != nil {
		return component.EmptyContentResponse, nil
	}

	response := component.NewContentResponse(nil)
	response.Add(view)
	return *response, nil
}
