package storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestMemStorageGaugeDelete(t *testing.T) {
	memStorage := New[string, Gauge]()
	memStorage.Set("gaugeTest", Gauge(1.1))
	memStorage.Delete("gaugeTest")

	_, ok := memStorage.Get("gaugeTest")
	assert.False(t, ok)
	assert.Equal(t, 0, memStorage.Size())
}

func TestMemStorageCounterDelete(t *testing.T) {
	memStorage := New[string, Counter]()
	memStorage.Set("counterTest", Counter(1))
	memStorage.Delete("counterTest")

	_, ok := memStorage.Get("counterTest")
	assert.False(t, ok)
	assert.Equal(t, 0, memStorage.Size())
}

func TestMemStorageKeys(t *testing.T) {
	memStorage := New[string, int]()
	memStorage.Set("test1", 1)
	memStorage.Set("test2", 2)

	assert.ElementsMatch(t, []string{"test1", "test2"}, memStorage.Keys())
	assert.Empty(t, New[string, int]().Keys())
}

func TestMemStorageIterDoesNotHoldLock(t *testing.T) {
	memStorage := New[string, int]()
	memStorage.Set("a", 1)

	done := make(chan error, 1)
	go func() {
		done <- memStorage.Iter(context.Background(), func(name string, val int) error {
			memStorage.Set("b", 2) // при удержании RLock внутри Iter был бы дедлок
			return nil
		})
	}()
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("Iter holds RLock during callback (deadlock)")
	}
	assert.Equal(t, 2, memStorage.Size())
}
