package main

import (
    "net/http"

    "appengine"
    

	 sc "github.com/pbberlin/tools/sharded_counter"
	"github.com/pbberlin/tools/util_err"
	
	"tpl_html"	 
)


const c_view_gbe = `
	{{range .}}
		{{$a := .Date}}
		{{$b := .Date  | df }}
		{{$c := df .Date}}
			<p>
		{{with .Author}}
			<b>{{.}}</b> wrote on {{$c}}<br>
		{{else}}
			An anonymous person wrote:   <br>
		{{end}}
			<span style='display:block; max-width:300px;font-size:12px;' >{{.Content}}</span>
		</p>
	{{end}}
`


const c_new_gbe = `
	<form action="/guest-save" method="post">
		<div><textarea name="content" rows="3" cols="60"></textarea></div>
		<div><input type="submit" value="Save Entry"></div>
	</form>
`





func guestEntry(w http.ResponseWriter, r *http.Request) {


	tplAdder,tplExec := tpl_html.FuncTplBuilder(w,r)
	tplAdder("n_html_title","New guest book entry",nil)
	tplAdder("n_cont_0",c_new_gbe,nil)

	tplExec(w,r)


}


func guestSave(w http.ResponseWriter, r *http.Request) {

	contnt := r.FormValue("content")
	entrySave(w,r,contnt)
	http.Redirect(w, r, "/guest-view", http.StatusFound)
    
}



func guestView(w http.ResponseWriter, r *http.Request) {



   c := appengine.NewContext(r)
	err := sc.Increment(c, "n_visitors_guestbook" )
	util_err.Err_http(w,r,err,false)	

	
	cntr, err := sc.Count(c, "n_visitors_guestbook" ); util_err.Err_http(w,r,err,false)

	
	gbEntries, report  := guestbookEntries(w,r)

/*
	mc  := map[string]string{
		"n_html_title"    :   "List of guest book entries",
		"n_cont_0"       :   c_view_gbe,
		"n_cont_1":   "<pre>" + report + "</pre>",
		"n_cont_2":   s_cntr ,
	}
	myTplExecute(w,r,mc,gbEntries)
*/


	tplAdder,tplExec := tpl_html.FuncTplBuilder(w,r)
	tplAdder("n_html_title","List of guest book entries",nil)
	tplAdder("n_cont_0",c_view_gbe,gbEntries)
	tplAdder("n_cont_1","<pre>{{.}}</pre>", report)
	tplAdder("n_cont_2","Visitors: {{.}}",cntr)
	tplExec(w,r)



}