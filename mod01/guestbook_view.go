package main

import (
    "fmt"
    "net/http"

    "appengine"
    

	 sc "github.com/pbberlin/tools/sharded_counter"
	"github.com/pbberlin/tools/util_err"
	
	 
)


func guestEntry(w http.ResponseWriter, r *http.Request) {
	mc  := map[string]string{
		"static_title"  :      "guestbook entry title",
		"static_content_1":     "",
		"static_content_2":     "",
	}	
	myTplExecute(w,mc,c_guestEntryHTML, nil)
}


func guestSave(w http.ResponseWriter, r *http.Request) {

	contnt := r.FormValue("content")
	entrySave(w,r,contnt)
	http.Redirect(w, r, "/guest-view", http.StatusFound)
    
}



func guestView(w http.ResponseWriter, r *http.Request) {



   c := appengine.NewContext(r)
	err := sc.Increment(c, "cGuestView" )
	util_err.Err_log(err)
	
	cntr, err := sc.Count(c, "cGuestView" ); util_err.Err_log(err)
	s_cntr := fmt.Sprint("<br>Counter Guest View is -",cntr,"-<br>\n")

	
	gbEntries, report  := guestbookEntries(w,r)

	mc  := map[string]string{
		"static_title"  :   "second title",
		"static_content_1":     "<pre>" + report + "</pre>",
		"static_content_2":     s_cntr ,
	}
	myTplExecute(w,mc,c_content_2, gbEntries)



}





