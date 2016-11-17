package main

import (
	"net/http"
	"time"

	"github.com/zew/logx"
	"github.com/zew/util"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type HtmlPage struct {
	Val      int
	Url      string
	Unixtime int64
	T        time.Time
	Body     string
}

const htmlPageKind = "HtmlPage"

func queryPages(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)

	q := datastore.NewQuery(htmlPageKind).Filter("Url >=", "").Order("Url").Order("Unixtime").Limit(5)
	var pages []HtmlPage
	_, err := q.GetAll(ctx, &pages)
	util.CheckErr(err)

	logx.Debugf(r, "found %v pages", len(pages))
	for _, p := range pages {
		logx.Debugf(r, "%s %s", p.Url, p.Unixtime)
	}
}

func put(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)

	pg := &HtmlPage{
		Val:      32168,
		Url:      "faz.net/aktuell/second.html",
		Unixtime: time.Now().Unix(),
		T:        time.Now(),
		Body:     "body",
	}
	_ = pg

	key := datastore.NewIncompleteKey(ctx, htmlPageKind, nil)
	key = datastore.NewKey(ctx, htmlPageKind, "", 3, nil)
	keyComplete, err := datastore.Put(ctx, key, pg)
	util.CheckErr(err)

	logx.Debugf(r, "keyComplete is %v", keyComplete)
	w.Write([]byte("put complete"))

}
