This directory contains a script used to convert feature gate status HTML tables in 
[https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/](https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/)
to JSON.

Running this script assumes an environment where `requests` and `pandas` are installed.

`feature_gates.go` will be generated with JSON embedded.

Before each release, the JSON files will need to be checked for changes and regenerate if needed.

A better solution in the future is to consume [https://github.com/kubernetes/kubernetes/blob/master/pkg/features/kube_features.go](https://github.com/kubernetes/kubernetes/blob/master/pkg/features/kube_features.go) directly.