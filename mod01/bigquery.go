package main

// https://godoc.org/code.google.com/p/google-api-go-client/bigquery/v2
// https://developers.google.com/bigquery/bigquery-api-quickstart
import (
	"net/http"
	bq			"code.google.com/p/google-api-go-client/bigquery/v2"
	oauth2_google "github.com/golang/oauth2/google"
	"appengine"
	"fmt"
	"github.com/pbberlin/tools/u_err"
	
)


func bqInit(w http.ResponseWriter, r *http.Request) {

	var q bq.QueryRequest = bq.QueryRequest{}
	q.Query = `
		SELECT 
		  repository_language   
		, LEFT(repository_pushed_at,7) monthx
		, CEIL( count(*)/1000) Tausend
		FROM githubarchive:github.timeline
		where 1=1
			AND  LEFT(repository_pushed_at,7) >= '2011-01'
			AND  repository_language in ('Go','go','Golang','golang','C','Java','PHP','JavaScript','C++','Python','Ruby')
			AND  type="PushEvent"
		group by monthx, repository_language
		order by repository_language   , monthx
		;	
	`


	context := appengine.NewContext(r)
	config  := oauth2_google.NewAppEngineConfig(context, []string{
		"https://www.googleapis.com/auth/bigquery",
	})
	// The following client will be authorized by the App Engine
	// app's service account for the provided scopes.
	client := http.Client{Transport: config.NewTransport()}
	//client.Get("...")	

	
	//oauthHttpClient := &http.Client{}
	bigqueryService, err := bq.New( &client )	
	util_err.Err_log(err)	
	
	fmt.Fprint(w,"s1<br>\n")
	
	// Create a query statement and query request object
	//  query_data = {'query':'SELECT TOP(title, 10) as title, COUNT(*) as revision_count FROM [publicdata:samples.wikipedia] WHERE wp_namespace = 0;'}
	//  query_request = bigquery_service.jobs()
	// Make a call to the BigQuery API
	//  query_response = query_request.query(projectId=PROJECT_NUMBER, body=query_data).execute()	



	js := bq.NewJobsService( bigqueryService )
	jqc := js.Query("347979071940", &q)

	fmt.Fprint(w,"s2 " + timeMarker()+" <br>\n")
	resp, err := jqc.Do()
	util_err.Err_log(err)
	
	
	
	rows := resp.Rows

	for i,v := range rows {

		s := fmt.Sprintf("Z%v \t\t",i)
		fmt.Fprint(w,s)
		
		cells := v.F
		for i1,v1 := range cells{
			val1 := v1.V
			s := fmt.Sprintf("c%v  %v  %T \t\t",i1,val1,val1)
			fmt.Fprint(w,s)
			//val2 := &v1.V
			//s  = fmt.Sprintf("\n\t\t%v  %T",val2,val2)
			//fmt.Fprint(w,s)
		}

		fmt.Fprint(w,"<br>\n")

	}
	
	fmt.Fprint(w,"s3 " + timeMarker()+" <br>\n")
	
	
}

func init() {
	http.HandleFunc("/bq-init", bqInit)	
}
