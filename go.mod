module github.com/vmware-tanzu/octant-plugin-for-kind

go 1.16

require (
	github.com/Microsoft/go-winio v0.4.15 // indirect
	github.com/Microsoft/hcsshim v0.8.10 // indirect
	github.com/containerd/containerd v1.4.1 // indirect
	github.com/containerd/continuity v0.0.0-20201204194424-b0f312dbb49a // indirect
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v17.12.0-ce-rc1.0.20200916142827-bd33bbf0497b+incompatible
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	github.com/vmware-tanzu/octant v0.23.0
	k8s.io/client-go v0.21.3
	sigs.k8s.io/kind v0.11.1
)

replace golang.org/x/sys => golang.org/x/sys v0.0.0-20190830141801-acfa387b8d69
