package big_query



import (

	"fmt"
	"bytes"
	"appengine"
	"github.com/pbberlin/tools/util"
	"github.com/pbberlin/tools/util_err"

	"net/http"
	
	"encoding/gob"

	
)



// chart data
type CData struct {
	M map[string]map[string]float64 	
	VPeriods []string
	VLangs   []string
	F_max float64
	unexported string
}




// http://stackoverflow.com/questions/12854125/go-how-do-i-dump-the-struct-into-the-byte-array-without-reflection



func SaveChartDataToDatastore(w http.ResponseWriter, r *http.Request, cd  CData, key string ) string{

	internalType := fmt.Sprintf("%T",cd)
	//buffBytes, _	 := StringToVByte(s)  // instead of []byte(s)

	// CData to []byte
	serializedStruct := new(bytes.Buffer)
	enc := gob.NewEncoder(serializedStruct)
	err := enc.Encode(cd)
	util_err.Err_http(w,r,err,false)

	c := appengine.NewContext(r)	
	key_combi,err  := util.Buf_put(c , util.WrapBlob{Name:key, VByte:serializedStruct.Bytes(), S:internalType } , key )
	util_err.Err_http(w,r,err,false)

	return key_combi

	
}


func GetChartDataFromDatastore(w http.ResponseWriter, r *http.Request,  key string )(  *CData){

	c := appengine.NewContext(r)
	
	dsObj, err  := util.Buf_get(c , "util.WrapBlob__" + key)
	util_err.Err_http(w,r,err,false)


	newCData := new(CData)
	serializedStruct := bytes.NewBuffer( dsObj.VByte )
	dec := gob.NewDecoder(serializedStruct)
	err = dec.Decode(newCData)
	util_err.Err_http(w,r,err,false)


	return	newCData
}





//  if we want to gob.Encode/Decode unexported fields,
//  like CData.unexported, then we have to implement 
//  every field ourselves 
//  => uncomment following ...
/*
func (d *CData)GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(d.unexported)
	if err!=nil { return nil, err }
	return w.Bytes(), nil
}

func (d *CData)GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	return  decoder.Decode(&d.unexported)
}
*/

func TestDecodeEncode(w http.ResponseWriter, r *http.Request ) {

	// without custom implementation 
	// everything is encoded/decoded except field unexported
	d := CData{
		M: make( map[string]map[string]float64 )  ,
		VPeriods:  []string{"2011-11","2014-11",} ,
		VLangs:    []string{"C","++",},
		F_max:      44.2,
		unexported: "val of unexported",
	}



	// writing to []byte
	serializedStruct := new(bytes.Buffer)
	enc := gob.NewEncoder(serializedStruct)
	err := enc.Encode(d)
	util_err.Err_http(w,r,err,false)
	fmt.Fprintf(w,"encoded CData is: %#v \n<br>", serializedStruct.String()[0:20] )



	// reading
	e := new(CData)
	secondBuf := bytes.NewBuffer( serializedStruct.Bytes() )
	dec := gob.NewDecoder(secondBuf)
	err = dec.Decode(e)
	util_err.Err_http(w,r,err,false)


	fmt.Fprintf(w,"decoded CData is %#v \n<br>", e)


}