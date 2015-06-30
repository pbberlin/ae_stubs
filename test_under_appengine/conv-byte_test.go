package test_under_appengine

import (
	"bytes"
	"image"
	"log"
	"testing"

	"github.com/pbberlin/tools/conv"
	"github.com/pbberlin/tools/util"
)

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

func TestStringToImgAndBack(t *testing.T) {

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
