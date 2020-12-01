package views

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_getFeatureGateList(t *testing.T) {
	_, err := getFeatureGateList()
	require.NoError(t, err)
}

func Test_Unique(t *testing.T) {
	cases := []struct {
		name         string
		featureGates []FeatureGate
		expected     []FeatureGate
	}{
		{
			name:         "empty list",
			featureGates: []FeatureGate{},
			expected:     nil,
		},
		{
			name: "duplicate",
			featureGates: []FeatureGate{
				{
					Feature: "test",
					Stage:   "alpha",
				},
				{
					Feature: "test",
					Stage:   "beta",
				},
			},
			expected: []FeatureGate{
				{
					Feature: "test",
					Stage:   "alpha",
				},
			},
		},
		{
			name: "unique",
			featureGates: []FeatureGate{
				{
					Feature: "a",
					Stage:   "alpha",
				},
				{
					Feature: "b",
					Stage:   "alpha",
				},
			},
			expected: []FeatureGate{
				{
					Feature: "a",
					Stage:   "alpha",
				},
				{
					Feature: "b",
					Stage:   "alpha",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Unique(tc.featureGates)
			require.EqualValues(t, got, tc.expected)
		})
	}
}
