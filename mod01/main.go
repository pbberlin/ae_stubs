package main

import (
	"net/http"

	_ "net/http/pprof"

	"bytes"
	"fmt"
	"io"
	"log"

	// not used - but init() functions wanted for
	// httpHandler registrations
	_ "github.com/pbberlin/tools/big_query"
	_ "github.com/pbberlin/tools/blobstore_mgt"
	_ "github.com/pbberlin/tools/dsu_ancestored_urls"
	_ "github.com/pbberlin/tools/dsu_persistent_cursor"
	_ "github.com/pbberlin/tools/email"
	_ "github.com/pbberlin/tools/fetch"
	_ "github.com/pbberlin/tools/fulltext"
	_ "github.com/pbberlin/tools/guestbook"
	_ "github.com/pbberlin/tools/instance_mgt"
	_ "github.com/pbberlin/tools/json"
	_ "github.com/pbberlin/tools/namespaces_taskqueues"
)

var sq func(a ...interface{}) string = fmt.Sprint
var sp func(format string, a ...interface{}) string = fmt.Sprintf
var fp func(w io.Writer, format string, a ...interface{}) (int, error) = fmt.Fprintf

// small helper
func wb(buf1 *bytes.Buffer, linktext, url string) {

	if url == "" {
		buf1.WriteString("<br>\n")
	}

	buf1.WriteString("<span style='display:inline-block; min-width:200px; margin: 6px 0px; margin-right:10px;'>\n")
	if url == "" {
		buf1.WriteString("\t" + linktext + "\n")
	} else {
		buf1.WriteString("\t<a target='_app' href='" + url + "' >" + linktext + "</a>\n")
	}
	buf1.WriteString("</span>\n")
}

func init() {

	http.HandleFunc("/login", login)
	//http.HandleFunc("/", util_err.Adapter(big_query.ViewHTML))
	log.Println("init() for mod01 (alias 'default') complete")
	//util_err.StackTrace(5)

}
