package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	_ "net/http/pprof"

	"github.com/pbberlin/tools/util"
	"github.com/pbberlin/tools/util_appengine"
	"github.com/pbberlin/tools/util_err"

	"appengine"
	"appengine/urlfetch"
)

var serveraddress string = `header of ["SERVER_ADDR"]`
var pos_01 int = strings.Index(serveraddress, "192.")

var dns_router string = "see func init"
var dns_cam string = "see func init"

const (
	path_get_alarm = "/cgi-bin/CGIProxy.fcgi?cmd=getMotionDetectConfig&usr=admin&pwd=pb165205"
	path_set_alarm = "/cgi-bin/CGIProxy.fcgi?cmd=setMotionDetectConfig&usr=admin&pwd=pb165205"
	AREAS          = `&area0=255&area1=255&area2=255&area3=255&area4=255&area5=255&area6=255&area7=255&area7=255&area8=255&area9=255`
)

type CGI_Result struct {
	Result          string `xml:"result"`
	IsEnable        string `xml:"isEnable"`
	Linkage         string `xml:"linkage"`
	SnapInterval    string `xml:"snapInterval"`
	Sensitivity     string `xml:"sensitivity"`
	TriggerInterval string `xml:"triggerInterval"`
	Schedule0       string `xml:"schedule0"`
	Schedule1       string `xml:"schedule1"`
	Schedule2       string `xml:"schedule2"`
	Schedule3       string `xml:"schedule3"`
	Schedule4       string `xml:"schedule4"`
	Schedule5       string `xml:"schedule5"`
	Schedule6       string `xml:"schedule6"`
	Area0           string `xml:"area0"`
	Area1           string `xml:"area1"`
	Area2           string `xml:"area2"`
	Area3           string `xml:"area3"`
	Area4           string `xml:"area4"`
	Area5           string `xml:"area5"`
	Area6           string `xml:"area6"`
	Area7           string `xml:"area7"`
	Area8           string `xml:"area8"`
	Area9           string `xml:"area9"`
	// XMLName  xml.Name `xml:"account"`
}

func foscamToggle(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	util.SetNocacheHeaders(w, true)
}

func foscamStatus(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	util.SetNocacheHeaders(w, false)

	ssecs := r.FormValue("sleep")
	if ssecs != "" {
		secs := util.Stoi(ssecs)
		opf(w, "sleeping %v secs ... <br><br>\n", secs)
		time.Sleep(time.Duration(secs) * time.Second)
	}

	opf(w, "sta foscam status<br>\n")

	c := appengine.NewContext(r)
	client := urlfetch.Client(c)

	url := spf(`http://%s/%s`, dns_cam, path_get_alarm)
	opf(w, "requesting %v<br>\n", url)
	resp1, err := client.Get(url)
	util_err.Err_http(w, r, err, false)

	opf(w, "reading resp1<br>\n")
	bcont, err := ioutil.ReadAll(resp1.Body)
	defer resp1.Body.Close()
	util_err.Err_http(w, r, err, false)

	cgiRes := CGI_Result{}
	xmlerr := xml.Unmarshal(bcont, &cgiRes)
	util_err.Err_http(w, r, xmlerr, false)

	psXml := util.IndentedDump(cgiRes)
	opf(w, "<pre>%v</pre>", *psXml)

	// scont := string(bcont)
	// // scont = util.Ellipsoider(scont, 250)
	// opf(w, "<pre>%v</pre>", scont)

	opf(w, "end foscam status<br>\n")

}

func init() {

	if util_appengine.IsLocalEnviron() {
		dns_router = "192.168.1.1"
		dns_cam = "192.168.1.4:8081"
	} else {
		dns_router = "ds7934.myfoscam.org"
		dns_cam = "ds7934.myfoscam.org:8081"
	}

	http.HandleFunc("/foscam-status", util_err.Adapter(foscamStatus))
	http.HandleFunc("/foscam-toggle", util_err.Adapter(foscamToggle))

}

func toggle_alarm(doSwitch bool) {

	// 	$_params = str_replace("&result=0","", $_params );
	// 	$pos2 = strpos( $_params , "&isEnable=0"  );
	// 	$pos3 = strpos( $_params , "&isEnable=1"  );
	// 	if ($pos2 === false  AND $pos3 === false){
	// 		echo "no enabled param<br>";
	// 		exit();
	// 	} else {
	// 		if( $_switch ){
	// 			//$_params = preg_replace('/&area[0-9]=[0-9]{4,4}/i',"",$_params);
	// 			$_params = preg_replace('/\&area[0-9]=[0-9]+/i',"",$_params);
	// 			$_params .= $_AREAS;

	// 			if( $pos2 !== false ){
	// 				$msg1 = "now switched ON";
	// 				$_params = str_replace("&isEnable=0","&isEnable=1", $_params );
	// 				$_state=1;
	// 			} else {
	// 				$msg1 = "now switched OFF";
	// 				$_params = str_replace("&isEnable=1","&isEnable=0", $_params );
	// 				$_state=0;
	// 			}
	// 			$x = url_get_contents_x2($_dns_cam . $_path_set_alarm . $_params, false, false);

	// 		} else {
	// 			if( $pos2 !== false ){
	// 				$_state=0;
	// 			} else {
	// 				$_state=1;
	// 			}
	// 			$msg1 = "no switching";
	// 		}
	// 	}    // endif has return param "enabled"
	// }       // endif has return param "result"
	// return array( "msg" => $msg1 , "params" => $_params, "state" => $_state );
}
