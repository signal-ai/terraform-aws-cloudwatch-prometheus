package main

import (
	"fmt"
	"testing"
)

type sanitizeTest struct {
	inputString, outputString string
}

var dimensions = map[string]string{
	"QueueName":    "asdasd-asd",
	"FunctionName": "bar",
}

func TestDimensions(t *testing.T) {
	output := createDimensionLabels(dimensions)
	fmt.Println(output)
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
