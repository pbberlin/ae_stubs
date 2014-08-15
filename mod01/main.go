package main

import (
	"net/http"

	"appengine"

	 _ "net/http/pprof"	

	"fmt"
	"bytes"
 	"io/ioutil"	
	"io"
	"github.com/pbberlin/tools"
	 
)

func init() {
	
	//go main_ftp()
	
	http.HandleFunc("/"	 , adapterAddC(homedir))
	http.HandleFunc("/login", login)
	http.HandleFunc("/guest-entry" , guestEntry)
	http.HandleFunc("/guest-save"  , guestSave)
	http.HandleFunc("/guest-view"  , guestView)
	http.HandleFunc("/guest-view-cursor" , guestViewCursor)
	
	http.Handle	("/json" , Servable_As_HTTP_JSON{Body:"myTitle",Title:"myBody"} )
	
	http.HandleFunc("/save-url" , saveURL_NoAnc)
	http.HandleFunc("/save-url-anc" , saveURL_WithAncestor)
	http.HandleFunc("/view-url" , viewURLAll)
	http.HandleFunc("/view-url-anc" , viewURLwithAncestors)
	
	
	  http.HandleFunc("/_ah/mail/"  , emailReceive1)
	//http.HandleFunc("/_ah/mail/"  , emailReceive2)
	http.HandleFunc("/email-view" , emailView)

	http.HandleFunc("/write-methods" , writeMethods)
	

}


func writeMethods(w http.ResponseWriter, r *http.Request) {
	
	
	w.Header().Set("Content-Type", "text/html")	

	defer r.Body.Close()


	// operations with a byte slice	
	var b2 []byte
	b2  = make( []byte , 100)
	b2[0] = 112
	b2[1] = 111
	b2[3] = 112
	b2[4] = 101
	b2[5] = 108
	b2[6] = 32

	w.Write(  b2          )
	w.Write(  []byte("<br>\n")    )
	
	bytesRead,_ := r.Body.Read( b2 )  // this reads 100 bytes from r.Body into b2, but r.Body is empty
	fmt.Fprint(w, "content of b2: " , string(b2) , " <br>bytes read " , bytesRead, "<br>")
	


	
	// operations with a bytes buffer
	var b1 *bytes.Buffer
	b1  = new(bytes.Buffer) // not optional on buffer pointer
	b1.ReadFrom( r.Body )
	b1.WriteString( "<br>\nstr_end_of_body<br>\n")
	fmt.Fprint(b1, "fmt.Fprinted into b1: number ", 222 , " EOL<br>\n")
	w.Write(  b1.Bytes()  )
	fmt.Fprint(w, b1.String() )
	// and copy the bytes.Buffer into w
	io.Copy(w, b1)


	// now using ioutil
	var content []byte
	content, _ = ioutil.ReadAll( r.Body )
	scont := string(content)
	cutoff := util.Min(20,len(scont))
	
	s := fmt.Sprintf("%v bytes => %v ...", len(scont), scont[:cutoff] )

	w.Write(  []byte("<br>\n==========<br>\n")  )
	w.Write(  []byte(s)  )

	fmt.Fprint(w, "<br>\n fmt.Print_string <br>\n")
	
	
	// finally: using io
	sio := "<br>\n finally using io.WriteString into w<br>\n"	
	io.WriteString(w, sio )


	
}


func homedir(w http.ResponseWriter, r *http.Request, dir , base string, c appengine.Context) {
	
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	b1 := new(bytes.Buffer)

	
	b1.WriteString( "<a target='_app' href='/guest-entry'>Eintrag ins Gästebuch</a><br>\n")
	b1.WriteString( "<a target='_app' href='/guest-view' >Alle Einträge</a><br>\n")
	b1.WriteString( "<a target='_app' href='/guest-view-cursor' >Alle Einträge, paged</a><br>\n")
	
	
	b1.WriteString( "<br>\n")
	b1.WriteString( "<a target='_app' href='/email-view' >Letzte Email</a><br>\n")
	b1.WriteString( "<a target='_app' href='/image-serve?mode=modified' >Chart</a><br>\n")



	b1.WriteString( "<br>\n")
	b1.WriteString( "<hr>\n")
	b1.WriteString( "<a target='_gae' href='https://console.developers.google.com/project/347979071940' ><b>global</b> developer console</a><br>\n")	
	b1.WriteString( " &nbsp; &nbsp; <a target='_gae' href='http://localhost:8000/mail' >app console local</a><br>\n")	
	b1.WriteString( " &nbsp; &nbsp; <a target='_gae' href='https://appengine.google.com/settings?&app_id=s~libertarian-islands' >app console online</a><br>\n")

	b1.WriteString( "<br>\n")
	b1.WriteString( "<a target='_gae'   href='http://localhost:8085/' >app local</a><br>\n")
	b1.WriteString( "<a target='_gae_r' href='http://libertarian-islands.appspot.com/' >app online</a><br>\n")



	b1.WriteString( "<br>\n")
	b1.WriteString( "Dir: "+dir+" -  Base: "+base+" <br>\n")
	
	b1.WriteString( "<br>\n")
	s := fmt.Sprintf( "IntegerSequenes a, b: %v %v %v<br>\n",myIntSeq01(), myIntSeq01(), myIntSeq02())
	b1.WriteString( s)
	
	w.Header().Set("Content-Type", "text/html")	
	w.Write( b1.Bytes() ) 
	
}


