package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"runtime/pprof"
	"strings"
	"testing"

	"github.com/json-iterator/go"
	"github.com/mailru/easyjson"
	"github.com/pquerna/ffjson/ffjson"
)

func BenchmarkMarshal(b *testing.B) {

	NUM_TIMESERIES := 500
	NUM_DATAPOINTS := 100
	m := generateData(NUM_TIMESERIES, NUM_DATAPOINTS)

	// Single float64
	b.Run("SampleValue", func(b *testing.B) {
		marshalSomething(b, m[0].Values[0].Value)
	})

	// Single timestamp
	b.Run("Timestamp", func(b *testing.B) {
		marshalSomething(b, m[0].Values[0].Timestamp)
	})

	// Single pair of value and timestamp
	b.Run("SamplePair", func(b *testing.B) {
		marshalSomething(b, m[0].Values[0])
	})

	// Labelset for the metrics
	b.Run("labelset", func(b *testing.B) {
		marshalSomething(b, m[0].Metric)
	})

	// labelset + []SamplePair
	b.Run("SampleStream", func(b *testing.B) {
		marshalSomething(b, m[0])
	})

	// labelset + []SamplePair
	b.Run("Matrix", func(b *testing.B) {
		marshalSomething(b, m)
	})
}

var marshaled []byte
var str string

func marshalSomething(b *testing.B, v interface{}) {
	b.Run("encoding/json", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			marshaled, _ = json.Marshal(v)
		}
	})

	if false {
		b.Run("ffjson", func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				marshaled, _ = ffjson.Marshal(v)
			}
		})

		b.Run("jsoniter", func(b *testing.B) {
			var j = jsoniter.ConfigCompatibleWithStandardLibrary
			for n := 0; n < b.N; n++ {
				marshaled, _ = j.Marshal(v)
			}
		})
	}

	b.Run("easyjson", func(b *testing.B) {
		t := reflect.TypeOf(v)
		m, ok := v.(easyjson.Marshaler)
		if !ok {
			fmt.Println(t, "not easyjson.Marshaler")
			b.Skip()
		}
		
		tName := t.Name()
		if t.Kind() == reflect.Ptr {
		    tName = t.Elem().Name()
		}

        fname := "/tmp/prof/" + strings.Replace(tName, "*", "_", -1)
		f, err := os.Create(fname)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
		
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			easyjson.MarshalToWriter(m, ioutil.Discard)
		}
	})
}
