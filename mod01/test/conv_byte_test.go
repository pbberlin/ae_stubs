package main

import (

	"testing"
	"appengine/aetest"
	
	"log"
	"image"
	"bytes"

	"github.com/pbberlin/tools/conv"
	"github.com/pbberlin/tools"

)

var str_b64 string = testString()
var debug bool = false
var use_test_string bool = false


func testString() string{
	b1 := new(bytes.Buffer)
	b1.WriteString("START_")
	for i:=1;i<12;i++{
		b1.WriteString("Line"+util.Itos(i)+"__a-b-c-d-e-f-g-h-i-j-k-l-m-n-o-p-q-r-s-t-u-v-w-x-y-z")	
	}
	b1.WriteString("_END")


	
	if use_test_string {
		return b1.String()			
	} else {
		return conv.Img_rgba_base64
	}
	
}


func Test_put_get(t *testing.T){

	c, err := aetest.NewContext(nil)	
	if err != nil {
		t.Errorf("could not get a context")		
	}


	vvbyte,_     := conv.String_to_vvbyte(str_b64)
	key_combi, _ := conv.Buf_put(c , conv.SWrapp{"test",vvbyte} , "test" )

	sw, _    := conv.Buf_get(c , key_combi)
	
	if debug {
		for i,v := range 	sw.Vvbyte {
			c.Errorf("%v  %s\n",i,v)	
		}
	}
	
	buff1,_  := conv.Vvbyte_to_string(sw.Vvbyte)
	if buff1.String() != str_b64{
		t.Errorf("put - get yields = %s", buff1.String() )
	}

	
}

func Test_string_to_img_and_back(t *testing.T){

	img,whichformat := conv.Base64_str_to_img( str_b64 )	
	log.Printf( "Retrieved img from base64 string: format %v - type %T\n" , whichformat, img )

	imgRGBA,ok := img.(*image.RGBA)
	if !ok {
		t.Errorf("could not cast loaded image to RGBA - is it png ?" )
		return
	}

	inverse := conv.Rgba_img_to_base64_str( imgRGBA )
	if debug {
		log.Printf("First 33 chars: %v", inverse[:33])
	}

	if inverse != str_b64 {
		t.Errorf("base64 - encode - decode yields = %s", inverse[:40] )
	}

	
}

// long string to slice of slice of byte (vector of vector of byte)
//   two inverse functions
func Test_string_to_vvbyte_and_back(t *testing.T) {


	vvbyte,_ := conv.String_to_vvbyte(str_b64)
	for i,v := range vvbyte {
		if debug { log.Printf("%v -  %s \n",i,v) }
	}	
	

	buff1,_ := conv.Vvbyte_to_string(vvbyte)
	if buff1.String() != str_b64{
		t.Errorf("encode - decode yields = %s", buff1.String() )
	}

}













// UNUSED
// long string to map of sbyte
//   two inverse functions
func Test_string_to_mapvbyte_and_back(t *testing.T) {

	// converting to map is pointless
	return
	


	mapr,_ := conv.String_to_mapvbyte(str_b64)
	for i,v := range mapr {
		log.Printf("%v -  %s \n",i,v)
	}	
	

	buff1,_ := conv.Mapvbyte_to_string(mapr)
	if buff1.String() != str_b64{
		t.Errorf("encode - decode yields = %s", buff1.String() )
	}

	vbyte,_ := conv.Mapvbyte_to_string_app(mapr)
	inverse := string(vbyte)
	if inverse != str_b64{
		t.Errorf("encode - decode yields = %s", inverse )
	}


}





