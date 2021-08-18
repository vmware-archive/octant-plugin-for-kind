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
			"v1.21.1": "kindest/node:v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6",
			"v1.20.7": "kindest/node:v1.20.7@sha256:cbeaf907fc78ac97ce7b625e4bf0de16e3ea725daf6b04f930bd14c67c671ff9",
			"v1.19.11": "kindest/node:v1.19.11@sha256:07db187ae84b4b7de440a73886f008cf903fcf5764ba8106a9fd5243d6f32729",
			"v1.18.19": "kindest/node:v1.18.19@sha256:7af1492e19b3192a79f606e43c35fb741e520d195f96399284515f077b3b622c",
			"v1.17.17": "kindest/node:v1.17.17@sha256:66f1d0d91a88b8a001811e2f1054af60eef3b669a9a74f9b6db871f2f1eeed00",
			"v1.16.15": "kindest/node:v1.16.15@sha256:83067ed51bf2a3395b24687094e283a7c7c865ccc12a8b1d7aa673ba0c5e8861",
			"v1.15.12": "kindest/node:v1.15.12@sha256:b920920e1eda689d9936dfcf7332701e80be12566999152626b2c9d730397a95",
			"v1.14.10": "kindest/node:v1.14.10@sha256:f8a66ef82822ab4f7569e91a5bccaf27bceee135c1457c512e54de8c6f7219f8",
		},
		keys: []string{"v1.21.1", "v1.20.7", "v1.19.11", "v1.18.19", "v1.17.17", "v1.16.15", "v1.15.12", "1.14.10"},
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
