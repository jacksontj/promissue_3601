package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/common/model"
)

var NUM_TIMESERIES = 5000
var NUM_DATAPOINTS = 17280

func generateData(timeseries, datapoints int) model.Matrix {

	// ~177x slower to marshal than to generate

	// Create the top-level matrix
	m := make(model.Matrix, 0)

	for i := 0; i < timeseries; i++ {
		lset := map[model.LabelName]model.LabelValue{
			model.MetricNameLabel: model.LabelValue("timeseries_" + strconv.Itoa(i)),
		}

		now := model.Now()

		values := make([]model.SamplePair, datapoints)

		for x := datapoints; x > 0; x-- {
			values[x-1] = model.SamplePair{
				Timestamp: now.Add(time.Second * -15 * time.Duration(x)), // Set the time back assuming a 15s interval
				Value:     model.SampleValue(float64(x)),
			}
		}

		ss := &model.SampleStream{
			Metric: model.Metric(lset),
			Values: values,
		}

		m = append(m, ss)
	}
	return m
}

func test() {

	start := time.Now()
	m := generateData(NUM_TIMESERIES, NUM_DATAPOINTS)
	took := time.Now().Sub(start)

	fmt.Println("done generatingData took:", took)

	start = time.Now()
	json.Marshal(m)
	took = time.Now().Sub(start)
	fmt.Println("done marshaling took:", took)
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		fmt.Println("response handler took", time.Now().Sub(start))
	}()

	//m := generateData(NUM_TIMESERIES, NUM_DATAPOINTS)
	m := generateData(5, 100)

	w.Header().Set("Content-Type", "application/json")
	if false {
		enc := json.NewEncoder(w)
		w.Write([]byte{'['})

		for i, item := range m {
			if err := enc.Encode(item); err != nil {
				fmt.Println(err)
				return
			}
			if i < len(m)-1 {
				w.Write([]byte{','})
			}
		}
		w.Write([]byte{']'})
	} else {
		b, _ := json.Marshal(m)
		w.Write(b)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
