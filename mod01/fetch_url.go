package main

import (
	"fmt"
	"net/http"

	"appengine"
	"appengine/urlfetch"
	
	"io/ioutil"	
	"html"
	"github.com/pbberlin/tools/util"
	
)

func fetchURL(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	rURL := "http://www.google.com/"

	if r.PostFormValue("url") != "" {
		rURL = r.PostFormValue("url")		
	}


	

	tplAdder,tplExec := funcTplBuilder(w,r)
	tplAdder("static_title","Fetch some http data",nil)
	tplAdder("n_cont_0",c_formFetchURL,nil)
	tplExec(w,r)




	client := urlfetch.Client(c)
	resp, err := client.Get(rURL)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "HTTP GET returned status %v<br>\n\n", resp.Status)



	var s1, s2 string 
	if true && false {
		var b2 []byte
		b2 = make( []byte , 100)
		n, err := resp.Body.Read(b2 )	
		if err != nil {
			s1 =err.Error()
		} else {
			s1 =fmt.Sprintf("%v bytes read<br>", n )		
			s2 =string(b2)
		}
	} else {

		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			s1 = err.Error()
			c.Errorf("%s", err)
		} else {
			s1 =fmt.Sprintf("%v bytes read<br>", len(contents) )		
			s2 =string(contents)
		}
	}
	s2 = html.EscapeString(s2)
	
	fmt.Fprintf(w,s1)
	fmt.Fprintf(w,"\n\n")
	cutoff := util.Min(100,len(s2))
	fmt.Fprintf(w, "content is: <pre>" +  s2[:cutoff] + " ... " + s2[len(s2)-cutoff:] + "</pre>")	
	
}


func init() {
	http.HandleFunc("/fetch-url", fetchURL)
}