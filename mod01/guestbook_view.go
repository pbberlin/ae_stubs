package main

import (
    "net/http"

    "appengine"
    

	 sc "github.com/pbberlin/tools/sharded_counter"
	"github.com/pbberlin/tools/util_err"
	
	 
)



func guestEntry(w http.ResponseWriter, r *http.Request) {


	tplAdder,tplExec := funcTplBuilder(w,r)
	tplAdder("static_title","New guest book entry",nil)
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
		"static_title"    :   "List of guest book entries",
		"n_cont_0"       :   c_view_gbe,
		"n_cont_1":   "<pre>" + report + "</pre>",
		"n_cont_2":   s_cntr ,
	}
	myTplExecute(w,r,mc,gbEntries)
*/


	tplAdder,tplExec := funcTplBuilder(w,r)
	tplAdder("static_title","List of guest book entries",nil)
	tplAdder("n_cont_0",c_view_gbe,gbEntries)
	tplAdder("n_cont_1","<pre>{{.}}</pre>", report)
	tplAdder("n_cont_2","Visitors: {{.}}",cntr)
	tplExec(w,r)



}