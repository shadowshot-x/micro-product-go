package monitormodule

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

const endpoint string = "http://localhost:9090/metrics"
const outputEndpoint string = "http://localhost:9090/checkRoutine"
const metricConfFile string = "./monitormodule/config.yaml"

// this will have endpoints to monitor as per config.yaml
type confMap struct {
	Metadataregion     string
	Metadatapipelineid int
	Signin             string
	Signup             string
}

// we need some metadata of the record
// this can be populated using this struct
type transformationOutput struct {
	Time          time.Time
	Region        string
	PipelineId    int
	MetricDetails []string
}

// dataAggregator combines data from the microservices provided
func dataAggregator(log *zap.Logger, trackData confMap, key string, aggregateChan chan<- []string) error {
	// data is fetched every 5 seconds from the microservice specified.
	fetchInterval := time.NewTicker(5 * time.Second)
	// we start an infinite loop as we want to monitor always.
	for {
		// the above fetchInterval gives us a channel which we query
		select {
		case <-fetchInterval.C:
			// we send the get request to prometheus http endpoint
			resp, err := http.Get(endpoint)
			// if there is no error, move ahead, else stop the goroutine immediately
			if err != nil {
				log.Error("Error", zap.Any("message", err))
				return err
			}
			// read the response body from the server
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Error(err.Error())
			}
			// now we pick out the metrics that we want
			// this is based on the conf yaml file
			contents := string(body)
			aggregate := []string{}
			// i have used this to get value dynamically from a struct
			reflection := reflect.ValueOf(trackData)
			reflectionField := reflect.Indirect(reflection).FieldByName(key)
			// get the keys we want to fetch
			allKeys := strings.Split(reflectionField.String(), ",")
			for _, metricValue := range strings.Split(contents, "\n") {
				for _, metricKey := range allKeys {
					// check if the given line has the metric we are looking for
					if strings.Contains(metricValue, metricKey) {
						// ignore the first character
						if metricValue[0] != '#' {
							aggregate = append(aggregate, metricValue)
						}
					}
				}
			}
			// save the result in the aggregateChannel
			aggregateChan <- aggregate
		}
	}
}

// dataTransformer reads from aggregate channel and transorms into a struct. This is output to a unique channel for this goroutine only.
func dataTransformer(log *zap.Logger, aggeregateChan <-chan []string, transformChan chan<- string, routineId int, region string, pipelineid int) error {
	// read from the channel
	for val := range aggeregateChan {
		tr := transformationOutput{
			Time:          time.Now(),
			Region:        region,
			PipelineId:    pipelineid,
			MetricDetails: val,
		}
		// marshal into byte array
		op, err := json.Marshal(tr)
		if err != nil {
			log.Error("Error Encountered", zap.Int("routine id", routineId), zap.Any("error", err))
		}
		// send the output to the respective channel
		transformChan <- string(op)
	}
	return nil
}

// dataTransportation send a POST request to our endpoint which is on our server only
func dataTransportation(log *zap.Logger, transportChan <-chan string, successChan chan<- string) error {
	for op := range transportChan {
		// this is fanned in channel
		// we send a POST request with our metric data
		response, err := http.Post(outputEndpoint, "application/json", bytes.NewBuffer([]byte(op)))
		if err != nil {
			log.Error("Error in transportation", zap.Error(err))
		}
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Error("Error in Reading Output body", zap.Error(err))
			return err
		}
		// this is a temp channel to show how waitgroups function
		successChan <- string(responseBody)
	}
	return nil
}

// simple function that prints the channels value available. It takes in the Waitgroup
func showDemo(log *zap.Logger, wg *sync.WaitGroup, finalChan <-chan string) {
	log.Info("final Chan ", zap.String("message", <-finalChan))
	wg.Done()
}

// MonitorBinder binds all the goroutines and combines the output
func MonitorBinder(log *zap.Logger) error {
	// read the yaml file and unmarshal to the metricMap
	content, err := ioutil.ReadFile(metricConfFile)
	if err != nil {
		return err
	}
	// we read the info from config yaml file
	metricEndpoints := confMap{}
	yaml.Unmarshal(content, &metricEndpoints)
	log.Info("Contents", zap.Any("Any", metricEndpoints))

	// setting up unbuffered channel
	aggregateChan := make(chan []string)
	// Go routintes to monitor metrcis for 2 sevices
	// there are signin and signup
	go dataAggregator(log, metricEndpoints, "Signin", aggregateChan)
	go dataAggregator(log, metricEndpoints, "Signup", aggregateChan)

	// these are the output channel for transormer function
	transform1Chan := make(chan string)
	transform2Chan := make(chan string)

	// these 2 goroutines transform the data and send output to transform1Chan and transform2Chan
	go dataTransformer(log, aggregateChan, transform1Chan, 1, metricEndpoints.Metadataregion, metricEndpoints.Metadatapipelineid)
	go dataTransformer(log, aggregateChan, transform2Chan, 2, metricEndpoints.Metadataregion, metricEndpoints.Metadatapipelineid)

	transportChan := make(chan string)
	// now we need fanin of 2 channels to a single channel. This final channel will be passed for transportation
	// we can write another goroutine for this
	go func() {
		for {
			select {
			case op1 := <-transform1Chan:
				transportChan <- op1
			case op2 := <-transform2Chan:
				transportChan <- op2
			case <-time.After(3 * time.Second):
				// the last statement imposes a timeout that maybe needed if the above routines are taking time to execute
			}
		}
	}()
	// this has te output from our server
	successChan := make(chan string)

	// spin 2 gorotines in parallel to get the successChan fast
	go dataTransportation(log, transportChan, successChan)
	go dataTransportation(log, transportChan, successChan)

	// now we implement a WaitGroup to get the first 5 values only of the response body
	// this will
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		// add to the waitgroup. Acts like a semaphore
		wg.Add(1)
		go showDemo(log, &wg, successChan)
	}
	// this goroutine only runs after the first 5 responses have been printed
	go func() {
		wg.Wait()
		log.Info("=========== WE ARE DONE FOR THE DEMO ===========")
	}()
	return nil
}
