package test_under_appengine

import (
	"log"
	"os"
	"testing"

	"github.com/pbberlin/tools/conv"
	"github.com/pbberlin/tools/dsu"

	"appengine/aetest"
)

var str_b64 string = testString()
var debug bool = false

func TestPutGet(t *testing.T) {

	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Errorf("could not get a context")
		log.Printf("could not get a context")
		os.Exit(1)
	}
	defer c.Close()

	vVBytes, _ := conv.String_to_VVByte(str_b64)
	wb := dsu.WrapBlob{}
	wb.Name = "test"
	wb.VVByte = vVBytes
	key_combi, _ := dsu.BufPut(c, wb, "test")

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

type SomeStruct struct {
	S1 string `json:"s1"`
	S2 string `json:"s2"`
}

func Test_meMcacheGet_set(t *testing.T) {

	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Errorf("could not get a context")
		log.Printf("could not get a context")
		os.Exit(1)
	}
	defer c.Close()

	rvb1 := SomeStruct{}
	ok := dsu.McacheGet(c, "key1", &rvb1)
	if debug {
		log.Printf("  -------%#v--%T--%v--- \n\n", rvb1, rvb1, ok)
	}

	rvb2 := ""
	ok2 := dsu.McacheGet(c, "key2", "otherstr")
	if debug {
		log.Printf("  -------%#v--%T--%v--- \n\n", rvb2, rvb2, ok2)
	}

	rvb3 := 2323
	ok3 := dsu.McacheGet(c, "key3", 22323)
	if debug {
		log.Printf("  -------%#v--%T--%v--- \n\n", rvb3, rvb3, ok3)
	}

	dsu.McacheSet(c, "key1", "just a scalar stupid string")

	myStruct1 := SomeStruct{"this content", "is structured"}
	dsu.McacheSet(c, "key2", myStruct1)

	myStruct2 := SomeStruct{"wonderbar", "is not wonderbra"}
	dsu.McacheSet(c, "key3", myStruct2)

	// rva2 := 22323
	// ok2b := dsu.McacheGet(c, "key1", 22323)
	// if debug {
	// 	log.Printf("  -------%#v--%T--%v--- \n\n", rva2, rva2, ok2b)
	// }

	// rva3 := "string"
	// ok3b := dsu.McacheGet(c, "key2", "string")
	// if debug {
	// 	log.Printf("  -------%#v--%T--%v--- \n\n", rva3, rva3, ok3b)
	// }

	rva4 := SomeStruct{}
	ok4b := dsu.McacheGet(c, "key3", &rva4)
	if debug {
		log.Printf("  -------%#v--%T--%v--- \n\n", rva4, rva4, ok4b)
	}

	rva5 := SomeStruct{}
	dsu.McacheGet(c, "key2", &rva5)
	if debug {
		log.Printf("  --%#v--- \n\n", rva5)
	}

	// want1 := "just a scalar stupid string"
	// if rva1 != want1 {
	// 	t.Errorf("memache get - set - want %s, got ", want1, rva1)
	// }

	// want2 := 0
	// if rva2 != want2 {
	// 	t.Errorf("memache get - set - want %s, got ", want2, rva2)
	// }

	// want3 := "{\"s1\":\"this content\",\"s2\":\"is structured\"}"
	// if rva3 != want3 {
	// 	t.Errorf("memache get - set - want %s, got ", want3, rva3)
	// }

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
