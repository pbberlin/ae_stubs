package big_query

// https://godoc.org/code.google.com/p/google-api-go-client/bigquery/v2
// https://developers.google.com/bigquery/bigquery-api-quickstart
import (

	"bytes"
	"math/rand"
	"time"

	"net/http"
	bq			"code.google.com/p/google-api-go-client/bigquery/v2"
	oauth2_google "github.com/golang/oauth2/google"
	"appengine"
	"fmt"
	"github.com/pbberlin/tools/util_err"
	"github.com/pbberlin/tools/util"
	"github.com/pbberlin/tools/parsetools"
	"github.com/pbberlin/tools/adapter"


)

// chart data
type cdata struct {
	M map[string]map[string]float64 	
	VPeriods []string
	VLangs   []string
	F_max float64
}






// print it to http writer
func printPlaintextTable(w http.ResponseWriter, r *http.Request, vvDest [][]byte  ) {

	//c := appengine.NewContext(r)
	b1 := new(bytes.Buffer)
	defer func(){
		w.Header().Set("Content-Type", "text/plain")
		w.Write( b1.Bytes() )
	}()


	for i0 := 0 ; i0 < len(vvDest); i0++ {
		b1.WriteString("--" )
		b1.Write( vvDest[i0] )
		b1.WriteString("--" )
		b1.WriteString( "\n" )
	}

}


func queryIntoDatastore(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

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


	c := appengine.NewContext(r)
	config  := oauth2_google.NewAppEngineConfig(c, []string{
		"https://www.googleapis.com/auth/bigquery",
	})
	// The following client will be authorized by the App Engine
	// app's service account for the provided scopes.
	client := http.Client{Transport: config.NewTransport()}
	//client.Get("...")


	//oauthHttpClient := &http.Client{}
	bigqueryService, err := bq.New( &client )
	util_err.Err_http(w,r,err,false)


	fmt.Fprint(w,"s1<br>\n")

	// Create a query statement and query request object
	//  query_data = {'query':'SELECT TOP(title, 10) as title, COUNT(*) as revision_count FROM [publicdata:samples.wikipedia] WHERE wp_namespace = 0;'}
	//  query_request = bigquery_service.jobs()
	// Make a call to the BigQuery API
	//  query_response = query_request.query(projectId=PROJECT_NUMBER, body=query_data).execute()



	js := bq.NewJobsService( bigqueryService )
	jqc := js.Query("347979071940", &q)

	fmt.Fprint(w,"s2 " + util.TimeMarker()+" <br>\n")
	resp, err := jqc.Do()
	util_err.Err_http(w,r,err,false)



	rows := resp.Rows
	var vvDest [][]byte = make( [][]byte, len(rows) )

	c.Errorf("%#v",rows)

	for i0,v0 := range rows {

		cells := v0.F

		b_row := new(bytes.Buffer)
		b_row.WriteString( fmt.Sprintf("r%0.2d -- ",i0) )
		for i1,v1 := range cells{
			val1 := v1.V
			b_row.WriteString( fmt.Sprintf("c%0.2d: %v  ",i1,val1) )
		}
		vvDest[i0] = []byte( b_row.Bytes() )
	}

	key_combi,_  := util.Buf_put(c , util.WrapBlob{Name:"bq_res1",Vvbyte:vvDest} , "bq_res1" )
	dsObj,_  := util.Buf_get(c , key_combi)

	printPlaintextTable(w, r ,  dsObj.Vvbyte)


	fmt.Fprint(w,"s3 " + util.TimeMarker()+" <br>\n")


}

func mockDateIntoDatastore(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {


	c := appengine.NewContext(r)

	rand.Seed(time.Now().UnixNano())


	row_max := 100
	col_max := 3

	var languages[]string  = []string{"C","C++","Rambucto"}

	var vvDest [][]byte = make( [][]byte, row_max )
	for i0 := 0 ; i0 < row_max; i0++ {

		vvDest[i0] = make( []byte, col_max)

		b_row := new(bytes.Buffer)
		b_row.WriteString( fmt.Sprintf("r%0.2d -- ",i0) )

		for i1 := 0 ; i1 < col_max; i1++ {
			if i1 == 0 {
				val := languages[ i0/10 % 3 ]
				b_row.WriteString( fmt.Sprintf(" c%0.2d: %-10.8v  ",i1,val) )
			} else if i1 == 2 {
				val := rand.Intn(300)
				b_row.WriteString( fmt.Sprintf(" c%0.2d: %10v  ",i1,val) )
			} else {

				f2 := "2006-01-02 15:04:05"
				f2  = "2006-01"
				tn := time.Now()
				//tn  = tn.Add( - time.Hour * 85 *24 )
				tn  = tn.Add( - time.Hour * time.Duration(i0) *24 )
				val := tn.Format( f2 )
				b_row.WriteString( fmt.Sprintf(" c%0.2d: %v  ",i1,val) )
			}
		}
		vvDest[i0] = []byte( b_row.Bytes() )

	}

	key_combi,_  := util.Buf_put(c , util.WrapBlob{Name:"bq_res_test",Vvbyte:vvDest} , "bq_res_test" )
	dsObj,_  := util.Buf_get(c , key_combi)

	printPlaintextTable(w, r ,dsObj.Vvbyte)

}


func regroupFromDatastore01(w http.ResponseWriter, r *http.Request, m map[string]interface{} ) {

	c := appengine.NewContext(r)
	b1 := new(bytes.Buffer)
	defer func(){
		w.Header().Set("Content-Type", "text/html")
		w.Write( b1.Bytes() )
	}()


	var vvSrc [][]byte

	if util.IsLocalEnviron() {
		vvSrc = bq_statified_res1	
	} else {
		dsObj1,_ := util.Buf_get(c , "util.WrapBlob__bq_res1")
		vvSrc = dsObj1.Vvbyte
	}


	if r.FormValue("mock") != "" {
		dsObj1,_ := util.Buf_get(c , "util.WrapBlob__bq_res_test")
		vvSrc = dsObj1.Vvbyte
	}


	var vvDest [][]byte = make( [][]byte, len(vvSrc) )


	for i0 := 0 ; i0 < len(vvSrc); i0++ {


		s_row := string(vvSrc[i0])
		v_row := parsetools.SplitByWhitespace(s_row)
		b_row := new(bytes.Buffer)


		b_row.WriteString( fmt.Sprintf("%16.12s   ", v_row[3]) )
		b_row.WriteString( fmt.Sprintf("%16.12s   ", v_row[5]) )
		b_row.WriteString( fmt.Sprintf("%16.8s"    , v_row[7]) )

		vvDest[i0] = []byte( b_row.Bytes() )

	}


	key_combi,_  := util.Buf_put(c , util.WrapBlob{Name:"res_processed_01",Vvbyte:vvDest} , "res_processed_01" )
	dsObj2,_  := util.Buf_get(c , key_combi)

	printPlaintextTable(w, r ,dsObj2.Vvbyte)


}



func regroupFromDatastore02(w http.ResponseWriter, r *http.Request, m map[string]interface{} ) {

	c := appengine.NewContext(r)
	b1 := new(bytes.Buffer)
	defer func(){
		w.Header().Set("Content-Type", "text/html")
		w.Write( b1.Bytes() )
	}()


	var vvSrc [][]byte
	dsObj1,err := util.Buf_get(c , "util.WrapBlob__res_processed_01")
	util_err.Err_http(w,r,err,false)	
	vvSrc = dsObj1.Vvbyte



	d  := make(    map[string]map[string]float64 )
	
	
	distinctLangs   := make(map[string]interface{})
	distinctPeriods := make(map[string]interface{})
	f_max := 0.0
	for i0 := 0 ; i0 < len(vvSrc); i0++ {
		//vvDest[i0] = []byte( b_row.Bytes() )
		s_row := string(vvSrc[i0])		
		v_row := parsetools.SplitByWhitespace(s_row)

		lang   := v_row[0]
		period := v_row[1]
		count  := v_row[2]
		fCount := util.Stof(count)
		if fCount > f_max { f_max = fCount }
		

		distinctLangs[ lang ] = 1
		distinctPeriods[ period ] = 1

		if _,ok := d[period] ; ! ok {
			d[period] = map[string]float64{}
		}
		d[period][lang] = fCount
		
	}
	//fmt.Fprintf(w,"%#v\n",d2)	
	//fmt.Fprintf(w,"%#v\n",f_max)	

	sortedPeriods := util.KeysToSortedArray(distinctPeriods)
	sortedLangs := util.KeysToSortedArray(distinctLangs)

	cd := cdata{}
	_ = cd
	
	cd.M = d
	cd.VPeriods = sortedPeriods 
	cd.VLangs   = sortedLangs
	cd.F_max    = f_max
	


	if r.FormValue("f") == "table" {
		showAsTable(w,r,cd)
	} else {
		showAsChart(w,r,cd)		
	}




}

func showAsTable(w http.ResponseWriter, r *http.Request, cd cdata){

	span := util.GetSpanner()
	// Header row
	fmt.Fprintf(w,span(" ",164)	)		
	for _,lg := range cd.VLangs {
		fmt.Fprintf(w,span(lg,88)	)		
	}
	fmt.Fprintf(w,"<br>")		
	
	for _,period := range cd.VPeriods {
		fmt.Fprintf(w,span(period,164)	)		
		for _,lg := range cd.VLangs {
			fmt.Fprintf(w,span( cd.M[period][lg]  ,88)	)		
		}
		fmt.Fprintf(w,"<br>")		
	}

	
}

func init() {

	http.HandleFunc("/big-query/query-into-datastore", adapter.Adapter(queryIntoDatastore) )
	http.HandleFunc("/big-query/mock-data-into-datastore", adapter.Adapter(mockDateIntoDatastore) )
	http.HandleFunc("/big-query/regroup-data-01", adapter.Adapter(regroupFromDatastore01) )
	http.HandleFunc("/big-query/regroup-data-02", adapter.Adapter(regroupFromDatastore02) )
}


