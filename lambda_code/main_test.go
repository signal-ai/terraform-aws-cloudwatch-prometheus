package main

import (
	"cmp"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/prometheus/prometheus/prompb"
)

type sanitizeTest struct {
	inputString, outputString string
}

var sanitizeTests = []sanitizeTest{
	{"foo", "foo"},
	{"  ", "__"},
	{" ,=", "___"},
	{"count%", "count_percent"},
	{"foo\tbar", "foo_bar"},
	{"foo,bar%", "foo_bar_percent"},
	{"/\\/:@<>“", "________"},
	{"“", "_"},
	{"prometheus metric count % is: 200", "prometheus_metric_count__percent_is__200"},
}

func sortLabelsFunc(a, b *prompb.Label) int {
	return cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
}

func TestSanitize(t *testing.T) {
	for _, test := range sanitizeTests {
		if output := sanitize(test.inputString); output != test.outputString {
			t.Errorf("Output %s not as expected %s", output, test.outputString)
		}
	}
}

func TestSnakeCase(t *testing.T) {
	expected := "foo_bar"
	output := toSnakeCase("FooBar")

	if expected != output {
		t.Errorf("Output %s not as expected %s", output, expected)
	}
}

func TestCreateMetricNameLabels(t *testing.T) {
	customOutput := createMetricNameLabels("foo", "bar", "count", "eu-west-1", "dev")
	expected := []*prompb.Label{
		{
			Name:  "__name__",
			Value: "aws_custom_bar_foo_count",
		},
		{
			Name:  "region",
			Value: "eu-west-1",
		},
		{
			Name:  "account",
			Value: "dev",
		},
	}

	if !reflect.DeepEqual(expected, customOutput) {
		t.Errorf("Output %v not as expected %v", customOutput, expected)
	}

	output := createMetricNameLabels("foo", "AWS/bar", "count", "eu-west-1", "dev")
	expected = []*prompb.Label{
		{
			Name:  "__name__",
			Value: "aws_bar_foo_count",
		},
		{
			Name:  "region",
			Value: "eu-west-1",
		},
		{
			Name:  "account",
			Value: "dev",
		},
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Output %v not as expected %v", output, expected)
	}
}

func TestCreateDimensionLabels(t *testing.T) {
	output := createDimensionLabels(map[string]string{
		"foo":    "bar",
		"baz":    "qux",
		"ignore": "",
	})
	expected := []*prompb.Label{
		{
			Name:  "foo",
			Value: "bar",
		},
		{
			Name:  "baz",
			Value: "qux",
		},
	}

	slices.SortFunc(expected, sortLabelsFunc)

	slices.SortFunc(output, sortLabelsFunc)

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Output %v not as expected %v", output, expected)
	}
}

func TestHandleAddLabels(t *testing.T) {
	dimensions := map[string]string{
		"foo": "bar",
		"baz": "qux",
	}
	output := handleAddLabels("count", "foo", "bar", dimensions, "eu-west-1", "dev")
	expected := []*prompb.Label{
		{
			Name:  "foo",
			Value: "bar",
		},
		{
			Name:  "baz",
			Value: "qux",
		},
		{
			Name:  "__name__",
			Value: "aws_custom_bar_foo_count",
		},
		{
			Name:  "region",
			Value: "eu-west-1",
		},
		{
			Name:  "account",
			Value: "dev",
		},
	}

	slices.SortFunc(expected, sortLabelsFunc)

	slices.SortFunc(output, sortLabelsFunc)

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Output %v not as expected %v", output, expected)
	}
}

func TestHandleAddSample(t *testing.T) {
	output := handleAddSamples("count", Value{Count: 42}, 1234)
	expected := prompb.Sample{
		Value:     42,
		Timestamp: 1234,
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Output %v not as expected %v", output, expected)
	}
}
