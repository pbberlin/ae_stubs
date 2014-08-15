package main


import (
	"net/http"

	"appengine"
	 
	"net/url"
	"path"
	"regexp"
	"errors"
)


/* this is a closure, 
	  an anonymous func 
	  with surrounding variables 
	  
	  myIntSeq01(), myIntSeq02()
		 yield independent values for i
*/
func intSeq() func() int {
	 i := 0
	 return func() int {
		  i += 1
		  return i
	 }
}
var myIntSeq01 func()(int) = intSeq()
var myIntSeq02 func()(int) = intSeq()


/*
	http://golang.org/doc/articles/wiki/

		1.)  requi(a1) 
		2.)  given(a1,a2)
	=> 3.)  requi(a1) = adapter( given(a1,a2) )
	
	func adapter(	 given func( t1, t2)	){
		return func( a1 ) {						  // signature of requi
			a2 := something							// set second argument 
			given( a1, a2)
		}
	} 
	
	No chance for closure context variables.
	They can not flow into given(), 
	   because given() is not anonymous.
	
	
*/
func	adapterAddC(  given func(http.ResponseWriter, *http.Request, string, string, appengine.Context)	) http.HandlerFunc {
	
	
	return func(w http.ResponseWriter, r *http.Request) {


		c := appengine.NewContext(r)
		conditionTotal := regexp.MustCompile("^/([/a-zA-Z0-9\\.-]*)$") 
		m := conditionTotal.FindStringSubmatch( r.URL.Path )
		if m == nil {
			err := errors.New("illegal chars in path: " + r.URL.Path )
			c.Errorf("%v",err)
			return
		}

		s,_  := url.Parse( r.URL.String() )
		dir  := path.Dir( s.Path)
		base := path.Base(s.Path)
		given(w, r, dir , base,c)
	}
}

