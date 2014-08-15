package main

import (
	"time"
	 "net/http"
	 "appengine"
	ds "appengine/datastore"	
	"fmt"
	"bytes"
	"github.com/pbberlin/tools/util_err"
	
)






type LastURL struct {
	Value string
}


func ancKey( c appengine.Context ) *ds.Key {
	return ds.NewKey(c, "kindLastURLParent", "LastURLParent1", 0, nil)	
}


func saveURL_NoAnc(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)


	k := ds.NewKey(c, "childLastURL", "strKeyChildLastURL", 0, nil )
	

	e := new(LastURL)
   
   
   err := ds.Get(c, k, e)
	if err == ds.ErrNoSuchEntity {
		util_err.Err_log(err)
	} else if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	old := e.Value
	e.Value = r.URL.Path +"--"+ r.URL.RawQuery

	if _, err := ds.Put(c, k, e); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write( []byte("old="+old    +"\n" ) )
	w.Write( []byte("new="+e.Value+"\n" ) )


	 
}



func saveURL_WithAncestor(w http.ResponseWriter, r *http.Request) {

   c := appengine.NewContext(r)

	//k := ds.NewKey(c, "childLastURL", "strKeyChildLastURL", 0, ancKey(c) )
	k := ds.NewKey(c, "childLastURL", "", 0, ancKey(c) )

	lastURL_fictitious_1 := LastURL{"url_with_anc_1 " + timeMarker()}
	_, err := ds.Put(c, k, &lastURL_fictitious_1)
	check(w,err)

	lastURL_fictitious_2 := LastURL{"url_with_anc_2 " + timeMarker()}
	_, err = ds.Put(c, k, &lastURL_fictitious_2)
	check(w,err)
}

	

// get all URLs
func viewURLAll(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	
	b1 := new(bytes.Buffer)
	q := ds.NewQuery("childLastURL").
		Filter("Value >", "/save").
		Order("-Value")
	
	cntr := 0	
	for t := q.Run(c); ; {
		cntr++
		fmt.Fprint(b1, "q loop ", cntr, "\n")
		var lu LastURL
		key, err := t.Next(&lu)
		if err == ds.Done {
			b1.WriteString("\tq ds.Done\n")
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(b1, "\tKey=%v  \n\t%v\n", key, lu.Value)
	}

	 //fmt.Fprint(w, b1) 	
	 //io.Copy(w, b1)   	 

	 w.Write(  b1.Bytes()  )	 

	
}


func viewURLwithAncestors(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	
	q := ds.NewQuery("childLastURL").Ancestor(ancKey(c))
	var vURLs []LastURL
	keys, err := q.GetAll(c, &vURLs)
	check(w,err)
	
	for i,v := range vURLs{
		s := fmt.Sprint( i, keys[i], v,"\n")
		w.Write( []byte(s) )		
	}
	
	
}


func timeMarker() string{
	f2 := "2006-01-02 15:04:05"
	tn := time.Now()
	tn  = tn.Add( - time.Hour * 85 *24 )
	s2 := tn.Format( f2 ) 
	return s2
	
}