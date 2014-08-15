
// This program makes a call to the specified API, authenticated with OAuth2.
// Examples: https://code.google.com/oauthplayground/
// 			 https://developers.google.com/oauthplayground/?code=4/If8JRiGJHmIj_SmiVKuWnzZhpLqh.4rKSn04vDfEeyjz_MlCJoi0i4oB-jwI
// https://code.google.com/p/goauth2/source/browse/oauth/oauth.go
// https://developers.google.com/appengine/docs/go/oauth/
// https://developers.google.com/console/help/new/#generatingoauth2
// https://developers.google.com/accounts/docs/OAuth2
// https://developers.google.com/accounts/docs/OAuth2?csw=1
// https://developers.google.com/accounts/docs/OAuth2ServiceAccount
// => https://github.com/golang/oauth2
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"net/http"
	"code.google.com/p/goauth2/oauth"
)

var (
	// service account client id - https://developers.google.com/accounts/docs/OAuth2?csw=1#serviceaccount
	clientId     = flag.String("id", "347979071940-u7mogetlcu3l47ocqk1mqm6b6tl8sbiu.apps.googleusercontent.com", "Client ID")
	clientSecret = flag.String("secret", "AwR-DDoRIK3iG9ai-4KP7rJm", "Client Secret")

	// for example "https://www.googleapis.com/auth/buzz",
	scope	= flag.String("scope", "https://www.googleapis.com/auth/userinfo.profile", "OAuth scope")
	
	authURL	= flag.String("auth_url" , "https://accounts.google.com/o/oauth2/auth", "Authentication URL")
	tokenURL     = flag.String("token_url", "https://accounts.google.com/o/oauth2/token", "Token URL")

	//RedirectURL  "http://you.example.org/handler",
	redirectURL  = flag.String("redirect_url", "oob", "Redirect URL")

	requestURL   = flag.String("request_url", "https://www.googleapis.com/oauth2/v1/userinfo", "API request")

	code	 = flag.String("code", "", "Authorization Code")
	cachefile    = flag.String("cache", "cache.json", "Token cache file")
)

const usageMsg = `
To obtain a request token you must specify both -id and -secret.

To obtain Client ID and Secret, see the "OAuth 2 Credentials" section under
the "API Access" tab on this page: https://code.google.com/apis/console/

Once you have completed the OAuth flow, the credentials should be stored inside
the file specified by -cache (cache.json) and you may run without the -id and -secret flags.
`

// Set up a configuration.
var config *oauth.Config  = &oauth.Config{
	ClientId:     *clientId,
	ClientSecret: *clientSecret,
	RedirectURL:  *redirectURL,
	Scope:	     *scope,
	AuthURL:	     *authURL,
	TokenURL:     *tokenURL,
	TokenCache:   oauth.CacheFile(*cachefile),
}


func oauthFlow() {
	flag.Parse()


	// Set up a Transport using the config.
	transport := &oauth.Transport{Config: config}

	// Try to pull the token from the cache; if this fails, we need to get one.
	token, err := config.TokenCache.Token()
	if err != nil {
		if *clientId == "" || *clientSecret == "" {
			flag.Usage()
			fmt.Fprint(os.Stderr, usageMsg)
			os.Exit(2)
		}
		if *code == "" {
			// Get an authorization code from the data provider.
			// ("Please ask the user if I can access this resource.")
			url := config.AuthCodeURL("")
			fmt.Print("Visit this URL to get a code, then run again with -code=YOUR_CODE\n\n")
			fmt.Println(url)
			return
		}
		// Exchange the authorization code for an access token.
		// ("Here's the code you gave the user, now give me a token!")
		token, err = transport.Exchange(*code)
		//           transport.Exchange(r.FormValue("code"))
		if err != nil {
			log.Fatal("Exchange:", err)
		}
		// The Transport now has a valid Token. Create an *http.Client
		// with which we can make authenticated API requests.
		// c := transport.Client()
		// c.Post(...)
		
		
		
		// (The Exchange method will automatically cache the token.)
		fmt.Printf("Token is cached in %v\n", config.TokenCache)
	}

	// Make the actual request using the cached token to authenticate.
	// ("Here's the token, let me in!")
	transport.Token = token

	// Make the request.
	r, err := transport.Client().Get(*requestURL)
	if err != nil {
		log.Fatal("Get:", err)
	}
	defer r.Body.Close()

	// Write the response to standard output.
	io.Copy(os.Stdout, r.Body)

	// Send final carriage return, just to be neat.
	fmt.Println()
}



// A landing page redirects to the OAuth provider to get the auth code.
func oauthLanding(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, config.AuthCodeURL("foo"), http.StatusFound)
}

// The user will be redirected back to this handler, that takes the
// "code" query parameter and Exchanges it for an access token.
func oauthHandler(w http.ResponseWriter, r *http.Request) {

	oauthFlow()

}


func init() {
	http.HandleFunc("/oauth-handler", oauthHandler)	
	http.HandleFunc("/oauth-landing", oauthLanding)	
}
