package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	_ "github.com/pbberlin/tools/net/http/proxy1"
)

var (
	otreeRoot = "/otree"
	otreeData = "/otree/data"
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

func paramInt(r *http.Request, key string) int {
	str := param(r, key)
	integer, err := strconv.Atoi(str)
	if err != nil {
		integer = 0
	}
	return integer
}

type FlyingPathT struct {
	OTreeUserId        int    `json:"otree_user_id"`
	OTreeUserIdComment string `json:"otree_user_id_comment"`

	Draws   int `json:"number_of_draws"`
	Periods int `json:"number_of_periods"`

	MinSlow        int    `json:"min_number_slow_sims"`
	MinSlowComment string `json:"min_number_slow_sims_comment"`
	MinFast        int    `json:"min_number_fast_sims"`

	YMin int `json:"y_min"`
	YMax int `json:"y_max"`

	Wealth0         int `json:"wealth0"` // wealth at the start of the simulation
	WealtMultiplier int `json:"wealth_multiplier"`

	DataPointNumbers        bool   `json:"indicator_numbers"` // chart: show numbers next to datapoints
	DataPointNumbersComment string `json:"indicator_numbers_comment"`

	PaymentsWealthComment string  `json:"payments_wealth_comment"`
	Payments              [][]int `json:"payments"` // payment
	Wealth                [][]int `json:"wealth"`   // summed up payments; redundant
}

var FP = FlyingPathT{}

func init() {

	//
	http.HandleFunc(otreeRoot, func(w http.ResponseWriter, r *http.Request) {
		cnt := fmt.Sprintf("<a href='%v'> Get random flying path data in JSON format</a> <br />", otreeData)
		w.Write([]byte(cnt))
		cnt = fmt.Sprintf("<a href='%v?keep=10'>Keep first 10 realizations</a> <br />", otreeData)
		w.Write([]byte(cnt))
		cnt = fmt.Sprintf("<a href='%v?draws=10&periods=10'>10 draws - 10 periods</a> <br />", otreeData)
		w.Write([]byte(cnt))
	})
	http.HandleFunc(otreeData, flyingPathData)

	FP.OTreeUserId = 32168
	FP.OTreeUserIdComment = `Save into global JavaScript variable so that code can be instrumented with it.

			var lnk = document.getElementById("xx");
			lnk.onclick  = function(){
			  console.log("i was clicked");
			};
			
			
			function wrapFunc(func, message) {
				return function () {
					func();
					console.log(message);
				}
			}

			lnk.onclick = wrapFunc(lnk.onclick, "now some log msg");

	`

	FP.Draws = 10000
	FP.Periods = 35
	FP.PaymentsWealthComment = "Wealth is of course redundant. It can be created by summing up the payments."
	FP.initStructSlices()

	FP.MinSlow = 10
	FP.MinSlowComment = "If zero: Start immediately in fast mode. Otherwise indicates appearance of button allowing user to advance to fast mode."
	FP.MinFast = 50

	FP.YMin = -100
	FP.YMax = 100

	FP.Wealth0 = 0
	FP.WealtMultiplier = 10

	FP.DataPointNumbersComment = "Show numbers next to datapoints in Charts"

	FP.DataPointNumbers = false

}

func (f *FlyingPathT) initStructSlices() {

	// We could re-adjust slice lengths for already existing slices, but we are too lazy
	f.Payments = make([][]int, f.Draws, f.Draws)
	f.Wealth = make([][]int, f.Draws, f.Draws)
	for i := 0; i < f.Draws; i++ {
		f.Payments[i] = make([]int, f.Periods, f.Periods)
		f.Wealth[i] = make([]int, f.Periods, f.Periods)
	}

}

func flyingPathData(w http.ResponseWriter, r *http.Request) {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	iKeep := paramInt(r, "keep")

	draws := paramInt(r, "draws")
	periods := paramInt(r, "periods")
	if draws > 0 {
		FP.Draws = draws
	}
	if periods > 0 {
		FP.Periods = periods
	}
	if draws > 0 || periods > 0 {
		FP.initStructSlices()
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
	// w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Write(js)

}
