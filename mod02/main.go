package main

import (
	"net/http"

	"appengine"

	"fmt"
	
	"log"

	// not used - but init() functions wanted for 
	// httpHandler registrations
	_ "github.com/pbberlin/tools/instance_mgt"

)


func init() {
	// DISTINCT from other modules
	
	http.HandleFunc("/mod02" , mainMod02)

	// _ah/start and _ah/stop  seem to be 
	// the instance start and stop requests

	log.Println("mod02 init complete")

}


func mainMod02(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)	
	module := appengine.ModuleName(c)

	w.Header().Set("Content-Type", "text/html")	

	instanceId := appengine.InstanceID()	
	
	w.Write( []byte(  fmt.Sprintf("Module -%v- <br>\n",module)   ) ) 
	w.Write( []byte(  fmt.Sprintf("Instance -%v- <br>\n",instanceId)   ) ) 


}

