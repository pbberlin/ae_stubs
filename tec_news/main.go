package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"appengine"

	"github.com/pbberlin/tools/appengine/login"
	"github.com/pbberlin/tools/appengine/login/gitkit1"
	"github.com/pbberlin/tools/appengine/memcachepb"
	"github.com/pbberlin/tools/dsu"
	"github.com/pbberlin/tools/net/http/coinbase"
	"github.com/pbberlin/tools/net/http/fileserver"
	"github.com/pbberlin/tools/net/http/htmlfrag"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/tplx"
	"github.com/pbberlin/tools/net/http/upload" // upload receive
	"github.com/pbberlin/tools/os/fsi/dsfs"
	"github.com/pbberlin/tools/os/fsi/memfs"
	"github.com/pbberlin/tools/os/fsi/webapi"
)

var fs1 = memfs.New(
	memfs.Ident(tplx.TplPrefix[1:]), // a closured variable in init() did not survive map-pointer reallocation
)

var (
	successLandingURL = "/auth/signin-landing"
	signoutLandingURL = "/auth/signout-landing"
)

func init() {

	upload.InitHandlers()
	coinbase.InitHandlers()
	tplx.InitHandlers()
	login.InitHandlers()
	gitkit1.InitHandlers()
	memcachepb.InitHandlers()
	http.HandleFunc(webapi.UriDeleteSubtree, loghttp.Adapter(webapi.DeleteSubtree))

	http.HandleFunc("/backend-reduced", backendHandler)

	dynSrv := func(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

		lg, b := loghttp.BuffLoggerUniversal(w, r)
		_ = b

		r.Header.Set("X-Custom-Header-Counter", "nocounter")
		htmlfrag.SetNocacheHeaders(w)

		err := r.ParseForm()
		lg(err)

		if !strings.Contains(r.URL.Path, "/member/") {
			serveFromRoot(w, r, map[string][]byte{})
			return
		}

		//
		// ok, usr, msg := login.CheckForNormalUser(r)
		// if msg != "" {
		// 	msg += "<br>"
		// }
		// if !ok {
		// 	w.Write([]byte(msg))
		// 	return
		// }

		usr := gitkit1.CurrentUser(r)
		if ok := gitkit1.IsSignedIn(r); !ok {
			usr = nil
		}
		if usr == nil {
			http.Redirect(w, r, gitkit1.WidgetSigninAuthorizedRedirectURL+"?mode=select&user=wasNil&red="+r.URL.Path, http.StatusFound)
		}

		//
		//
		usrID := "32168-unknown-user" // prevent method call on nil further down
		if usr != nil {
			usrID = usr.ID
		}

		someUserGreeting := new(bytes.Buffer)
		wpf(someUserGreeting, tplx.ExecTplHelper(gitkit1.UserInfoHTML, map[string]interface{}{
			"User": usr,
		}))

		invoice, err := dsu.BufGet(appengine.NewContext(r), "dsu.WrapBlob__"+usrID+r.URL.Path)
		lg(err)

		origTransactionMap := ""
		if len(invoice.VByte) > 0 && r.Form.Get("transaction-details") != "" {
			origTransactionMap = "<br>\n<pre>" + string(invoice.VByte) + "</pre>"
		}

		buyStatus := ""
		if err != nil || invoice.I == 0 ||
			invoice.Name == "" || invoice.Desc != "completed" {
			buyStatus += "For this article we request a small contribution.<br>"
			buyStatus += "You have not bought this article yet.<br>"
		} else {

			if r.Form.Get("noredirect") == "" {

				idc := "<br>" + gitkit1.GetIDCardTpl(w, r, usr)

				tm := time.Unix(int64(invoice.I), 0).Format("02 Jan 2006 at 15:04")

				Satoshi := int64(invoice.F * (1000 * 1000 * 100))

				// buyStatus += fmt.Sprintf("You bought this article %v for %2.8f BTC.<br>", tm, invoice.F)
				buyStatus += "<div style='font-size:80%;line-height:135%;'>"
				buyStatus += fmt.Sprintf("You bought this article on %v for %v Satoshi.<br>", tm, Satoshi)
				buyStatus += "Thanks for your business. Read your purchased article below."
				buyStatus += "</div>"
				buyStatus += "<hr>"

				backPath := strings.Replace(r.URL.Path, "/member", "", 1)
				backAnch := fmt.Sprintf("<br><a href='%v'>Back to introduction</a>.<br>", backPath)

				serveFromRoot(w, r, map[string][]byte{
					"</head>": []byte(gitkit1.Headers + "\n</head>"),
					"<span id='REPLACE_BEFORE_CONTENT'></span>": []byte(idc + buyStatus + origTransactionMap),
					"<span id='REPLACE_AFTER_CONTENT'></span>":  []byte(backAnch),
					// "<span id='REPLACE_BEFORE_CONTENT'></span>": []byte(buyStatus),
				})
				return
			}

		}

		btnHTML := fmt.Sprintf(coinbase.BtnLiveFormat, r.URL.Path, usrID)
		if r.Form.Get("dbg") > "05" {
			// coded hint, that we are in testing mode
			btnHTML = fmt.Sprintf(coinbase.BtnTestFormat, r.URL.Path, usrID)
			btnHTML += "<div style='display:inline-block;font-size:12px;margin-top:-20px;vertical-align:middle'>coinbase</div>"
		}

		backPath := strings.Replace(r.URL.Path, "/member", "/transition", 1)
		backAnch := fmt.Sprintf("<a href='%v'>Back to introduction</a><br>", backPath)

		bstpl := tplx.TemplateFromHugoPage(w, r)

		wpf(w, tplx.ExecTplHelper(bstpl, map[string]interface{}{
			"HtmlTitle":       "Access restricted",
			"HtmlDescription": "", // reminder
			"HtmlHeaders":     template.HTML(gitkit1.Headers),
			"HtmlContent": template.HTML(
				"<br>" +
					gitkit1.GetIDCardTpl(w, r, usr) + "<br>" +
					buyStatus +
					btnHTML + "<br>" + "<br>" +
					someUserGreeting.String() + "<br>" +
					backAnch +
					origTransactionMap +
					"<br>"),
		},
		),
		)

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

var wpf = fmt.Fprintf

func backendHandler(w http.ResponseWriter, r *http.Request) {

	lg, b := loghttp.BuffLoggerUniversal(w, r)
	_ = lg
	closureOverBuf := func(bUnused *bytes.Buffer) {
		loghttp.Pf(w, r, b.String())
	}
	defer closureOverBuf(b) // the argument is ignored,

	r.Header.Set("X-Custom-Header-Counter", "nocounter")

	if ok, _, msg := login.CheckForAdminUser(r); !ok {
		w.Write([]byte(msg))
		return
	}

	wpf(w, tplx.ExecTplHelper(tplx.Head,
		map[string]interface{}{
			"HtmlTitle": "Static uploading and file serving",
		}),
	)
	defer wpf(w, tplx.Foot)

	htmlfrag.Wb(w, "secret backend", "")
	htmlfrag.Wb(w, "to root", "/", " ")

	wpf(w, upload.BackendUIRendered().String())

	htmlfrag.Wb(w, "fsi tools", "")
	htmlfrag.Wb(w, "remove subtr", webapi.UriDeleteSubtree, " ")
	htmlfrag.Wb(w, "memfs dump", "/dump-memfs", " ")
	htmlfrag.Wb(w, "memfs reset", "/reset-memfs", " ")

	wpf(w, coinbase.BackendUIRendered().String())
	wpf(w, tplx.BackendUIRendered().String())
	wpf(w, login.BackendUIRendered().String())
	wpf(w, gitkit1.BackendUIRendered().String())

}

func serveFromRoot(w http.ResponseWriter, r *http.Request, replacements map[string][]byte) {

	appID := appengine.AppID(appengine.NewContext(r))
	if appID == AllowedAppID {

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
		fileserver.FsiFileServer(w, r, fileserver.Options{FS: fs1, Prefix: tplx.TplPrefix, Replacements: replacements})
	} else {
		w.Write([]byte("wrong app id: " + appID + "- "))
	}

}
