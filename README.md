# Octant Plugin for kind

![CI](https://github.com/vmware-tanzu/octant-plugin-for-kind/workflows/ci/badge.svg?branch=main)

Octant Plugin for kind provides a visual interface for [kind](https://github.com/kubernetes-sigs/kind).
It provides a quick way to create and delete local development cluster from Octant's UI.

What this plugin can do:
- Set Kubernetes version
- Create HA control planes and multi-node clusters
- Enable feature gates
- Manage creating and deleting kind clusters
- Load and delete Docker images from a kind cluster

## Requirements
- Docker 19.03 or above
- Octant 0.18.0 or above
- kind 0.10 or above

## Known Issues
- After creating a new cluster, we need to find a way to know when the cluster is ready besides a context change
- What should happen to Octant when the current context is deleted?

## Feature Roadmap
- Export v1alpha4.Config as yaml
- Show disk usage and warn user to prune
- Reactive forms to show relevant feature gates based on k8s version
- Logging during cluster creation

## Contributing

Contributors will need to sign a DCO (Developer Certificate of Origin) with all changes. We also ask that a changelog entry is included with your pull request. Details are described in our [contributing](CONTRIBUTING.md) documentation.

See our [hacking](HACKING.md) guide for getting your development environment setup.

## License

Apache License v2.0: see [LICENSE](./LICENSE.txt) for details.
