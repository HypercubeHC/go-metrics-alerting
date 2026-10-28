package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemStorage(t *testing.T) {
	memStorage := New[int, any]()

	tests := []struct {
		value any
	}{
		{
			value: true,
		},
		{
			value: false,
		},
		{
			value: -1,
		},
		{
			value: 0,
		},
		{
			value: 1,
		},
		{
			value: -0.9,
		},
		{
			value: 2.2,
		},
		{
			value: "",
		},
		{
			value: 't',
		},
		{
			value: "test",
		},
		{
			value: map[string]interface{}{
				"bool":   true,
				"zero":   0,
				"float":  0.122,
				"string": "123abctest",
			},
		},
		{
			value: []string{"test0", "test1", "test2", "test3"},
		},
	}

	for name, val := range tests {
		memStorage.Set(name, val.value)
	}

	keys := make([]int, 0)
	for name, val := range tests {
		keys = append(keys, name)
		valFromStorage, ok := memStorage.Get(name)
		require.True(t, ok)
		require.Equal(t, val.value, valFromStorage)
	}
	require.Equal(t, len(keys), memStorage.Size())
}
