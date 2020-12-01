package actions

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_FlexInt(t *testing.T) {
	var got FlexInt
	err := json.Unmarshal([]byte("1"), &got)
	require.NoError(t, err)

	require.Equal(t, int(got), 1)
}

func Test_FlexInt_invalid(t *testing.T) {
	var got FlexInt
	err := json.Unmarshal([]byte("abc"), &got)
	require.Error(t, err)
}