package main

import (
	"log"

	_ "net/http/pprof" // profiling

	// not used - but init() functions wanted for
	// httpHandler registrations
	_ "github.com/pbberlin/tools/backend"
	_ "github.com/pbberlin/tools/big_query"
	_ "github.com/pbberlin/tools/blobstore_mgt"
	_ "github.com/pbberlin/tools/dsu/ancestored_urls"
	_ "github.com/pbberlin/tools/dsu/persistent_cursor"
	_ "github.com/pbberlin/tools/email"
	_ "github.com/pbberlin/tools/foscam"
	_ "github.com/pbberlin/tools/fulltext"
	_ "github.com/pbberlin/tools/guestbook"
	_ "github.com/pbberlin/tools/instance_mgt"
	_ "github.com/pbberlin/tools/json"
	_ "github.com/pbberlin/tools/namespaces_taskqueues"
	_ "github.com/pbberlin/tools/pbfetch"
	_ "github.com/pbberlin/tools/pboauth"
	_ "github.com/pbberlin/tools/write_methods"
)

func init() {

	log.Println("init() for mod01 (alias 'default') complete")
	//util_err.StackTrace(5)

}
