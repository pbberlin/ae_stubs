package main

import (
	"github.com/pbberlin/dom/clean"
	"github.com/pbberlin/dom/ui"
)

var cf clean.Config = clean.GetDefaultConfig()

func init() {

	opt1 := func(c *clean.Config) { c.HtmlTitle = "Proxify http requests" }
	cf.Apply(opt1, opt1)

	ui.ExplicitInit(nil)

}
