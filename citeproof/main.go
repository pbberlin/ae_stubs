package main

import (
	"net/http"

	"github.com/zew/exceldb/dom/clean"
	"github.com/zew/exceldb/dom/ui"
)

var cf clean.Config = clean.GetDefaultConfig()

func init() {

	opt1 := func(c *clean.Config) { c.HtmlTitle = "Proxify http requests" }
	cf.Apply(opt1, opt1)

	ui.ExplicitInit()

	http.HandleFunc("/put", put)
	http.HandleFunc("/query-pages", queryPages)

}
