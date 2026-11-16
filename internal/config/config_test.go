package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	addressVarName     = "ADDRESS"
	reportVarName      = "REPORT_INTERVAL"
	pollingVarName     = "POLL_INTERVAL"
	defaulTestUrl      = "localhost:8080"
	defaultTestReport  = 10
	defaultTestPolling = 2
)

func TestNewConfigDefaults(t *testing.T) {
	testConfig := NewConfig()
	assert.Equal(t, defaulTestUrl, testConfig.URL)
	assert.Equal(t, defaultTestReport, testConfig.Report)
	assert.Equal(t, defaultTestPolling, testConfig.Polling)
}

func TestSetVars(t *testing.T) {
	tests := []struct {
		name            string
		addressENV      string
		reportENV       string
		pollingENV      string
		urlFlag         string
		reportFlag      int
		pollingFlag     int
		expectedUrl     string
		expectedReport  int
		expectedPolling int
	}{
		{
			name:            "Piority for values from ENV",
			addressENV:      "localhost:1111",
			reportENV:       "15",
			pollingENV:      "6",
			urlFlag:         defaulTestUrl,
			reportFlag:      defaultTestReport,
			pollingFlag:     defaultTestPolling,
			expectedUrl:     "localhost:1111",
			expectedReport:  15,
			expectedPolling: 6,
		},
		{
			name:            "Priority for values from flags if from ENV incorrect",
			reportENV:       "test",
			pollingENV:      "test",
			addressENV:      "",
			urlFlag:         defaulTestUrl,
			reportFlag:      defaultTestReport,
			pollingFlag:     defaultTestPolling,
			expectedUrl:     defaulTestUrl,
			expectedReport:  defaultTestReport,
			expectedPolling: defaultTestPolling,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv(addressVarName, test.addressENV)
			t.Setenv(reportVarName, test.reportENV)
			t.Setenv(pollingVarName, test.pollingENV)

			url, report, polling := setVars(test.urlFlag, test.reportFlag, test.pollingFlag)
			assert.Equal(t, test.expectedUrl, url)
			assert.Equal(t, test.expectedReport, report)
			assert.Equal(t, test.expectedPolling, polling)
		})
	}
}
