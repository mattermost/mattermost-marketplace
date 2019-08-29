package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewStatikStore(t *testing.T) {
	_, err := newStatikStore(logger)
	require.NoError(t, err)
}
