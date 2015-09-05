package main

import (
	"net/http"

	"github.com/pbberlin/tools/net/http/upload"

	_ "net/http/pprof" // profiling  ... http://.../debug/pprof

	// not used - but init() functions wanted for
	// httpHandler registrations
	_ "github.com/pbberlin/tools/appengine/backend"
	_ "github.com/pbberlin/tools/appengine/blobstore_mgt"
	_ "github.com/pbberlin/tools/appengine/fulltext"
	_ "github.com/pbberlin/tools/appengine/guestbook"
	_ "github.com/pbberlin/tools/appengine/instance_info"
	_ "github.com/pbberlin/tools/appengine/namespaced_taskqueued_cntr"
	_ "github.com/pbberlin/tools/big_query"
	_ "github.com/pbberlin/tools/dsu/ancestored_urls"
	_ "github.com/pbberlin/tools/dsu/persistent_cursor"
	_ "github.com/pbberlin/tools/email"
	_ "github.com/pbberlin/tools/foscam"
	_ "github.com/pbberlin/tools/net/http/proxy1"
	_ "github.com/pbberlin/tools/oauthpb"
	// _ "github.com/pbberlin/tools/os/fsi/dsfs"
	_ "github.com/pbberlin/tools/util" // counter reset
	_ "github.com/pbberlin/tools/write_methods"
)

func init() {
	http.HandleFunc("/backend1", backendHandler)
}

func backendHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I am backend<br>\n<br>\n"))
	uiUpload := upload.BackendUIRendered()
	w.Write(uiUpload.Bytes())
}
