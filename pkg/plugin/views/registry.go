package views

import (
	"fmt"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stringid"
	units "github.com/docker/go-units"
	"github.com/vmware-tanzu/octant-plugin-for-kind/pkg/docker"
	"github.com/vmware-tanzu/octant-plugin-for-kind/pkg/plugin/actions"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"strconv"
	"strings"
	"time"
)

func BuildRegistryView(request service.Request) (component.Component, error) {
	ctx := request.Context()
	client := docker.NewDockerClient()
	name := strings.TrimPrefix(request.Path(), "/")

	flexLayout := component.NewFlexLayout("Local Images")

	table := component.NewTable("Docker Images", "No images found",
		component.NewTableCols("Repository", "Tag", "Image ID", "Age", "Size"))

	dockerImages, err := client.ListDockerImages(ctx)
	if err != nil {
		return flexLayout, nil
	}

	for _, image := range dockerImages {
		rowPrinter(name, image, table)
	}

	kindTable := component.NewTable("Kind Images", "No images found",
		component.NewTableCols("Image", "Image ID", "Size"))

	kindImages, err := client.ListKindImages(ctx, name)
	if err != nil {
		return flexLayout, nil
	}

	for _, image := range kindImages.Images {
		for _, repoTag := range image.RepoTags {
			kindTable.Add(kindPrinter(name, image, repoTag))
		}
	}

	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: table},
		{Width: component.WidthFull, View: kindTable},
	})
	return flexLayout, nil
}

func kindPrinter(clusterName string, image docker.KindImage, repoTag string) component.TableRow {
	row := component.TableRow{}

	var humanReadableSize string
	size, err := strconv.ParseFloat(image.Size, 64)
	if err != nil {
		humanReadableSize = "<unknown>"
	} else {
		humanReadableSize = units.HumanSize(size)
	}
	row["Image"] = component.NewText(fmt.Sprintf("%s", repoTag))
	row["Image ID"] = component.NewText(stringid.TruncateID(image.ID))
	row["Size"] = component.NewText(humanReadableSize)

	confirmation := &component.Confirmation{
		Title: "Are you sure?",
		Body:  fmt.Sprintf("Do you want to delete %s from your kind images?", repoTag),
	}

	action := component.GridAction{
		Name:       "Delete",
		ActionPath: actions.DeleteImageAction,
		Payload: action.Payload{
			"action":      actions.DeleteImageAction,
			"imageID":     stringid.TruncateID(image.ID),
			"clusterName": clusterName,
		},
		Confirmation: confirmation,
		Type:         component.GridActionDanger,
	}

	row.AddAction(action)
	return row
}

func rowPrinter(clusterName string, image types.ImageSummary, table *component.Table) {
	repoTags := map[string][]string{}

	for _, refString := range image.RepoTags {
		ref, err := reference.ParseNormalizedNamed(refString)
		if err != nil {
			continue
		}
		if nt, ok := ref.(reference.NamedTagged); ok {
			familiarRef := reference.FamiliarName(ref)
			repoTags[familiarRef] = append(repoTags[familiarRef], nt.Tag())
		}
	}

	for repo, tags := range repoTags {
		for _, tag := range tags {
			row := component.TableRow{
				"Repository": component.NewText(repo),
				"Tag":        component.NewText(tag),
				"Image ID":   component.NewText(stringid.TruncateID(image.ID)),
				"Age":        component.NewTimestamp(time.Unix(image.Created, 0)),
				"Size":       component.NewText(units.HumanSize(float64(image.VirtualSize))),
			}

			action := component.GridAction{
				Name:       "Load into Kind",
				ActionPath: actions.LoadImageAction,
				Payload: action.Payload{
					"action":      actions.LoadImageAction,
					"imageName":   repo + ":" + tag,
					"clusterName": clusterName,
				},
				Type: component.GridActionPrimary,
			}
			row.AddAction(action)
			table.Add(row)
		}
	}
	return
}
