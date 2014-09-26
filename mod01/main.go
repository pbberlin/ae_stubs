package main

import (
	"net/http"

	
	 _ "net/http/pprof"	
	"os"

	"fmt"
	"log"
	"bytes"
 	"io/ioutil"	
	"io"
	"github.com/pbberlin/tools/util"
	"github.com/pbberlin/tools/conv"
	"github.com/pbberlin/tools/util_err"

	// not used - but init() functions wanted for 
	// httpHandler registrations
	_ "github.com/pbberlin/tools/big_query"
	_ "github.com/pbberlin/tools/blobstore_mgt"
	_ "github.com/pbberlin/tools/instance_mgt"
	_ "github.com/pbberlin/tools/guestbook"
	_ "github.com/pbberlin/tools/last_url"
)


func init() {

	
	//go main_ftp()
	
	http.HandleFunc("/"	 , util_err.Adapter(homedir))
	http.HandleFunc("/login", login)
	
	http.Handle	("/json" , Servable_As_HTTP_JSON{Body:"myTitle",Title:"myBody"} )
	
	
	
	  http.HandleFunc("/_ah/mail/"  , emailReceive1)
	//http.HandleFunc("/_ah/mail/"  , emailReceive2)
	http.HandleFunc("/email-view" , emailView)

	http.HandleFunc("/write-methods" , writeMethods)


	log.Println("mod01 (default) init complete")

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






func homedir(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)



	b1 := new(bytes.Buffer)



	wb(b1, "Diverse", "" )
	wb(b1, "Letzte Email", "/email-view" )
	wb(b1, "Blob List", "/blob/list" )
	wb(b1, "Template Demo 1", "/tpl/demo1" )
	wb(b1, "Template Demo 2", "/tpl/demo2" )
	wb(b1, "Http fetch", "/fetch-url" )
	wb(b1, "Instance Info", "/instance-info" )
	wb(b1, "Gob encode decode", "/big-query/test-gob-codec" )



	wb(b1, "Guest Book", "" )
	wb(b1, "Login", "/login" )
	wb(b1, "Eintrag hinzufügen", "/guest-entry" )
	wb(b1, "Einträge auflisten", "/guest-view" )
	wb(b1, "Einträge auflisten - paged - serialized cursor", "guest-view-cursor" )

	wb(b1, " ", "" )
	wb(b1, "Drawing a static chart", "/image/draw-lines-example" )
	

	wb(b1, "Big Query ...", "" )
	wb(b1, "Get real data", "/big-query/query-into-datastore" )
	wb(b1, "Get mocked data", "/big-query/mock-data-into-datastore" )
	wb(b1, "  &nbsp; &nbsp; &nbsp; ... with Chart", "" )
	wb(b1, "Process Data 1 (mock=1)",  "/big-query/regroup-data-01?mock=0" )
	wb(b1, "Process Data 2",  "/big-query/regroup-data-02?f=table" )
	wb(b1, "Show as Table",  "/big-query/show-table" )
	wb(b1, "Show as Chart",  "/big-query/show-chart" )
	wb(b1, "As HTML",  "/big-query/html" )



	wb(b1, "Request Images ", "" )
	wb(b1, "WrapBlob from Datastore", "/image/img-from-datastore?p=chart1" )
	wb(b1, "base64 from Datastore", "/image/base64-from-datastore?p=chart1" )
	wb(b1, "base64 from Variable", "/image/base64-from-var?p=1" )
	wb(b1, "base64 from File", "/image/base64-from-file?p=static/pberg1.png" )


	wb(b1, "Namespaces + Task Queues", "" )
	wb(b1, "Increment", "/namespaced-counters/increment" )
	wb(b1, "Read", "/namespaced-counters/read" )
	wb(b1, "Push to task-queue", "/namespaced-counters/queue-push" )

	wb(b1, "URLs with/without ancestors", "" )
	wb(b1, "Backend", "/save-url/backend" )




	b1.WriteString( "<br>\n")
	b1.WriteString( "<hr>\n")
	b1.WriteString( "<a target='_gae' href='https://console.developers.google.com/project/347979071940' ><b>global</b> developer console</a><br>\n")	
	b1.WriteString( " &nbsp; &nbsp; <a target='_gae' href='http://localhost:8000/mail' >app console local</a><br>\n")	
	b1.WriteString( " &nbsp; &nbsp; <a target='_gae' href='https://appengine.google.com/settings?&app_id=s~libertarian-islands' >app console online</a><br>\n")

	b1.WriteString( "<br>\n")
	b1.WriteString( "<a target='_gae'   href='http://localhost:8085/' >app local</a><br>\n")
	b1.WriteString( "<a target='_gae_r' href='http://libertarian-islands.appspot.com/' >app online</a><br>\n")



	dir  := m["dir"].(string)
	base := m["base"].(string)
	b1.WriteString( "<br>\n")
	b1.WriteString( "Dir: --"+dir+"-- &nbsp; &nbsp; &nbsp; &nbsp;   Base: --"+base+"-- <br>\n")
	
	b1.WriteString( "<br>\n")
	s := fmt.Sprintf( "IntegerSequenes a, b: %v %v %v<br>\n",util_err.MyIntSeq01(), util_err.MyIntSeq01(), util_err.MyIntSeq02())
	b1.WriteString( s)

	b1.WriteString( "<br>\n")
	b1.WriteString( fmt.Sprintf( "Temp dir is %s<br>\n",os.TempDir() ))


	b1.WriteString( "<br>\n")
	b2 := new(bytes.Buffer)
	b2.WriteString("data:image/png;base64,...")
	b1.WriteString( fmt.Sprintf( "Mime from %q is %q<br>\n",b2.String(),conv.MimeFromBase64(b2) ))


	b1.WriteString( "<br>\n")

	io.WriteString(b1, "Date: " + util.TimeMarker() + "  - " )
	b1.WriteString( fmt.Sprintf( "Last Month %q - 24 Months ago is %q<br>\n",util.MonthsBack(0),
	 util.MonthsBack(24) ))


	b1.WriteString( "<br>\n")
	sEnc := "Theo - wir fahrn nach Łódź."	
	b1.WriteString( fmt.Sprint(  "restore string string(  []byte(sEnc) ): ",  string(  []byte(sEnc) ),"<br>" ) )

	
	w.Header().Set("Content-Type", "text/html")	
	w.Write( b1.Bytes() ) 
	
}
