package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"appengine"

	"github.com/pbberlin/tools/dsu"
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
	memfs.Ident(tplx.TplPrefix[1:]), // a closured variable in init() did not survive map-pointer reallocation
)

func init() {

	upload.InitHandlers()
	coinbase.InitHandlers()
	http.HandleFunc(webapi.UriDeleteSubtree, loghttp.Adapter(webapi.DeleteSubtree))

	http.HandleFunc("/backend-secret", backendHandler)

	dynSrv := func(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

		lg, b := loghttp.BuffLoggerUniversal(w, r)
		_ = b

		r.Header.Set("X-Custom-Header-Counter", "nocounter")
		htmlfrag.SetNocacheHeaders(w)

		if !strings.Contains(r.URL.Path, "/member/") {
			serveFromRoot(w, r)

		} else {

			auth, usr, msg := oauthpb.Auth(r)
			if msg != "" {
				msg += "<br>"
			}

			usrID := "32168-unknown-user" // prevent method call on nil further down
			if usr != nil {
				usrID = usr.ID
			}

			if auth == false {

				w.Write([]byte(msg))
				return

			} else {

				invoice, err := dsu.BufGet(appengine.NewContext(r), "dsu.WrapBlob__"+usrID+r.URL.Path)
				lg(err)
				buyStatus := ""
				fullJSONData := ""
				if err != nil {
					buyStatus = "You have not bought this article yet.<br>"
				} else {

					tm := time.Unix(int64(invoice.I), 0)
					buyStatus = fmt.Sprintf("status %v - UID %v Amount %v - at %v<br>",
						invoice.Desc, invoice.Name, invoice.F, tm)
					// fullJSONData = "<pre>" + string(invoice.VByte) + "</pre>"

					if invoice.Desc == "completed" {
						serveFromRoot(w, r)
						return
					}

				}

				//

				btnTest := `
					<div style='height:10px;'>&nbsp;</div>
					<a class="coinbase-button" 
						data-code="0025d69ea925b48ba2b7adeb2a911ca2" 
						data-custom="productID=` + r.URL.Path + `&uID=` + usrID + `" 
						data-env="sandbox" 
						href="https://sandbox.coinbase.com/checkouts/0025d69ea925b48ba2b7adeb2a911ca2" 
					>Pay With Bitcoin</a>
					<script src="https://sandbox.coinbase.com/assets/button.js" type="text/javascript"></script>				
				`

				btnLive := `
					<div style='height:10px;'>&nbsp;</div>
					<a class="coinbase-button" 
						data-code="aa4e03abbc5e2f5321d27df32756a932" 
						data-custom="productID=` + r.URL.Path + `&uID=` + usrID + `" 
						href="https://www.coinbase.com/checkouts/aa4e03abbc5e2f5321d27df32756a932" 
					>Pay With Bitcoin</a>
					<script src="https://www.coinbase.com/assets/button.js" type="text/javascript"></script>

				`

				_, _ = btnLive, btnTest

				backPath := strings.Replace(r.URL.Path, "/member", "", 1)
				backAnch := fmt.Sprintf("<a href='%v'>Back to introduction</a><br>", backPath)

				bstpl := tplx.BootstrapTemplate(w, r)
				wpf(w,
					tplx.ExecTplHelper(bstpl, map[string]interface{}{
						"HtmlTitle":       "Access restricted",
						"HtmlDescription": "", // reminder
						"HtmlContent": template.HTML("Access is restricted<br>" +
							msg +
							btnLive + "<br>" +
							backAnch +
							buyStatus +
							fullJSONData +
							"<br>")}))

				return
			}
		}

	}
	http.HandleFunc("/", loghttp.Adapter(dynSrv))

	//
	dmpMemfs := func(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
		htmlfrag.SetNocacheHeaders(w)
		r.Header.Set("X-Custom-Header-Counter", "nocounter")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte("<pre>"))

		fs2 := dsfs.New(
			dsfs.MountName(tplx.TplPrefix[1:]),
			dsfs.AeContext(appengine.NewContext(r)),
		)
		fs1.SetOption(
			memfs.ShadowFS(fs2),
		)

		w.Write(fs1.Dump())
	}
	http.HandleFunc("/dump-memfs", loghttp.Adapter(dmpMemfs))

	resetMemfs := func(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
		fs1 = memfs.New(
			memfs.Ident(tplx.TplPrefix[1:]),
		)
	}
	http.HandleFunc("/reset-memfs", loghttp.Adapter(resetMemfs))
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

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]interface{}{"HtmlTitle": "Static uploading and file serving"}))
	defer wpf(w, tplx.Foot)

	htmlfrag.Wb(w, "secret backend", "")
	htmlfrag.Wb(w, "to root", "/", " ")

	wpf(w, upload.BackendUIRendered().String())

	htmlfrag.Wb(w, "fsi tools", "")
	htmlfrag.Wb(w, "remove subtr", webapi.UriDeleteSubtree, " ")
	htmlfrag.Wb(w, "memfs dump", "/dump-memfs", " ")
	htmlfrag.Wb(w, "memfs reset", "/reset-memfs", " ")

	wpf(w, coinbase.BackendUIRendered().String())

}

func serveFromRoot(w http.ResponseWriter, r *http.Request) {

	appID := appengine.AppID(appengine.NewContext(r))
	if appID == "tec-news" {

		fs2 := dsfs.New(
			dsfs.MountName(tplx.TplPrefix[1:]),
			dsfs.AeContext(appengine.NewContext(r)),
		)
		fs1.SetOption(
			memfs.ShadowFS(fs2),
		)

		//
		// TRICK
		// Making FsiFileServer dream, that the request path contained the mount prefix
		r.URL.Path = tplx.TplPrefix + r.URL.Path
		fileserver.FsiFileServer(fs1, tplx.TplPrefix, w, r)
	} else {
		w.Write([]byte("wrong app id: " + appID + "- "))
	}

}
