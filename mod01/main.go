package main

import (
	"net/http"

	"appengine"

	 _ "net/http/pprof"	
	"os"

	"fmt"
	"bytes"
 	"io/ioutil"	
	"io"
	"github.com/pbberlin/tools/util"
	"github.com/pbberlin/tools/conv"
	"github.com/pbberlin/tools/adapter"
	"github.com/pbberlin/tools/charting"
	 
)

func init() {
	
	//go main_ftp()
	
	http.HandleFunc("/"	 , adapter.AdapterAddC(homedir))
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


	fmt.Fprintln(w,"written via Fprintln<br>")


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


	b3 := new(bytes.Buffer)
	defer func(){
		w.Header().Set("Content-Type", "text/plain")
		w.Write( b3.Bytes() )		
	}()

	
}


func homedir(w http.ResponseWriter, r *http.Request, dir , base string, c appengine.Context) {
	
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)



	b1 := new(bytes.Buffer)


	wb(b1, "Guest Book", "" )
	wb(b1, "Eintrag hinzufügen", "/guest-entry" )
	wb(b1, "Einträge auflisten", "/guest-view" )
	wb(b1, "Einträge auflisten - paged - serialized cursor", "guest-view-cursor" )
	

	wb(b1, "Big Query", "" )
	wb(b1, "Get real data", "/big-query/query-into-datastore" )
	wb(b1, "Get mocked data", "/big-query/mock-data-into-datastore" )
	wb(b1, " ", "" )
	wb(b1, "Process Step 1 (optionally ?mock=1)",  "/big-query/regroup-data-01" )
	wb(b1, "Process Step 2",  "/big-query/regroup-data-02?f=table" )


	wb(b1, "Charts", "" )
	wb(b1, "Drawing a chart", "/image/draw-lines-example" )
	wb(b1, " ", "" )
	wb(b1, "Get image from Datastore", "/image/img-from-datastore?p=chart1" )
	wb(b1, "Get base64 from Datastore", "/image/base64-from-datastore?p=chart1" )
	wb(b1, "Get base64 from Variable", "/image/base64-from-var?p=1" )
	wb(b1, "Get base64 from File", "/image/base64-from-file?p=static/pberg1.png" )







	wb(b1, "Diverse", "" )
	wb(b1, "Letzte Email", "/email-view" )
	wb(b1, "Blob List", "/blob/list" )
	wb(b1, "Template Demo", "/tpl/demo" )
	wb(b1, "Http fetch", "/fetch-url" )



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
	s := fmt.Sprintf( "IntegerSequenes a, b: %v %v %v<br>\n",adapter.MyIntSeq01(), adapter.MyIntSeq01(), adapter.MyIntSeq02())
	b1.WriteString( s)

	b1.WriteString( "<br>\n")
	b1.WriteString( fmt.Sprintf( "Temp dir is %s<br>\n",os.TempDir() ))


	b1.WriteString( "<br>\n")
	b2 := new(bytes.Buffer)
	b2.WriteString("data:image/png;base64,...")
	b1.WriteString( fmt.Sprintf( "Mime from %q is %q<br>\n",b2.String(),conv.MimeFromBase64(b2) ))


	b1.WriteString( "<br>\n")
	b1.WriteString( fmt.Sprintf( "Last Month %q - 24 Months ago is %q<br>\n",util.MonthsBack(0),
	 util.MonthsBack(24) ))


	b1.WriteString( "<br>\n")
	sEnc := "Theo - wir fahrn nach Łódź. c a ff ee - trink nicht so viel Kaffee. "	
	buf1 ,msg1 := charting.StringToVByte(sEnc)
	//b1.WriteString( fmt.Sprint("string to byte in chunks says:",buf1.String(),"<br>" ) )
	b1.WriteString( fmt.Sprint(msg1) )

	b1.WriteString( fmt.Sprint(  "restore 1 s:= string([]bytes): ",  string(buf1.Bytes()),"<br>" ) )

	var bEnc []byte = []byte(sEnc)
	b1.WriteString( fmt.Sprint(  "restore 2 - from []byte(sEnc): ",  string(bEnc),"<br>" ) )
	


	


/*
*/

	
	w.Header().Set("Content-Type", "text/html")	
	w.Write( b1.Bytes() ) 
	
}


func wb( b1 *bytes.Buffer, linktext ,url string){

	if url == "" {
		b1.WriteString( "<br>\n")	
	}
	
	b1.WriteString( "<span style='display:inline-block; min-width:200px; margin: 6px 0px; margin-right:10px;'>\n")	
	if url == "" {
		b1.WriteString( "\t"+ linktext+"\n")	
	} else {
		b1.WriteString( "\t<a target='_app' href='"+url+"' >"+linktext+"</a>\n")	
	}
	b1.WriteString( "</span>\n")	
}
