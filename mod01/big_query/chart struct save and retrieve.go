package big_query

/*
	Normally, datastore types are restricted.
	For intstance a
	   map[string]map[string]float64 	
	can not be a datastore field.

	Therefore, this package takes a complex struct
	and *globs* it into a byte array,
	quasi normalizing it.
	
	It then saves the byte array within a util.WrapBlob 
	
	This way, any struct can be saved into into datastore
	using util.WrapBlob.


	Another aspect is the memoryInstanceStore.
	Data flows now as follows:
	
	instance[1]Memory < 
	instance[2]Memory < memCache < dataStore < bigQueryDB
	instance[3]Memory < 
	
	The invalidation of the leftmost tier is still an issue.
	Upon each 
	  memoryInstanceStore[key_combi] = newCData
	there should be a message sent to all other instances.
	
	Furthermore we would need a datastore versioning.
	
	
	
	
*/


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


var memoryInstanceStore map[string]*CData = make( map[string]*CData )


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

	memoryInstanceStore[key_combi] = &cd

	return key_combi
}


func GetChartDataFromDatastore(w http.ResponseWriter, r *http.Request,  key string )(  *CData){

	key_combi := "util.WrapBlob__" + key
	newCData := new(CData)

	newCData, ok := 	memoryInstanceStore[key_combi]
	util_err.Err_http(w,r,ok,true,"could not get it from memory")
	
	if !ok {
		c := appengine.NewContext(r)
		
		dsObj, err  := util.Buf_get(c , key_combi)
		util_err.Err_http(w,r,err,false)
	
		serializedStruct := bytes.NewBuffer( dsObj.VByte )
		dec := gob.NewDecoder(serializedStruct)
		err = dec.Decode(newCData)
		util_err.Err_http(w,r,err,false)

		memoryInstanceStore[key_combi] = newCData
		
	}


	return	newCData
}






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


//  if we wanted to gob.Encode/Decode unexported fields,
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
