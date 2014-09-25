package main

import   "testing"
import   "appengine/aetest"
import	"github.com/pbberlin/tools/util"
import   "github.com/pbberlin/tools/util_err"



func Test_print_map1(t *testing.T){
	
	c, err := aetest.NewContext(nil)	
	if err != nil {
		t.Errorf("could not get a context")		
	}
	
	
	s := util.PrintMap(util.Map_example_right)
	c.Infof("externally testing print map recursive ...")		
	if  util.Test_want != s {
		c.Errorf("want \n%s \ngot \n%s", 
			util.Test_want[0:22]+" ... " + util.Test_want[len(util.Test_want)-22:], 
			s[0:22] 		        +" ... " + s[len(s)-22:] )
	}
}


