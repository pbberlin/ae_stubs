package main

import (
	"html"
	tt "html/template"    
	"net/http"    
	"github.com/pbberlin/tools/util_err"
	
)







const T0 = `
I am T0  --{{template "T1"}}--<br>
`

const T1 = `
    {{define "T1"}}
    	|| 
    		<span style="color:#f22;">I am T1</span> 
			<span style="font-weight:bold;">{{template "T2"}}</span>
		||
	{{end}}
`




func templatesCompileDemo( w http.ResponseWriter ) {

	funcMap := tt.FuncMap{ 
		"unescape": html.UnescapeString, 
		"escape":   html.EscapeString,
	}	


	var t_base *tt.Template
	var err error = nil

	t_base = tt.Must(tt.New("str_t_outmost").Funcs(funcMap).Parse(T0))
	util_err.Err_log(err)

   t_base , err = t_base.Parse(T1)
	util_err.Err_log(err)


	t_1, err := t_base.Clone()
	util_err.Err_log(err)

	t_2, err := t_base.Clone()
	util_err.Err_log(err)

	t_1, err = t_1.Parse("{{define `T2`}}T2, version A{{end}}")
	util_err.Err_log(err)


	t_2, err = t_2.Parse("{{define `T2`}}T2, version B{{end}}")
	util_err.Err_log(err)


	err = t_1.ExecuteTemplate(w, "str_t_outmost", nil)
	util_err.Err_log(err)

	err = t_2.ExecuteTemplate(w, "str_t_outmost", nil)
	util_err.Err_log(err)

   
}






