package insightfinder

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const METRIC_DATA_API = "/api/v2/metric-data-receive"
const CHUNK_SIZE = 2 * 1024 * 1024
const MAX_PACKET_SIZE = 10000000
const HTTP_RETRY_TIMES = 15
const HTTP_RETRY_INTERVAL = 60

type InsightFinderClient struct {
	Url        string
	Username   string
	LicenseKey string
	Project    string
}

func CreateInsightFinderClient(url, username, licenseKey, project string) *InsightFinderClient {
	return &InsightFinderClient{
		Url:        url,
		Username:   username,
		LicenseKey: licenseKey,
		Project:    project,
	}
}

func (client *InsightFinderClient) SendMetricData(instanceDataMap *InstanceDataMap) {
	curTotal := 0
	var newPayload = MetricDataReceivePayload{
		ProjectName:     client.Project,
		UserName:        client.Username,
		InstanceDataMap: *instanceDataMap,
	}
	for instanceName, istData := range *instanceDataMap {
		instanceData, ok := newPayload.InstanceDataMap[instanceName]
		if !ok {
			// Current NodeInstance didn't exist
			instanceData = InstanceData{
				InstanceName:       istData.InstanceName,
				ComponentName:      istData.ComponentName,
				DataInTimestampMap: make(map[int64]DataInTimestamp),
			}
			newPayload.InstanceDataMap[instanceName] = instanceData
		}
		for timeStamp, tsData := range istData.DataInTimestampMap {
			// Need to send out the data in the same timestamp in one payload
			dataBytes, err := json.Marshal(tsData)
			if err != nil {
				panic("[ERORR] There's issue form json data for DataInTimestampMap.")
			}
			// Add the data into the payload
			instanceData.DataInTimestampMap[timeStamp] = tsData
			// The json.Marshal transform the data into bytes so the length will be the actual size.
			curTotal += len(dataBytes)
			if curTotal > CHUNK_SIZE {
				request := IFMetricPostRequestPayload{
					LicenseKey: client.LicenseKey,
					UserName:   client.Username,
					Data:       newPayload,
				}
				jData, err := json.Marshal(request)
				if err != nil {
					panic(err)
				}
				client.sendDataToIF(jData, METRIC_DATA_API)
				curTotal = 0
				newPayload = MetricDataReceivePayload{
					ProjectName:     client.Project,
					UserName:        client.Username,
					InstanceDataMap: make(map[string]InstanceData),
				}
				newPayload.InstanceDataMap[instanceName] = InstanceData{
					InstanceName:       istData.InstanceName,
					ComponentName:      istData.InstanceName,
					DataInTimestampMap: make(map[int64]DataInTimestamp),
				}
			}
		}
	}
	request := IFMetricPostRequestPayload{
		LicenseKey: client.LicenseKey,
		UserName:   client.Username,
		Data:       newPayload,
	}
	jData, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	client.sendDataToIF(jData, METRIC_DATA_API)
}

func (client *InsightFinderClient) sendDataToIF(data []byte, receiveEndpoint string) {
	slog.Info("-------- Sending data to InsightFinder --------")

	if len(data) > MAX_PACKET_SIZE {
		panic("[ERROR]The packet size is too large.")
	}

	endpoint := FormCompleteURL(client.Url, receiveEndpoint)
	var response []byte
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	slog.Info("[LOG] Prepare to send out " + fmt.Sprint(len(data)) + " bytes data to IF:" + endpoint)
	response, _ = SendRequest(
		http.MethodPost,
		endpoint,
		bytes.NewBuffer(data),
		headers,
	)
	var result map[string]interface{}
	json.Unmarshal(response, &result)
	slog.Info(string(response))
}

func SendRequest(operation string, endpoint string, form io.Reader, headers map[string]string) ([]byte, http.Header) {
	newRequest, err := http.NewRequest(
		operation,
		endpoint,
		form,
	)
	if err != nil {
		panic(err)
	}
	for k := range headers {
		newRequest.Header.Add(k, headers[k])
	}
	// Skip certificate verification.
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	var res *http.Response
	for i := 0; i < HTTP_RETRY_TIMES; i++ {
		res, err = client.Do(newRequest)
		if err == nil {
			break // Request successful, exit the loop
		}
		fmt.Printf("Error occurred: %v\n", err)
		time.Sleep(HTTP_RETRY_INTERVAL * time.Second)
		fmt.Printf("Sleep for " + fmt.Sprint(HTTP_RETRY_INTERVAL) + " seconds and retry .....")
	}
	if err != nil {
		slog.Info("[ERROR] HTTP connection failure after 10 times of retry.")
		panic(err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	return body, res.Header
}
