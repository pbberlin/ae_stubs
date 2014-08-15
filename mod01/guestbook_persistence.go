package main

import (
    "time"
     "net/http"
     "appengine"
    ds "appengine/datastore"    
    "appengine/user"

    "fmt"
    "strings"
    "bytes"
    "reflect"    
)


const (
	kind_guestbk     string = "Guestbook"         // "classname" of the parent
	kind_entry       string = "entityGreeting"    // "classname" of a guestbookk entry

	key_str_guestbk  string = "default_guestbook" // string Key

)



type Greeting1 struct {
    Author  string
    Content string      `datastore:"Content,noindex" json:"content"`
    Date    time.Time
    unsaved string      `datastore:"-"`
}

type Greeting2 struct {
    Author  string
    Content string
    Date    time.Time
    Field2  string
}


	// returns entity group - or parent - key 
	//   to store and retrieve all guestbook entries.
	//   the content of this parent is nil
	//   it only servers as umbrella for the entries
	func key_entity_group_key_parent(c appengine.Context)(r *ds.Key) {
	    // key_str_guestbk could be varied for multiple guestbooks.
	    // Either key_str_guestbk XOR key_int_guestbk must be zero
	    var key_int_guestbk int64 = 0  
	    var key_parent *ds.Key = nil
	    r = ds.NewKey(c, kind_guestbk, key_str_guestbk, key_int_guestbk, key_parent)
	    return
	}



func (g Greeting1) String() string {

	b1 := new(bytes.Buffer)
	
	b1.WriteString(g.Author + "<br>\n")
	b1.WriteString(g.Content+ "<br>\n")
	f2  := "2006-01-02 (Jan 02)"
	s2  := g.Date.Format( f2 ) 
	b1.WriteString(s2+ "<br>\n")

	return b1.String()	
}


func check(w http.ResponseWriter, e error){
	if e != nil {
		s := fmt.Sprint( "err: "+ e.Error() + "<br>\n" )
		w.Write( []byte(s) )
		//http.Error(w, err.Error(), http.StatusInternalServerError)
	}	
}




func entrySave(w http.ResponseWriter, r *http.Request, contnt string) {

    c := appengine.NewContext(r)
    
    g := Greeting1{
        Content: contnt,
        Date:    time.Now(),
    }
    if u := user.Current(c); u != nil {
        g.Author = u.String()
    }

    /* We set the same parent key on every Greeting entity 
       to ensure each Greeting is in the same entity group. 
       
       Queries across the single entity group will be consistent. 
       However, you should limit write rate to single entity group ~1/sec. 
    */
    
    // NewIncompleteKey(appengine.Context, kind string , parent *Key       ) 
    //  it has neither a string key, nor integer key 
    //  only a "kind" (classname) and a parent
    // Upon usage the datastore generates an integer key
    key := ds.NewIncompleteKey(         c, kind_entry, key_entity_group_key_parent(c) )
    discardedNewKey, err := ds.Put(c, key, &g)
    _ = discardedNewKey  // we query entries via key_entity_group_key_parent - via parent
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
}


func guestbookEntries(w http.ResponseWriter, r *http.Request)( greetings []Greeting2, report string ){

	c := appengine.NewContext(r)
	/* High Replication Datastore:
		Ancestor queries are strongly consistent. 
		Queries spanning MULTIPLE entity groups are EVENTUALLY consistent. 
		If .Ancestor was omitted from this query, there would be slight chance 
		that recent Greeting would not show up in a query.    */
	q := ds.NewQuery(kind_entry).Ancestor( key_entity_group_key_parent(c) ).Order("-Date").Limit(10)
	greetings  = make([]Greeting2, 0, 10)
	keys, err := q.GetAll(c, &greetings)
	if err != nil {	

		errtest,ok := err.(*ds.ErrFieldMismatch)
		_ = errtest
		if ok {
			// "types identical"
			err = nil
				
		}


		typeErr   := reflect.TypeOf(err)
		typeAsked := reflect.TypeOf(  new(ds.ErrFieldMismatch) )
		if typeErr == typeAsked {
			// "types identical"
			err = nil
		}
		
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


	
	var b1 bytes.Buffer
	var sw string
	var descrip []string = []string{"class","path","key_int_guestbk"}
	for i0,v0 := range keys{
		sKey := fmt.Sprintf("%v",v0)
		v1 := strings.Split( sKey,",")
	   sw = fmt.Sprintf("key %v",i0) ;b1.WriteString(sw)
		for i2, v2 := range v1{
			d := descrip[i2]
		   sw = fmt.Sprintf(" \t %v:  %q ",d,v2) ;b1.WriteString(sw)
		}
	   b1.WriteString("\n")
	}
	report = b1.String()

	return
}


