package main

import (

	"html"
	tt "html/template"    
	"net/http"    
	"time"
	"github.com/pbberlin/tools/util_err"

)

/*
two prefixes 
   n_ ... ist the string NAME    of a template
   c_ ... ist the string CONTENT of a template


example

type GBEntry struct {
    Author  string
    Content string      
}

const c_tpl_gbe = `
		{{range .}}
			<b>{{.Author}}</b> wrote: 
			{{.Content}} <br>
		{{end}}
`

gbe1 := GBEntry{
	Content: "gb entry contnt. gb entry contnt. gb entry contnt.",
	Author:   "John Dos Passos",
}

vgbe := []GBEntry{gbe1,gbe1}


mc  := map[string]string{
	"static_title"  :   "second title",
	"n_cont_0"     :   c_tpl_gbe,

}
myTplExecute(w,mc,vgbe)


*/




const c_page_scaffold_01 = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <link rel="icon" href="data:;base64,=">
    <title>{{template "static_title"  }}</title>
  </head>
  <body>
    {{template "n_cont_0" .n_cont_0}}
    {{template "n_cont_1" .n_cont_1 }}
    {{template "n_cont_2" .n_cont_2 }}
    <p><a href='/'>Back to root</a></p>
  </body>
</html>
`


// must contain all subtemplates demanded by c_page_scaffold_01 
var map_default map[string]string = map[string]string{
	"static_title"    :     "",
	"n_cont_0":     "",
	"n_cont_1":     "",
	"n_cont_2":     "",
}



const c_contentpl_extended = `
	--{{.}}--
`

const c_view_gbe = `
	{{range .}}
		{{$a := .Date}}
		{{$b := .Date  | df }}
		{{$c := df .Date}}
			<p>
		{{with .Author}}
			<b>{{.}}</b> wrote on {{$c}}<br>
		{{else}}
			An anonymous person wrote:   <br>
		{{end}}
			<span style='display:block; max-width:300px;font-size:12px;' >{{.Content}}</span>
		</p>
	{{end}}
`


const c_new_gbe = `
	<form action="/guest-save" method="post">
		<div><textarea name="content" rows="3" cols="60"></textarea></div>
		<div><input type="submit" value="Save Entry"></div>
	</form>
`



const c_formFetchURL = `

    <form action="/fetch-url" method="post">
      <div><input name="url"    size="160"  value="{{.}}"></div>
      <div><input type="submit" value="Fetch" accesskey='f'></div>
    </form>

`

// 		

var t_base *tt.Template = nil


func cloneFromBase(w http.ResponseWriter, r *http.Request) *tt.Template {

	funcMap := tt.FuncMap{ 
		"unescape": html.UnescapeString, 
		"escape"  : html.EscapeString,
		"df"      : func(g time.Time) string{
			return g.Format( "2006-01-02 (Jan 02)" ) 
		},
	}	

	if t_base == nil {
		t_base = tt.Must(tt.New("n_page_scaffold_01").Funcs(funcMap).Parse( c_page_scaffold_01 ))
	}

	t_derived, err := t_base.Clone()
	util_err.Err_http(w,r,err,false)

	return t_derived	
}


func templatesExtend(w http.ResponseWriter, r *http.Request,  m map[string]string ) *tt.Template {

	var err error = nil
	tder  := cloneFromBase(w,r)
	

	for k,v := range m{
		tder, err = tder.Parse( `{{define "` + k  +`"}}`   + v + `{{end}}` )
		util_err.Err_http(w,r,err,false)
	}

	return tder
}




func funcTplBuilder(w http.ResponseWriter, r *http.Request)( f1  func( string, string, interface{} ) , 
  f2 func(http.ResponseWriter, *http.Request) ){

	// map template contents, map template data
	mtc := map[string]string      {}
	mtd := map[string]interface{} {}

	// template key - template content, template data
	f1 = func( tk string, tc string, td interface{} ) {

		_,ok := map_default[tk]
		util_err.Err_http(w,r,ok,false,"template key must be one of " , map_default )

		mtc[tk] = tc
		mtd[tk] = td

	}
	
	
	f2 = func(w http.ResponseWriter, r *http.Request){

		// merge arguments with defaults
		map_result  := map[string]string{}
		for k,v := range map_default {
			if _,ok := mtc[k]; ok {
				map_result[k] = mtc[k]			
			} else {
				map_result[k] = v
			}
		}
	
		tpl_extended  :=  templatesExtend( w,r,map_result )
	
		err  :=  tpl_extended.ExecuteTemplate(w, "n_page_scaffold_01", mtd)
		util_err.Err_http(w,r,err,false)


	}
	
	return f1,f2
}

// myTemplateAdder, myTplExec := funcTplBuilder()
// myTemplateAdder("n_content","--{{.}}--","some string")
// myTplExec(w,r)


