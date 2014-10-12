package main

import (
	"appengine"
	sc "github.com/pbberlin/tools/dsu_distributed_unancestored"
	"github.com/pbberlin/tools/tpl_html"
	"github.com/pbberlin/tools/util_err"
	"net/http"
)

type b0 struct {
	NumXSectors int
	NumYSectors int
	VB1         []interface{}
	NumB1       int
	NumB2       int
}

type b1 struct {
	Heading string
	VB2     []interface{}
}

type b2 struct {
	Linktext string
	Url      string
	Target   string
}

var myB0 = b0{NumYSectors: 2, VB1: x}

var x = []interface{}{
	b1{
		Heading: "Diverse",
		VB2: []interface{}{
			b2{Linktext: "Login", Url: "/login"},
			b2{Linktext: "Schreib-Methoden", Url: "/write-methods"},
			b2{Linktext: "Letzte Email", Url: "/email-view"},
			b2{Linktext: "Blob List", Url: "/blob/list"},
			b2{Linktext: "Template Demo 1", Url: "/tpl/demo1"},
			b2{Linktext: "Template Demo 2", Url: "/tpl/demo2"},
			b2{Linktext: "Http fetch", Url: "/fetch-url"},
			b2{Linktext: "Instance Info", Url: "/instance-info/view"},
			b2{Linktext: "Gob encode decode", Url: "/big-query/test-gob-codec"},

			b2{Linktext: "JSON encode", Url: "/json-encode"},
			b2{Linktext: "JSON decode", Url: "/json-decode"},

			b2{Linktext: "Fulltext put", Url: "/fulltext-search/put"},
			b2{Linktext: "Fulltext get", Url: "/fulltext-search/get"},
		},
	},

	b1{
		Heading: "Guestbook",
		VB2: []interface{}{
			b2{Linktext: "Einträge auflisten", Url: "/guest-view"},
			b2{Linktext: "Einträge auflisten - paged - serialized cursor", Url: "/guest-view-cursor"},
		},
	},
	b1{
		Heading: "Drawing",
		VB2: []interface{}{
			b2{Linktext: "Drawing a static chart", Url: "/image/draw-lines-example"},
		},
	},
	b1{
		Heading: "Big Query",
		VB2: []interface{}{
			b2{Linktext: "Get real data", Url: "/big-query/query-into-datastore"},
			b2{Linktext: "Get mocked data", Url: "/big-query/mock-data-into-datastore"},
		},
	},
	b1{
		Heading: "... with Chart",
		VB2: []interface{}{
			b2{Linktext: "Process Data 1 (mock=1},", Url: "/big-query/regroup-data-01?mock=0"},
			b2{Linktext: "Process Data 2", Url: "/big-query/regroup-data-02?f=table"},
			b2{Linktext: "Show as Table", Url: "/big-query/show-table"},
			b2{Linktext: "Show as Chart", Url: "/big-query/show-chart"},
			b2{Linktext: "As HTML", Url: "/big-query/html"},
		},
	},
	b1{
		Heading: "Request Images",
		VB2: []interface{}{
			b2{Linktext: "WrapBlob from Datastore", Url: "/image/img-from-datastore?p=chart1"},
			b2{Linktext: "base64 from Datastore", Url: "/image/base64-from-datastore?p=chart1"},
			b2{Linktext: "base64 from Variable", Url: "/image/base64-from-var?p=1"},
			b2{Linktext: "base64 from File", Url: "/image/base64-from-file?p=static/pberg1.png"},
		},
	},
	b1{
		Heading: "Namespaces and Task Queues",
		VB2: []interface{}{
			b2{Linktext: "Increment", Url: "/namespaced-counters/increment"},
			b2{Linktext: "Read", Url: "/namespaced-counters/read"},
			b2{Linktext: "Push to task-queue", Url: "/namespaced-counters/queue-push"},
		},
	},
	b1{
		Heading: "URLs with/without Ancestors",
		VB2: []interface{}{
			b2{Linktext: "Backend", Url: "/save-url/backend"},
		},
	},
	b1{
		Heading: "x",
		VB2:     []interface{}{},
	},
}

func backend3(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	c := appengine.NewContext(r)

	myB0.NumB1 = len(myB0.VB1)
	for _, v := range myB0.VB1 {

		myB1, ok := v.(b1)
		util_err.Err_http(w, r, ok, false)

		myB0.NumB2 += len(myB1.VB2)
	}

	path := m["dir"].(string) + m["base"].(string)

	err := sc.Increment(c, path)
	util_err.Err_http(w, r, err, false)

	cntr, err := sc.Count(w, r, path)
	util_err.Err_http(w, r, err, false)

	add, tplExec := tpl_html.FuncTplBuilder(w, r)
	add("n_html_title", "Backend", nil)

	add("n_cont_0", "<pre>{{.}}</pre>", myB0.NumB2)
	add("n_cont_1", tpl_html.PrefixLff+"backend3_body", myB0)
	add("tpl_legend", tpl_html.PrefixLff+"backend3_body_embed01", "")
	add("n_cont_2", "<p>{{.}} views</p>", cntr)

	tplExec(w, r)

}

func prepareb0(l b0) {

}

func init() {
	prepareb0(myB0)
	http.HandleFunc("/backend3", util_err.Adapter(backend3))
}
