package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewStaticStore(t *testing.T) {
	_, err := newStaticStore(logger)
	require.NoError(t, err)
}
