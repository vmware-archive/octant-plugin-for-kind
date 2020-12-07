package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	dockerClient "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/pkg/system"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

type KindImage struct {
	ID          string   `json:"id"`
	UID         ID       `json:"uid"` // UID can be null
	RepoTags    []string `json:"repoTags"`
	RepoDigests []string `json:"repoDigests"`
	Size        string   `json:"size"`
	Username    string   `json:"username"`
}

type ID struct {
	value string
}

type KindImages struct {
	Images []KindImage `json:"images"`
}

type Image struct {
	Containers   string
	CreatedAt    string
	CreatedSince string
	Digest       string
	ID           string
	Repository   string
	SharedSize   string
	Size         string
	Tag          string
	UniqueSize   string
	VirtualSize  string
}

// ContainerDetails contains metadata of a docker container
type ContainerDetails struct {
	Version string
	Created int64
	State   string
}

// DockerClient is a docker client
type Client struct {
	client *dockerClient.Client
	ctx    context.Context
}

func NewDockerClient() *Client {
	client, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv, dockerClient.WithAPIVersionNegotiation())
	if err != nil {
		return nil
	}

	return &Client{
		client: client,
	}
}

func (d *Client) KindControlPlaneContainer(ctx context.Context, clusterName string) (*ContainerDetails, error) {
	if d.client == nil {
		return nil, errors.New("docker client is nil")
	}

	defer d.client.Close()
	containerDetails := &ContainerDetails{}

	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}

	for _, container := range containers {
		containerDetails.Created = container.Created
		containerDetails.State = container.State

		for _, name := range container.Names {
			if "/"+clusterName+"-control-plane" == name {
				r, _ := regexp.Compile("v\\d+\\.\\d+\\.\\d+")
				containerDetails.Version = r.FindString(container.Image)
				return containerDetails, nil
			}
		}
	}
	return containerDetails, errors.Errorf("cannot find details for cluster: %s", clusterName)
}

func (d *Client) ListKindImages(ctx context.Context, clusterName string) (KindImages, error) {
	if d.client == nil {
		return KindImages{}, errors.New("docker client is nil")
	}

	defer d.client.Close()
	execConfig := &types.ExecConfig{
		Tty:          false,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		// User: "root",
		Cmd: []string{
			"crictl",
			"images",
			"--output=json",
		},
	}
	response, err := d.client.ContainerExecCreate(ctx, clusterName+"-control-plane", *execConfig)
	if err != nil {
		return KindImages{}, err
	}

	h, err := d.client.ContainerExecAttach(ctx, response.ID, types.ExecStartCheck{})
	if err != nil {
		return KindImages{}, err
	}

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	if _, err := stdcopy.StdCopy(stdout, stderr, h.Reader); err != nil {
		return KindImages{}, err
	}

	if stderr.Len() > 0 {
		return KindImages{}, errors.New(stderr.String())
	}

	var images KindImages
	err = json.Unmarshal(stdout.Bytes(), &images)
	if err != nil {
		return KindImages{}, err
	}

	return images, nil
}

func (d *Client) ListDockerImages(ctx context.Context) ([]types.ImageSummary, error) {
	if d.client == nil {
		return nil, errors.New("docker client is nil")
	}
	defer d.client.Close()

	images, err := d.client.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (d *Client) DeleteKindImage(ctx context.Context, clusterName, imageID string) error {
	if d.client == nil {
		return errors.New("docker client is nil")
	}

	defer d.client.Close()
	execConfig := &types.ExecConfig{
		Tty:          false,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		// User: "root",
		Cmd: []string{
			"crictl",
			"rmi",
			imageID,
		},
	}
	response, err := d.client.ContainerExecCreate(ctx, clusterName+"-control-plane", *execConfig)
	if err != nil {
		return err
	}
	if _, err := d.client.ContainerExecAttach(ctx, response.ID, types.ExecStartCheck{}); err != nil {
		return err
	}
	return nil
}

func (d *Client) Save(ctx context.Context, destination, imageID string) error {
	if d.client == nil {
		return errors.New("docker client is nil")
	}

	defer d.client.Close()

	imageIDs := []string{
		imageID,
	}
	response, err := d.client.ImageSave(ctx, imageIDs)
	if err != nil {
		return err
	}

	defer response.Close()

	return 	copyToFile(destination, response)
}

// Copied from https://github.com/docker/cli/blob/51a091485d69c71bcbfd9a1b358d551c3d57504e/cli/command/utils.go
// copyToFile writes the content of the reader to the specified file
func copyToFile(outfile string, r io.Reader) error {
	// We use sequential file access here to avoid depleting the standby list
	// on Windows. On Linux, this is a call directly to ioutil.TempFile
	tmpFile, err := system.TempFileSequential(filepath.Dir(outfile), ".docker_temp_")
	if err != nil {
		return err
	}

	tmpPath := tmpFile.Name()

	_, err = io.Copy(tmpFile, r)
	tmpFile.Close()

	if err != nil {
		os.Remove(tmpPath)
		return err
	}

	if err = os.Rename(tmpPath, outfile); err != nil {
		os.Remove(tmpPath)
		return err
	}

	return nil
}
