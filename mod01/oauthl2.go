package main

import (
	oauth2_common "github.com/golang/oauth2"
	oauth2_google "github.com/golang/oauth2/google"
	"net/http"
	"appengine"
)





// google
func oauthHandler2(w http.ResponseWriter, r *http.Request) {

	// get credentials from  https://console.developers.oauth2_google.com
	config, err := oauth2_google.NewServiceAccountConfig(&oauth2_common.JWTOptions{
		Email: "347979071940-livp1bl405mqq71797t3m99g3ghumu6p@developer.gserviceaccount.com",
		// The path to the pem file. If you have a p12 file instead, you
		// can use `openssl` to export the private key into a pem file.
		// $ openssl pkcs12 -in key.p12 -out key.pem -nodes
		
		/*
openssl pkcs12                               -in key.p12 -out key.pem -nodes  
openssl pkcs12 -password pass:               -in key.p12 -out key.pem -nodes  
openssl pkcs12 -password pass:pb165205       -in key.p12 -out key.pem -nodes  
openssl pkcs12 -password pass:347979071940   -in key.p12 -out key.pem -nodes  
		
		
		*/
		
		PemFilename: "/path/to/pem/file.pem",
		Scopes: []string{
			"https://www.googleapis.com/auth/bigquery",
		},
	})
	if err != nil {
		panic(err)
	}
	
	// Initiate an http.Client, the following GET request will be
	// authorized and authenticated on the behalf of
	// xxx@developer.gserviceaccount.com.
	client := http.Client{Transport: config.NewTransport()}
	client.Get("...")
	
	// If you would like to impersonate a user, you can
	// create a transport with a subject. The following GET
	// request will be made on the behalf of user@example.com.
	// client = http.Client{Transport: config.NewTransportWithUser("user@example.com")}
	// client.Get("...")
}




func oauthHandler3(w http.ResponseWriter, r *http.Request) {

	context := appengine.NewContext(r)
	config  := oauth2_google.NewAppEngineConfig(context, []string{
		"https://www.googleapis.com/auth/bigquery",
	})
	// The following client will be authorized by the App Engine
	// app's service account for the provided scopes.
	client := http.Client{Transport: config.NewTransport()}
	_ = client
	client.Get("...")	
}

func init() {
	http.HandleFunc("/oauth-handler2", oauthHandler2)	
	http.HandleFunc("/oauth-handler3", oauthHandler3)	
}
