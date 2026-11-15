package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var gaugeCorrectTests = []struct {
	value    Gauge
	valueStr string
}{
	{
		value:    0,
		valueStr: "0",
	},
	{
		value:    0.1,
		valueStr: "0.1",
	},
	{
		value:    -0.1,
		valueStr: "-0.1",
	},
	{
		value:    0.00001,
		valueStr: "0.00001",
	},
	{
		value:    1.0,
		valueStr: "1",
	},
	{
		value:    -1.0,
		valueStr: "-1",
	},
	{
		value:    9999999999.01,
		valueStr: "9999999999.01",
	},
	{
		value:    -9999999999.01,
		valueStr: "-9999999999.01",
	},
	{
		value:    -9999999999.000000000000001,
		valueStr: "-9999999999",
	},
	{
		value:    0.00000000000000000000000000000000001,
		valueStr: "0.00000000000000000000000000000000001",
	},
	{
		value:    -0.00000000000000000000000000000000001,
		valueStr: "-0.00000000000000000000000000000000001",
	},
}

var gaugeIncorrectTests = []struct {
	value string
}{
	{value: "abc"},
	{value: "test"},
	{value: " -0.1b"},
	{value: " 0.1b"},
	{value: "0.2d "},
	{value: "a.0001"},
}

var counterCorrectTests = []struct {
	value    Counter
	valueStr string
}{
	{
		value:    0,
		valueStr: "0",
	},
	{
		value:    1,
		valueStr: "1",
	},
	{
		value:    -1,
		valueStr: "-1",
	},
	{
		value:    10000000,
		valueStr: "10000000",
	},
	{
		value:    -10000000,
		valueStr: "-10000000",
	},
	{
		value:    2982723873528,
		valueStr: "2982723873528",
	},
	{
		value:    -2982723873528,
		valueStr: "-2982723873528",
	},
}

var counterIncorrectTests = []struct {
	value string
}{
	{value: ""},
	{value: " "},
	{value: "abc"},
	{value: "test"},
	{value: " -1c"},
	{value: " 2b"},
	{value: "0.1"},
	{value: "--1"},
}

func TestGaugeString(t *testing.T) {
	for _, test := range gaugeCorrectTests {
		assert.Equal(t, test.valueStr, Gauge(test.value).String())
	}
}

func TestCounterString(t *testing.T) {
	for _, test := range counterCorrectTests {
		assert.Equal(t, test.valueStr, Counter(test.value).String())
	}
}

func TestGaugeBytesConvert(t *testing.T) {
	for _, test := range gaugeCorrectTests {
		assert.Equal(t, test.value, BytesToGauge(GaugeToBytes(test.value)))
	}
}

func TestCounterBytesConvert(t *testing.T) {
	for _, test := range counterCorrectTests {
		assert.Equal(t, test.value, BytesToCounter(CounterToBytes(test.value)))
	}
}

func TestParseCorrectGauge(t *testing.T) {
	for _, test := range gaugeCorrectTests {
		parsed, err := ParseGauge(test.value.String())
		require.NoError(t, err)
		assert.Equal(t, Gauge(test.value), parsed)
	}
}

func TestParseIncorrectGauge(t *testing.T) {
	for _, test := range gaugeIncorrectTests {
		_, err := ParseGauge(test.value)
		assert.Error(t, err, "Value: "+test.value)
	}
}

func TestParseCorrectCounter(t *testing.T) {
	for _, test := range counterCorrectTests {
		parsed, err := ParseCounter(test.value.String())
		require.NoError(t, err)
		assert.Equal(t, Counter(test.value), parsed)
	}
}

func TestParseIncorrectCounter(t *testing.T) {
	for _, test := range counterIncorrectTests {
		_, err := ParseCounter(test.value)
		assert.Error(t, err)
	}
}
