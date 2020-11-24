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
			"v1.19.1":  "kindest/node:v1.19.1@sha256:98cf5288864662e37115e362b23e4369c8c4a408f99cbc06e58ac30ddc721600",
			"v1.18.8":  "kindest/node:v1.18.8@sha256:f4bcc97a0ad6e7abaf3f643d890add7efe6ee4ab90baeb374b4f41a4c95567eb",
			"v1.17.11": "kindest/node:v1.17.11@sha256:5240a7a2c34bf241afb54ac05669f8a46661912eab05705d660971eeb12f6555",
			"v1.16.15": "kindest/node:v1.16.15@sha256:a89c771f7de234e6547d43695c7ab047809ffc71a0c3b65aa54eda051c45ed20",
			"v1.15.12": "kindest/node:v1.15.12@sha256:d9b939055c1e852fe3d86955ee24976cab46cba518abcb8b13ba70917e6547a6",
			"v1.14.10": "kindest/node:v1.14.10@sha256:ce4355398a704fca68006f8a29f37aafb49f8fc2f64ede3ccd0d9198da910146",
			"v1.13.12": "kindest/node:v1.13.12@sha256:1c1a48c2bfcbae4d5f4fa4310b5ed10756facad0b7a2ca93c7a4b5bae5db29f5",
		},
		keys: []string{"v1.19.1", "v1.18.8", "v1.17.11", "v1.16.15", "v1.15.12", "v1.14.10", "v1.13.12"},
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

type FeatureGate struct {
	Feature string  `json:"feature"`
	Default bool    `json:"default"`
	Stage   string  `json:"stage"`
	Since   float32 `json:"since"`
	Until   float32 `json:"until"`
}

// TODO: Some features are listed multiple times due to differing stages.
//  Remove once reactive forms are implemented
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
