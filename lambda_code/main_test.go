package main

import (
	"testing"
)

type sanitizeTest struct {
	inputString, outputString string
}

var sanitizeTests = []sanitizeTest{
	sanitizeTest{"foo", "foo"},
	sanitizeTest{"  ", "__"},
	sanitizeTest{" ,=", "___"},
	sanitizeTest{"count%", "count_percent"},
	sanitizeTest{"foo\tbar", "foo_bar"},
	sanitizeTest{"foo,bar%", "foo_bar_percent"},
	sanitizeTest{"/\\/:@<>“", "________"},
	sanitizeTest{"“", "_"},
	sanitizeTest{"prometheus metric count % is: 200", "prometheus_metric_count__percent_is__200"},
}

func TestSanitize(t *testing.T) {
	for _, test := range sanitizeTests {
		if output := sanitize(test.inputString); output != test.outputString {
			t.Errorf("Output %s not as expected %s", output, test.outputString)
		}
	}
}

type NamespaceLabelTest struct {
	namespace string
	expected  string
}

type SampleTest struct {
	value             Value
	timestamp         int64
	expectedValue     Value
	expectedTimestamp int64
}

var values = &Value{Count: 1, Sum: 2, Max: 3, Min: 4}

var sampleTest = []SampleTest{
	SampleTest{*values, 123, *values, 123},
}

func TestCreateSamples(t *testing.T) {
	for _, test := range sampleTest {
		countOutput := createCountSample(test.value, test.timestamp)
		if countOutput.Value != test.expectedValue.Count && countOutput.Timestamp != test.expectedTimestamp {
			t.Errorf("Output is not as expected")
		}

		maxOutput := createMaxSample(test.value, test.timestamp)
		if maxOutput.Value != test.expectedValue.Max && maxOutput.Timestamp != test.expectedTimestamp {
			t.Errorf("Output is not as expected")
		}

		minOutput := createMinSample(test.value, test.timestamp)
		if minOutput.Value != test.expectedValue.Min && minOutput.Timestamp != test.expectedTimestamp {
			t.Errorf("Output is not as expected")
		}

		sumOutput := createSumSample(test.value, test.timestamp)
		if sumOutput.Value != test.expectedValue.Sum && sumOutput.Timestamp != test.expectedTimestamp {
			t.Errorf("Output is not as expected")
		}
	}
}
