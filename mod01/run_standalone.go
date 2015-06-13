// +build !appengine

// https://blog.golang.org/the-app-engine-sdk-and-workspaces-gopath

package main

import (
	"net/http"

	// not used - but init() functions wanted for
	// httpHandler registrations
	_ "github.com/pbberlin/tools/backend"
	// _ "github.com/pbberlin/tools/big_query"
	// _ "github.com/pbberlin/tools/blobstore_mgt"
	// _ "github.com/pbberlin/tools/dsu_ancestored_urls"
	// _ "github.com/pbberlin/tools/dsu_persistent_cursor"
	// _ "github.com/pbberlin/tools/email"
	// _ "github.com/pbberlin/tools/fetch"
	// _ "github.com/pbberlin/tools/fulltext"
	// _ "github.com/pbberlin/tools/guestbook"
	// _ "github.com/pbberlin/tools/instance_mgt"
	// _ "github.com/pbberlin/tools/json"
	// _ "github.com/pbberlin/tools/namespaces_taskqueues"
)

func main() {
	http.ListenAndServe("localhost:8086", nil)
}
