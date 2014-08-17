package main

import (
	"io"
	"net/http"
	//"net/url"
	
	"fmt"

	"appengine"
	"appengine/datastore"
	"appengine/taskqueue"
	"github.com/pbberlin/tools/u_err"
	
)

type Counter struct {
	Count int64
}

var sKey  string 	 = "SomeRequest"



func incrementCounterNamespaceAgnostic(c appengine.Context) error {
	key := datastore.NewKey(c, "NamespaceCounter", sKey, 0, nil)
	return datastore.RunInTransaction(c, func(c appengine.Context) error {
		var ctr Counter
		err := datastore.Get(c, key, &ctr)
		if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}
		ctr.Count++
		_, err = datastore.Put(c, key, &ctr)
		c.Infof("+1")
		return err
	}, nil)
}

func viewResetCounterNamespaceAgnostic(c appengine.Context, doReset bool)(int64,error ) {
	key := datastore.NewKey(c, "NamespaceCounter", sKey, 0, nil)

	var ctrRd Counter
	err := datastore.Get(c, key, &ctrRd)
	if err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}
	
	if doReset{
		var ctrSt Counter
		ctrSt.Count = -1
		_, err = datastore.Put(c, key, &ctrSt)		
	}
	
	return ctrRd.Count, err
}



func viewReset(w http.ResponseWriter, r *http.Request){
	
	var c1, c2 int64
	var s1, s2 string
	var err error
	var reset bool
	var pReset string
	
	pvReset  := r.URL.Query()["reset"]
	if len(pvReset) > 0 { 
			pReset = pvReset[0];
	}
	
	if pReset != "" { 
		reset = true 
	}
	
	c := appengine.NewContext(r)
	c1,err = viewResetCounterNamespaceAgnostic(c, reset )
	util_err.Err_log(err)


	{
		c , err = appengine.Namespace(c, "ns01")
		util_err.Err_log(err)
		c2,err = viewResetCounterNamespaceAgnostic(c, reset )
		util_err.Err_log(err)
	}

	s1 = fmt.Sprintf("%v",c1)
	s2 = fmt.Sprintf("%v",c2)


	io.WriteString(w, "|" + s1 + "|    |" + s2 + "|"   )
	if reset {  io.WriteString(w, "     and reset" )   }
	
}


func directIncBothNS(w http.ResponseWriter, r *http.Request){
	
	c := appengine.NewContext(r)
	err := incrementCounterNamespaceAgnostic(c)
	util_err.Err_log(err)


	{
		c , err := appengine.Namespace(c, "ns01")
		util_err.Err_log(err)
		err = incrementCounterNamespaceAgnostic(c)
		util_err.Err_log(err)
	}


	s:= `counters updates für ns=''  and ns='ns01'.` + "\n"
	io.WriteString(w,s )
	viewReset(w,r)
	
}




func enqueueIncrements(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	m     := map[string][]string{"counter_name": []string{"someRequest"},}
	t     := taskqueue.NewPOSTTask("/_ah/namespaced-counters/unqueue", m)


	taskqueue.Add(c, t, "")
	
	
	if true {
		c,err := appengine.Namespace(c, "ns01")
		util_err.Err_log(err)
		taskqueue.Add(c, t, "")
	}

	io.WriteString(w,"tasks enqueued\n" )
	viewReset(w,r)

}



func queuePop(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	err := incrementCounterNamespaceAgnostic(c)
	c.Infof("qp")
	util_err.Err_log(err)
}


func init() {
	 http.HandleFunc("/_ah/namespaced-counters/unqueue", queuePop)

	 http.HandleFunc("/namespaced-counters/direct-inc", directIncBothNS)

	 
	 http.HandleFunc("/namespaced-counters/queue-inc", enqueueIncrements)

	 http.HandleFunc("/namespaced-counters-view-reset", viewReset)
	
}

