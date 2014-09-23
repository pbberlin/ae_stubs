package big_query


import (
	"net/http"	
	"tpl_html"
	"github.com/pbberlin/tools/adapter"
	
)




const c_tpl_gbentry = `
		{{range .}}
			- <b>{{.Author}}</b> wrote: 
			{{.Content}} <br>
		{{end}}
`


func viewHTML( w http.ResponseWriter , r *http.Request, m map[string]interface{}) {

	b1,ml := disLegend(w,r)

	_ = ml
	_ = b1

	add, tplExec := tpl_html.FuncTplBuilder(w,r)
	
	add("n_html_title","The Battle of Computer Languages","")  

	add("n_cont_0"  , tpl_html.PrefixLff + "chart_body", map[string]map[string]string{"legend":ml} )
	add("tpl_legend", tpl_html.PrefixLff + "chart_body_embed01", "" )
	
	tplExec(w,r)
	
}


func init() {
	http.HandleFunc("/big-query/html", adapter.Adapter(viewHTML) )
}
