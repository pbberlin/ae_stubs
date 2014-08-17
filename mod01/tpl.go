package main

import (
	"html"
	tt "html/template"    
	"net/http"    
	"time"
	"github.com/pbberlin/tools/u_err"

)

/*
	two prefixes 
	   c_ ... ist the string content of a template
	   n_ ... ist the string NAME    of a template

*/








const c_page_scaffold_01 = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <link rel="icon" href="data:;base64,=">
    <title>{{template "static_title" }}</title>
  </head>
  <body>
    {{template "n_content" .}}
    {{template "static_content_1" }}
    {{template "static_content_2" }}
    <p><a href='/'>Back to root</a></p>
  </body>
</html>
`



const c_content_1 = `
	{{define "n_content"}}
	   --{{.}}--
	{{end}}
`

const c_content_2 = `
	{{define "n_content"}}
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
	{{end}}

`


const c_guestEntryHTML = `
{{define "n_content"}}
    <form action="/guest-save" method="post">
      <div><textarea name="content" rows="3" cols="60"></textarea></div>
      <div><input type="submit" value="Save Entry"></div>
    </form>
{{end}}
`



const c_formFetchURL = `
{{define "n_content"}}
    <form action="/fetch-url" method="post">
      <div><input name="url"    size="160"  value="{{.}}"></div>
      <div><input type="submit" value="Fetch" accesskey='f'></div>
    </form>
{{end}}
`

// 		

var t_base *tt.Template = nil


func cloneFromBase() *tt.Template {

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
	util_err.Err_log(err)

	return t_derived	
}


func templatesExtend(  m map[string]string , str_dyn_content string) *tt.Template {
	var err error = nil
	tder  := cloneFromBase()
	
	
	tder, err = tder.Parse( str_dyn_content )


	util_err.Err_log(err)

	for k,v := range m{
		tder, err = tder.Parse( `{{define "` + k  +`"}}`   + v + `{{end}}` )
		util_err.Err_log(err)
	}

	return tder
}


/*
	application example
	
	
	g := Greeting1{
		Content: "contnt contnt contnt contnt contnt contnt contnt contnt ",
		Author:   "dooodle",
		Date:    time.Now(),
	}
	vg := []Greeting1{g,g}

	mc  := map[string]string{
		"static_title"  :   "second title",
		"static_content_1":     "<pre>" + report + "</pre>",
		"static_content_2":     s_cntr ,
	}
	myTplExecute(w,mc,c_content_2, gbEntries)




*/


func myTplExecute( w http.ResponseWriter, m map[string]string , str_dyn_content string, v interface{}  ) {

	t_1  :=  templatesExtend( m, str_dyn_content )
	_ = t_1
	err  :=  t_1.ExecuteTemplate(w, "n_page_scaffold_01", v)
	util_err.Err_log(err)

}

