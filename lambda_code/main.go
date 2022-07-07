package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"
)

type Dimensions = map[string]string

type MetricStreamData struct {
	MetricStreamName string     `json:"metric_stream_name"`
	AccountID        string     `json:"account_id"`
	Region           string     `json:"region"`
	Namespace        string     `json:"namespace"`
	MetricName       string     `json:"metric_name"`
	Dimensions       Dimensions `json:"dimensions"`
	Timestamp        int64      `json:"timestamp"`
	Value            Value      `json:"value"`
	Unit             string     `json:"unit"`
}

type Value struct {
	Count float64 `json:"count"`
	Sum   float64 `json:"sum"`
	Max   float64 `json:"max"`
	Min   float64 `json:"min"`
}

type Values string

const (
	Count Values = "count"
	Sum          = "sum"
	Max          = "max"
	Min          = "min"
)

func HandleRequest(ctx context.Context, evnt events.KinesisFirehoseEvent) (events.KinesisFirehoseResponse, error) {

	var response events.KinesisFirehoseResponse
	var timeSeries []*prompb.TimeSeries
	// These are the 4 value types from Cloudwatch, each of which map to a Prometheus Gauge
	values := []Values{Count, Max, Min, Sum}

	for _, record := range evnt.Records {

		splitRecord := strings.Split(string(record.Data), string('\n'))
		for _, x := range splitRecord {

			// The Records includes an empty new line at the last position which becomes "" after parsing. Skipping over the empty string.
			if x == "" {
				continue
			}
			var metricStreamData MetricStreamData
			json.Unmarshal([]byte(x), &metricStreamData)

			// For each metric, the labels + valuetype is the __name__ of the sample, and the corresponding single sample value is used to create the timeseries.
			for _, value := range values {
				var samples []prompb.Sample
				currentLabels := handleAddLabels(value, metricStreamData.MetricName, metricStreamData.Namespace, metricStreamData.Dimensions, metricStreamData.Region, metricStreamData.AccountID)
				currentSamples := handleAddSamples(value, metricStreamData.Value, metricStreamData.Timestamp)
				samples = append(samples, currentSamples)

				singleTimeSeries := &prompb.TimeSeries{
					Labels:  currentLabels,
					Samples: samples,
				}

				timeSeries = append(timeSeries, singleTimeSeries)
			}
		}

		// No transformation occurs, just send OK response back to Kinesis
		var transformedRecord events.KinesisFirehoseResponseRecord
		transformedRecord.RecordID = record.RecordID
		transformedRecord.Result = events.KinesisFirehoseTransformedStateOk
		transformedRecord.Data = []byte(string(record.Data))

		response.Records = append(response.Records, transformedRecord)
	}

	err := createWriteRequestAndSendToAPS(timeSeries)
	return response, err
}

func main() {
	lambda.Start(HandleRequest)
}

// Taken directly from YACE: https://github.com/nerdswords/yet-another-cloudwatch-exporter/blob/1c7b3d7b7b64ce93bb4a27d8ef836e0c2b96b8e7/pkg/prometheus.go#L139
func sanitize(text string) string {
	replacer := strings.NewReplacer(
		" ", "_",
		",", "_",
		"\t", "_",
		"/", "_",
		"\\", "_",
		".", "_",
		"-", "_",
		":", "_",
		"=", "_",
		"“", "_",
		"@", "_",
		"<", "_",
		">", "_",
		"%", "_percent",
	)
	return replacer.Replace(text)
}

func toSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func handleAddLabels(valueType Values, metricName string, namespace string, dimensions Dimensions, region string, account string) []*prompb.Label {

	var labels []*prompb.Label

	metricNameLabels := createMetricNameLabels(metricName, namespace, valueType, region, account)
	dimensionLabels := createDimensionLabels(dimensions)
	labels = append(labels, dimensionLabels...)
	labels = append(labels, metricNameLabels...)
	return labels
}

func handleAddSamples(valueType Values, value Value, timestamp int64) prompb.Sample {
	var sample prompb.Sample
	switch valueType {
	case Count:
		sample = createCountSample(value, timestamp)
	case Min:
		sample = createMinSample(value, timestamp)
	case Max:
		sample = createMaxSample(value, timestamp)
	case Sum:
		sample = createSumSample(value, timestamp)
	}
	return sample
}

func createMetricNameLabels(metricName string, namespace string, valueType Values, region string, account string) []*prompb.Label {
	var labels []*prompb.Label
	metricNameLabel := &prompb.Label{
		Name:  "__name__",
		Value: strings.ToLower(sanitize(namespace) + "_" + sanitize(toSnakeCase(metricName)) + "_" + sanitize(string(valueType))),
	}
	labels = append(labels, metricNameLabel)
	regionLabel := &prompb.Label{
		Name:  "region",
		Value: region,
	}
	labels = append(labels, regionLabel)
	accountLabel := &prompb.Label{
		Name:  "account",
		Value: sanitize(account),
	}
	labels = append(labels, accountLabel)
	return labels
}

func createDimensionLabels(dimensions Dimensions) []*prompb.Label {
	var labels []*prompb.Label

	// for all dimensions in dimensions map, create a label with the dimension name and value
	// if element is not "" then create a label with the dimension name and value
	for key, value := range dimensions {
		if value != "" {
			dimensionLabel := &prompb.Label{
				Name:  sanitize(toSnakeCase(key)),
				Value: sanitize(value),
			}
			labels = append(labels, dimensionLabel)
		}
	}

	return labels
}

func createSumSample(value Value, timestamp int64) prompb.Sample {
	sumSample := prompb.Sample{
		Value:     value.Sum,
		Timestamp: timestamp,
	}
	return sumSample
}

func createCountSample(value Value, timestamp int64) prompb.Sample {
	countSample := prompb.Sample{
		Value:     value.Count,
		Timestamp: timestamp,
	}
	return countSample
}

func createMaxSample(value Value, timestamp int64) prompb.Sample {
	maxSample := prompb.Sample{
		Value:     value.Max,
		Timestamp: timestamp,
	}
	return maxSample
}

func createMinSample(value Value, timestamp int64) prompb.Sample {
	minSample := prompb.Sample{
		Value:     value.Min,
		Timestamp: timestamp,
	}
	return minSample
}

func createWriteRequestAndSendToAPS(timeseries []*prompb.TimeSeries) error {
	writeRequest := &prompb.WriteRequest{
		Timeseries: timeseries,
	}

	body := encodeWriteRequestIntoProtoAndSnappy(writeRequest)
	err := sendRequestToAPS(body)
	return err
}

func encodeWriteRequestIntoProtoAndSnappy(writeRequest *prompb.WriteRequest) *bytes.Reader {
	data, err := proto.Marshal(writeRequest)

	if err != nil {
		panic(err)
	}

	encoded := snappy.Encode(nil, data)
	body := bytes.NewReader(encoded)
	return body
}

func sendRequest(url string, bodyBytes []byte) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	var netClient = &http.Client{Timeout: time.Second * 5}
	resp, err := netClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		return fmt.Errorf("Error: statuscode %d sending data to endpoint: %s, %s", resp.StatusCode, url, bodyString)
	}
	return nil
}

func sendRequestToAPS(body *bytes.Reader) error {
	bodyBytes, _ := io.ReadAll(body)
	var errors []string
	endpoints := strings.Split(os.Getenv("PROMETHEUS_REMOTE_WRITE_URLS"), ",")
	for _, url := range endpoints {
		err := sendRequest(url, bodyBytes)
		if err != nil {
			errors = append(errors, err.Error())
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, ","))
	}
	return nil
}
