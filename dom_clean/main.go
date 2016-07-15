package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var (
	dataUrl = "/otree/data"
)

func param(r *http.Request, key string) string {
	values := r.URL.Query()
	for lpkey, sl := range values {
		if lpkey == key {
			last := len(sl)
			return sl[last-1]
		}
	}
	return ""
}

type FlyingPathT struct {
	Draws   int `json:"number_of_draws"`
	Periods int `json:"number_of_periods"`

	MinSlow int `json:"min_number_slow_sims"`
	MinFast int `json:"min_number_fast_sims"`

	YMin int `json:"y_min"`
	YMax int `json:"y_max"`

	WealtMultiplier int `json:"wealth_multiplier"`

	Wealth0 int `json:"wealth0"` // wealth at the start of the simulation

	DataPointNumbers bool `json:"indicator_numbers"` // chart: show numbers next to datapoints

	Payments [][]int `json:"payments"` // payment
	Wealth   [][]int `json:"wealth"`   // summed up payments; redundant
}

var FP = FlyingPathT{}

func init() {

	FP.Draws = 10000
	FP.Periods = 35

	FP.MinSlow = 10
	FP.MinFast = 50

	FP.YMin = -100
	FP.YMax = 100

	FP.Wealth0 = 0
	FP.WealtMultiplier = 10

	FP.DataPointNumbers = false

	FP.Payments = make([][]int, FP.Draws, FP.Draws)
	FP.Wealth = make([][]int, FP.Draws, FP.Draws)

	for i := 0; i < FP.Draws; i++ {
		FP.Payments[i] = make([]int, FP.Periods, FP.Periods)
		FP.Wealth[i] = make([]int, FP.Periods, FP.Periods)
	}

}

func flyingPathData(w http.ResponseWriter, r *http.Request) {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	keep := param(r, "keep")
	iKeep, err := strconv.Atoi(keep)
	if err != nil {
		iKeep = 0
	}

	for i := 0; i < FP.Draws; i++ {
		for j := iKeep; j < FP.Periods; j++ {
			draw := r1.Intn(2)
			if draw == 0 {
				draw = -1
			}
			FP.Payments[i][j] = draw * FP.WealtMultiplier
			if j > 0 {
				FP.Wealth[i][j] = FP.Wealth[i][j-1] + FP.Payments[i][j]
			} else {
				FP.Wealth[i][j] = FP.Wealth0 + FP.Payments[i][j]
			}
		}
	}

	// str := util.IndentedDump(FP)
	// log.Printf("%v", str)

	js, err := json.MarshalIndent(FP, "", "\t")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	// log.Printf("%v", js)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func init() {

	//
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cnt := fmt.Sprintf("<a href='%v'> Get random flying path data in JSON format</a> <br />", dataUrl)
		w.Write([]byte(cnt))
		cnt = fmt.Sprintf("<a href='%v?keep=10'>Keep first 10 realizations</a> <br />", dataUrl)
		w.Write([]byte(cnt))
	})
	http.HandleFunc("/xxx", flyingPathData)
	http.HandleFunc(dataUrl, flyingPathData)

}
