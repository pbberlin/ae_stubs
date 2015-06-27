package main

import (
	"testing"

	"appengine/aetest"

	"bytes"
	"image"
	"log"

	"github.com/pbberlin/tools/conv"
	"github.com/pbberlin/tools/dsu"
	"github.com/pbberlin/tools/util"
)

var str_b64 string = testString()
var debug bool = false
var use_test_string bool = false

func testString() string {
	b1 := new(bytes.Buffer)
	b1.WriteString("START_")
	for i := 1; i < 12; i++ {
		b1.WriteString("Line" + util.Itos(i) + "__a-b-c-d-e-f-g-h-i-j-k-l-m-n-o-p-q-r-s-t-u-v-w-x-y-z")
	}
	b1.WriteString("_END")

	if use_test_string {
		return b1.String()
	} else {
		return conv.Img_rgba_base64
	}

}

func Test_put_get(t *testing.T) {

	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Errorf("could not get a context")
	}

	VVByte, _ := conv.String_to_VVByte(str_b64)
	key_combi, _ := dsu.BufPut(c, dsu.WrapBlob{"test", VVByte}, "test")

	sw, _ := dsu.BufGet(c, key_combi)

	if debug {
		for i, v := range sw.VVByte {
			c.Errorf("%v  %s\n", i, v)
		}
	}

	buff1, _ := conv.VVByte_to_string(sw.VVByte)
	if buff1.String() != str_b64 {
		c.Errorf("put - get yields = %s", buff1.String())
	}

}

func Test_string_to_img_and_back(t *testing.T) {

	img, whichFormat := conv.Base64_str_to_img(str_b64)
	log.Printf("Retrieved img from base64 string: format %v - type %T\n", whichFormat, img)

	imgRGBA, ok := img.(*image.RGBA)
	if !ok {
		t.Errorf("could not cast loaded image to RGBA - is it png ?")
		return
	}

	inverse := conv.Rgba_img_to_base64_str(imgRGBA)
	if debug {
		log.Printf("First 33 chars: %v", inverse[:33])
	}

	if inverse != str_b64 {
		t.Errorf("base64 - encode - decode yields = %s", inverse[:40])
	}

}

// long string to slice of slice of byte (vector of vector of byte)
//   two inverse functions
func Test_string_to_VVByte_and_back(t *testing.T) {

	VVByte, _ := conv.String_to_VVByte(str_b64)
	for i, v := range VVByte {
		if debug {
			log.Printf("%v -  %s \n", i, v)
		}
	}

	buff1, _ := conv.VVByte_to_string(VVByte)
	if buff1.String() != str_b64 {
		t.Errorf("encode - decode yields = %s", buff1.String())
	}

}

type SomeStruct struct {
	S1 string `json:"s1"`
	S2 string `json:"s2"`
}

func Test_meMcacheGet_set(t *testing.T) {

	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Errorf("could not get a context")
	}

	rvb1, ok := dsu.McacheGet(c, "key1", &SomeStruct{})
	if debug {
		log.Printf("  -------%#v--%T--%v--- \n\n", rvb1, rvb1, ok)
	}

	rvb2, ok := dsu.McacheGet(c, "key2", "string")
	if debug {
		log.Printf("  -------%#v--%T--%v--- \n\n", rvb2, rvb2, ok)
	}

	rvb3, ok := dsu.McacheGet(c, "key3", 22323)
	if debug {
		log.Printf("  -------%#v--%T--%v--- \n\n", rvb3, rvb3, ok)
	}

	dsu.McacheSet(c, "key1", "just a scalar stupid string")

	myStruct1 := SomeStruct{"this content", "is structured"}
	dsu.McacheSet(c, "key2", myStruct1)

	myStruct2 := SomeStruct{"wonderbar", "is not wonderbra"}
	dsu.McacheSet(c, "key3", myStruct2)

	rva1, ok := dsu.McacheGet(c, "key1", "string")
	if debug {
		log.Printf("  -------%#v--%T--%v--- \n\n", rva1, rva1, ok)
	}

	rva2, ok := dsu.McacheGet(c, "key1", 22323)
	if debug {
		log.Printf("  -------%#v--%T--%v--- \n\n", rva2, rva2, ok)
	}

	rva3, ok := dsu.McacheGet(c, "key2", "string")
	if debug {
		log.Printf("  -------%#v--%T--%v--- \n\n", rva3, rva3, ok)
	}

	rva4 := SomeStruct{}
	_, ok = dsu.McacheGet(c, "key3", &rva4)
	if debug {
		log.Printf("  -------%#v--%T--%v--- \n\n", rva4, rva4, ok)
	}

	rva5 := SomeStruct{}
	dsu.McacheGet(c, "key2", &rva5)
	if debug {
		log.Printf("  --%#v--- \n\n", rva5)
	}

	want1 := "just a scalar stupid string"
	if rva1 != want1 {
		t.Errorf("memache get - set - want %s, got ", want1, rva1)
	}

	want2 := 0
	if rva2 != want2 {
		t.Errorf("memache get - set - want %s, got ", want2, rva2)
	}

	want3 := "{\"s1\":\"this content\",\"s2\":\"is structured\"}"
	if rva3 != want3 {
		t.Errorf("memache get - set - want %s, got ", want3, rva3)
	}

	want4a := "wonderbar"
	want4b := "is not wonderbra"

	if rva4.S1 != want4a || rva4.S2 != want4b {
		t.Errorf("memache get - set - wanted %s %s, got %#v", want4a, want4b, rva4)
	}

	want5a := "this content"
	want5b := "is structured"

	if rva5.S1 != want5a || rva5.S2 != want5b {
		t.Errorf("memache get - set - wanted %s %s, got %#v", want5a, want5b, rva5)
	}

}
