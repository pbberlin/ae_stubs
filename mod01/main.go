package main

import (
	"log"

	_ "net/http/pprof" // profiling

	// not used - but init() functions wanted for
	// httpHandler registrations
	_ "github.com/pbberlin/tools/appengine/backend"
	_ "github.com/pbberlin/tools/appengine/blobstore_mgt"
	_ "github.com/pbberlin/tools/appengine/fulltext"
	_ "github.com/pbberlin/tools/appengine/guestbook"
	_ "github.com/pbberlin/tools/appengine/instance_mgt"
	_ "github.com/pbberlin/tools/appengine/namespaces_taskqueues"
	_ "github.com/pbberlin/tools/big_query"
	_ "github.com/pbberlin/tools/dsu/ancestored_urls"
	_ "github.com/pbberlin/tools/dsu/persistent_cursor"
	_ "github.com/pbberlin/tools/email"
	_ "github.com/pbberlin/tools/foscam"
	_ "github.com/pbberlin/tools/godoc/vfs/gaefs"
	_ "github.com/pbberlin/tools/json"
	_ "github.com/pbberlin/tools/net/http/proxy1"
	_ "github.com/pbberlin/tools/oauthpb"
	_ "github.com/pbberlin/tools/util" // counter reset
	_ "github.com/pbberlin/tools/write_methods"
)

func init() {
	log.Println("init() for mod01 (alias 'default') complete")
	//util_err.StackTrace(5)
}
