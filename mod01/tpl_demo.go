package main

import (
	"html"
	tt "html/template"	
	"net/http"	
	"github.com/pbberlin/tools/u_err"
	
)







const T0 = `
	T0	<br>
		<span style='color:#aaf; line-height:200%;display:inline-block; margin:8px; margin-left:120px; border: 1px solid #aaf'>
				{{template "T1" .}}
		</span>
		<hr>
`

const T1 = `
{{define "T1"}}
	T1 <br>
		--{{.key1}}--<br>
		<span style='color:#faa; display:inline-block; margin:8px; margin-left:120px; border: 1px solid #faa'>
			{{template "T2" .key2 }}
		</span>
	
{{end}}
`


const iterOver = `{{ $mapOrArray := . }} 
{{range $index, $element := $mapOrArray }}
   <li><strong>$index</strong>: $element </li>
{{end}}`

const treatFirstIterDifferent = `{{if $index}},{{end}}`


func templatesCompileDemo( w http.ResponseWriter , r *http.Request, m map[string]interface{}) {

	w.Header().Set("Content-Type", "text/html")


	funcMap := tt.FuncMap{ 
		"unescape": html.UnescapeString, 
		"escape":   html.EscapeString,
	}	


	var t_base *tt.Template
	var err error = nil

	// creating T0 - naming it - adding func map
	t_base = tt.Must(tt.New("str_T0_outmost").Funcs(funcMap).Parse(T0))
	util_err.Err_http(w,r,err)

	// adding T1 definition
   t_base , err = t_base.Parse(T1)  // definitions must appear at top level - but not at the start of
	util_err.Err_http(w,r,err)
	


	// create two clones 
	// both contain T0 and T1
	tc_1, err := t_base.Clone()
	util_err.Err_http(w,r,err)
	tc_2, err := t_base.Clone()
	util_err.Err_http(w,r,err)


	// adding different T2 definitions
	tc_1, err = tc_1.Parse("{{define `T2`}}T2-A  <br>--{{.}}--  {{end}}")
	util_err.Err_http(w,r,err)
	tc_2, err = tc_2.Parse("{{define `T2`}}T2-B  <br>--{{.}}--  {{end}}")
	util_err.Err_http(w,r,err)



	// writing both clones to the response writer
	err = tc_1.ExecuteTemplate(w, "str_T0_outmost", nil)
	util_err.Err_http(w,r,err)

	// second clone is written with dynamic data on two levels
	dyndata := map[string]string{"key1":"dyn_val1","key2":"dyn_val2"}
	err = tc_2.ExecuteTemplate(w, "str_T0_outmost", dyndata)
	util_err.Err_http(w,r,err)

	// Note: it is important to pass the . 
	//		 {{template "T1" .}}
	//		 {{template "T2" .key2 }}
	//						 ^
	// otherwise "dyndata" can not be accessed by the inner templates...
 
 
 	// leaving T2 undefined => error 
	tc_3, err := t_base.Clone()
	util_err.Err_http(w,r,err)
	err = tc_3.ExecuteTemplate(w, "str_T0_outmost", dyndata)
	util_err.Err_http(w,r,err)

  
}


func init() {
	http.HandleFunc("/tpl/demo", adapter(templatesCompileDemo) )
}
