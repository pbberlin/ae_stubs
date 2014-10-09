package main

import (
	"net/http"

	_ "net/http/pprof"

	"appengine"
	"appengine/urlfetch"
	"bytes"
	"fmt"
	"github.com/pbberlin/tools/instance_mgt"
	"github.com/pbberlin/tools/util"
	"github.com/pbberlin/tools/util_err"
	"io"
	"io/ioutil"
)

func writeMethods(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	c := appengine.NewContext(r)

	client := urlfetch.Client(c)

	ii := instance_mgt.Get(w, r, m)
	resp2, err := client.Get(sp(`http://%s/write-methods-read`, ii.Hostname))
	util_err.Err_http(w, r, err, false)

	bufDemo := new(bytes.Buffer)
	bufDemo.WriteString("end of page")
	defer func() {
		//w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write(bufDemo.Bytes())

		resp2.Body.Close()
	}()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<pre>")

	//
	//
	fmt.Fprint(w, `Most basic:
		this is written via Fprintln into response writer`+"\n\n\n")

	//
	// byte slice into response writer
	const sz = 20
	var sB []byte
	sB = make([]byte, sz)
	sB[0] = 112
	sB[1] = 111
	sB[2] = '-'
	sB[3] = 112
	sB[4] = 101
	sB[5] = 108
	sB[6] = 32
	for i := 7; i < sz; i++ {
		sB[i] = ' '
	}
	sB[sz-1] = '!'

	w.Write([]byte("Byte slice into response writer: \n\t\t"))
	w.Write(sB)
	w.Write([]byte("\n\n\n"))

	//
	//
	// resp2.Body into byte slice,
	sB2 := make([]byte, sz)
	for i := 0; i < sz; i++ {
		sB2[i] = '-'
	}
	bytesRead, err := resp2.Body.Read(sB2)
	if err == nil {
		fmt.Fprintf(w, "Byte slice - reading %v bytes from response-body\n\t\t%q \n\n\n",
			bytesRead, string(sB2))
	} else {
		fmt.Fprintf(w, "err reading into byte slice  --%v-- \n\n\n", err)
	}

	//
	//
	//
	fp(w, "operations with a bytes buffer\n")
	var buf1 *bytes.Buffer
	buf1 = new(bytes.Buffer) // not optional on buffer pointer
	buf1.ReadFrom(resp2.Body)

	buf1 = new(bytes.Buffer)
	fp(buf1, "\t\tbuf1 content %v (filled via Fprintf)\n", 222)

	fp(w, "FOUR methods of dumping buf1 into resp.w:\n")
	fp(w, "\tw.Write\n")
	w.Write(buf1.Bytes())
	fp(w, "\tFprint\n")
	fp(w, buf1.String())
	fp(w, "\tio.WriteString\n")
	io.WriteString(w, buf1.String())
	fp(w, "\tio.Copy \n")
	io.Copy(w, buf1) // copy the bytes.Buffer into w
	fp(w, " \t\t\tio.copy exhausts buf1 - Fprinting again yields %q ", buf1.String())
	fp(w, buf1.String())
	fp(w, "\n\n\n")

	//
	//
	//
	fp(w, "ioutil.ReadAll\n")
	var content []byte
	resp3, err := client.Get(sp(`http://%s/write-methods-read`, ii.Hostname))
	util_err.Err_http(w, r, err, false)
	content, _ = ioutil.ReadAll(resp3.Body)
	scont := string(content)
	scont = util.Ellipsoider(scont, 20)
	w.Write([]byte(scont))

	fmt.Fprint(w, "</pre>")

}

// simple helper for reading http.response.Body
func writeMethodsResponder(w http.ResponseWriter, r *http.Request) {
	fp(w, "some http response body string")
}

func init() {
	http.HandleFunc("/write-methods", util_err.Adapter(writeMethods))
	http.HandleFunc("/write-methods-read", writeMethodsResponder)
}
