package views

import (
	"encoding/json"
)

// OrderedMap is an ordered map for kind images
type OrderedMap struct {
	m    map[string]string
	keys []string
}

// NewImageMap creates an instance of OrderedMap
func NewImageMap() *OrderedMap {
	return &OrderedMap{
		m: map[string]string{
			"v1.20.2": "kindest/node:v1.20.2@sha256:8f7ea6e7642c0da54f04a7ee10431549c0257315b3a634f6ef2fecaaedb19bab",
			"v1.19.7": "kindest/node:v1.19.7@sha256:a70639454e97a4b733f9d9b67e12c01f6b0297449d5b9cbbef87473458e26dca",
			"v1.18.15": "kindest/node:v1.18.15@sha256:5c1b980c4d0e0e8e7eb9f36f7df525d079a96169c8a8f20d8bd108c0d0889cc4",
			"v1.17.17": "kindest/node:v1.17.17@sha256:7b6369d27eee99c7a85c48ffd60e11412dc3f373658bc59b7f4d530b7056823e",
			"v1.16.15": "kindest/node:v1.16.15@sha256:a89c771f7de234e6547d43695c7ab047809ffc71a0c3b65aa54eda051c45ed20",
			"v1.15.12": "kindest/node:v1.15.12@sha256:d9b939055c1e852fe3d86955ee24976cab46cba518abcb8b13ba70917e6547a6",
			"v1.14.10": "kindest/node:v1.14.10@sha256:ce4355398a704fca68006f8a29f37aafb49f8fc2f64ede3ccd0d9198da910146",
		},
		keys: []string{"v1.20.2", "v1.19.7", "v1.18.15", "v1.17.17", "v1.16.15", "v1.15.12", "v1.14.10"},
	}
}

// Keys returns a list of keys
func (o *OrderedMap) Keys() []string {
	return o.keys
}

// Map returns a map
func (o *OrderedMap) Map() map[string]string {
	return o.m
}

func getFeatureGateList() ([]FeatureGate, error) {
	var featureGates []FeatureGate
	if err := json.Unmarshal([]byte(FeatureList), &featureGates); err != nil {
		return nil, err
	}

	return featureGates, nil
}

// FeatureGate contains metadata of a feature gate
type FeatureGate struct {
	Feature string  `json:"feature"`
	Default bool    `json:"default"`
	Stage   string  `json:"stage"`
	Since   float32 `json:"since"`
	Until   float32 `json:"until"`
}

// TODO: Some features are listed multiple times due to differing stages.
//  Remove once reactive forms are implemented

// Unique returns a list of unique feature gate names
func Unique(featureGates []FeatureGate) []FeatureGate {
	keys := make(map[string]bool)
	var list []FeatureGate
	for _, fg := range featureGates {
		if _, name := keys[fg.Feature]; !name {
			keys[fg.Feature] = true
			list = append(list, fg)
		}
	}
	return list
}
