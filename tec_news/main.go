package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"

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

		if strings.Contains(r.URL.Path, "/member/") {
			auth, usr, msg := oauthpb.Auth(r)
			if msg != "" {
				msg += "<br>"
			}

			if auth == false || true {
				r.Header.Set("X-Custom-Header-Counter", "nocounter")
				htmlfrag.SetNocacheHeaders(w)
				bstpl := tplx.BootstrapTemplate(w, r)

				usrID := "32168-unknown-user"
				if usr != nil {
					usrID = usr.ID
				}

				btnLive := `
					<div style='height:10px;'>&nbsp;</div>
					<a class="coinbase-button" data-code="aa4e03abbc5e2f5321d27df32756a932" 
						data-custom="productID=` + r.URL.Path + `&uID=` + usrID + `" 
						href="https://www.coinbase.com/checkouts/aa4e03abbc5e2f5321d27df32756a932" 
					>Pay With Bitcoin</a>
					<script src="https://www.coinbase.com/assets/button.js" type="text/javascript"></script>

				`
				btnTest := `
					<div style='height:10px;'>&nbsp;</div>
					<a class="coinbase-button" 
						data-code="0025d69ea925b48ba2b7adeb2a911ca2" 
						data-custom="productID=` + r.URL.Path + `&uID=` + usrID + `" 
						data-env="sandbox" 
						href="https://sandbox.coinbase.com/checkouts/0025d69ea925b48ba2b7adeb2a911ca2" 
					>Pay With Bitcoin</a>
					<script src="https://sandbox.coinbase.com/assets/button.js" type="text/javascript"></script>				`

				_, _ = btnLive, btnTest

				backPath := strings.Replace(r.URL.Path, "/member", "", 1)
				backAnch := fmt.Sprintf("<a href='%v'>Back to introduction</a><br>", backPath)

				retrieveAgain, err := dsu.BufGet(appengine.NewContext(r), "dsu.WrapBlob__"+usr.ID)
				lg(err)
				buyStatus := ""
				fullJSONData := ""
				if err != nil {
					buyStatus = "You have not bought this article yet.<br>"
				} else {
					buyStatus = fmt.Sprintf("status %v - UID %v Amount %v<br>",
						retrieveAgain.Desc, retrieveAgain.Name, retrieveAgain.F)
					// fullJSONData = "<pre>" + string(retrieveAgain.VByte) + "</pre>"
				}

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

		c := appengine.NewContext(r)
		appID := appengine.AppID(c)
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
