package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"appengine"

	"github.com/pbberlin/tools/net/http/coinbase"
	"github.com/pbberlin/tools/net/http/fileserver"
	"github.com/pbberlin/tools/net/http/htmlfrag"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/tplx"
	"github.com/pbberlin/tools/net/http/upload" // upload receive
	"github.com/pbberlin/tools/oauthpb"
	"github.com/pbberlin/tools/os/fsi/dsfs"
	"github.com/pbberlin/tools/os/fsi/memfs"
	"github.com/pbberlin/tools/os/fsi/webapi"
)

var fs1 = memfs.New(
	memfs.Ident("mnt02"), // a closured variable in init() did not survive map-pointer reallocation
)

func init() {

	upload.InitHandlers()
	coinbase.InitHandlers()
	http.HandleFunc(webapi.UriDeleteSubtree, loghttp.Adapter(webapi.DeleteSubtree))

	http.HandleFunc("/backend-secret", backendHandler)

	dynSrv := func(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

		if strings.Contains(r.URL.Path, "/member/") {
			auth, msg := oauthpb.Auth(r)
			if auth == false {
				w.Write([]byte(msg))
				return
			}
		}

		c := appengine.NewContext(r)
		appID := appengine.AppID(c)
		if appID == "tec-news" {

			prefix := "/mnt02"
			// prefix = "/xxx"

			fs2 := dsfs.New(
				dsfs.MountName(prefix[1:]),
				dsfs.AeContext(appengine.NewContext(r)),
			)

			fs1.SetOption(
				memfs.ShadowFS(fs2),
			)

			//
			// TRICK
			// making FsiFileServer dream, that the mount prefix was
			r.URL.Path = prefix + r.URL.Path
			fileserver.FsiFileServer(fs1, prefix, w, r)
		} else {
			w.Write([]byte("app id is -" + appID + "- "))
		}
	}
	http.HandleFunc("/", loghttp.Adapter(dynSrv))

	//
	dmpMemfs := func(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
		htmlfrag.SetNocacheHeaders(w)
		r.Header.Set("X-Custom-Header-Counter", "nocounter")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte("<pre>"))
		w.Write(fs1.Dump())
	}
	http.HandleFunc("/memfsdmp", loghttp.Adapter(dmpMemfs))

}

var wpf = fmt.Fprint

func backendHandler(w http.ResponseWriter, r *http.Request) {

	lg, b := loghttp.BuffLoggerUniversal(w, r)
	_ = lg
	closureOverBuf := func(bUnused *bytes.Buffer) {
		loghttp.Pf(w, r, b.String())
	}
	defer closureOverBuf(b) // the argument is ignored,

	r.Header.Set("X-Custom-Header-Counter", "nocounter")

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]string{"HtmlTitle": "Static uploading and file serving"}))
	defer wpf(w, tplx.Foot)

	htmlfrag.Wb(w, "secret backend", "")
	htmlfrag.Wb(w, "to root", "/", " ")

	wpf(w, upload.BackendUIRendered().String())

	htmlfrag.Wb(w, "fsi tools", "")
	htmlfrag.Wb(w, "remove subtr", webapi.UriDeleteSubtree, " ")
	htmlfrag.Wb(w, "memfs dump", "/memfsdmp", " ")

	wpf(w, coinbase.BackendUIRendered().String())

}
