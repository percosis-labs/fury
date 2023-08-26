package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetDiff(t *testing.T) {
	tests := []struct {
		name     string
		setA     []string
		setB     []string
		expected []string
	}{
		{"empty", []string{}, []string{}, []string(nil)},
		{"diff equal sets", []string{"busd", "usdf"}, []string{"busd", "usdf"}, []string(nil)},
		{"diff set empty", []string{"bnb", "ufury", "usdf"}, []string{}, []string{"bnb", "ufury", "usdf"}},
		{"input set empty", []string{}, []string{"bnb", "ufury", "usdf"}, []string(nil)},
		{"diff set with common elements", []string{"bnb", "btcb", "usdf", "xrpb"}, []string{"bnb", "usdf"}, []string{"btcb", "xrpb"}},
		{"diff set with all common elements", []string{"bnb", "usdf"}, []string{"bnb", "btcb", "usdf", "xrpb"}, []string(nil)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, setDifference(tt.setA, tt.setB))
		})
	}
}
