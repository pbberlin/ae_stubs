package blobstore_image_resize

import (
	"html/template"
	"net/http"

	"appengine"
	"appengine/blobstore"
	
	"fmt"
	
	"appengine/image"
	
	"github.com/pbberlin/tools/parsetools"

	"github.com/pbberlin/tools/u_err"

	"time"
	"bytes"
	"appengine/datastore"
	"strings"
	"unicode/utf8"
	
	"path"
)


type BlobInfo struct{
	Content_type string		`datastore:"content_type"`
	Creation time.Time		`datastore:"creation"`
	Filename string			`datastore:"filename"`
	Md5_hash string			`datastore:"md5_hash"`
	Size int                `datastore:"size"`
}




var rootTemplate = template.Must(template.New("root").Parse(rootTemplateHTML))


const rootTemplateHTML = `<html>
	<body>
		<form action="{{.}}" method="POST" enctype="multipart/form-data">
			Upload File: <input type="file" name="file"><br>
			<input type="submit" name="submit" value="Submit">
		</form>
	</body>
</html>
`



func blobUpload(w http.ResponseWriter, r *http.Request) {

	parsetools.SplitByWhitespace("a b")

	c := appengine.NewContext(r)
	uploadURL, err := blobstore.UploadURL(c, "/blob/server-process", nil)
	if err != nil {
		util_err.ServeError(c, w, err)
		return
	}
	
	w.Header().Set("Content-Type", "text/html")
	err = rootTemplate.Execute(w, uploadURL)
	if err != nil {
		c.Errorf("%v", err)
	}
}


func blobServerProcess(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	blobs, _, err := blobstore.ParseUpload(r)
	if err != nil {
		w.Write(  []byte("<a href='/blob/upload' >Fehler beim Parsing</a>")  )
		//util_err.ServeError(c, w, err)
		return
	} 
	
	file   := blobs["file"]
	lFile  := len(file)      // this always yields (int)1 
	slFile := fmt.Sprintf("-%v-  %v",lFile, file[:])
	if lFile == 0 {
		c.Errorf("no file uploaded")
		//http.Redirect(w, r, "/blob/upload", http.StatusFound)
		w.Write(  []byte("<a href='/blob/upload' >Keine oder leere Datei, Try again</a>")  )
		return
	}

	w.Header().Set("Content-Type", "text/html")
	file[0].Filename = strings.ToLower(file[0].Filename)
	w.Write(  []byte( "filename is " + file[0].Filename +"<br>\n")  )
	
	
	
	urlSuccessX     := "/blob/serve?blobkey="+string(file[0].BlobKey)
	urlSuccessThumb := "/blob/thumb-serve?blobkey="+string(file[0].BlobKey)
	w.Write(  []byte("<a href='"+urlSuccessX+"' >Erfolg ("+slFile+"Bytes)- view it</a><br>\n")  )
	w.Write(  []byte("<a href='"+urlSuccessThumb+"' >Erfolg ("+slFile+"Bytes)- view Thumbnail</a><br>\n")  )
	w.Write(  []byte("<a href='/blob/list' >List</a>")  )
	
	//http.Redirect(w, r, urlSuccess, http.StatusFound)
}


func blobServe(w http.ResponseWriter, r *http.Request) {
	blobstore.Send(w, appengine.BlobKey(r.FormValue("blobkey")))
}

func blobThumbServe(w http.ResponseWriter, r *http.Request) {


		c := appengine.NewContext(r)
		k := appengine.BlobKey(r.FormValue("blobkey"))
		
		var o image.ServingURLOptions = *new(image.ServingURLOptions )
		o.Size = 200
		o.Crop = true
		url,err := image.ServingURL(c,k,&o)

		if err != nil {
			util_err.ServeError(c, w, err)
			return
		}
		
		http.Redirect(w, r, url.String(), http.StatusFound)		
}


// working differently as in Python
//		//blobs := blobstore.BlobInfo.gql("ORDER BY creation DESC")
func blobList(w http.ResponseWriter, r *http.Request) {

	b1 := new(bytes.Buffer)
	s1 := ""

	s1 = `
		<style>
			.ib {
				vertical-align:middle;
				display:inline-block;
				width:95px;
			}
		</style>
	`
	b1.WriteString( s1 )


	namefilter := r.FormValue("name")


	c := appengine.NewContext(r)
	q := datastore.NewQuery("__BlobInfo__")
	if namefilter != "" {


		namefilter = strings.ToLower(namefilter)
		b1.WriteString( fmt.Sprintf("Filtering by %v<br>",namefilter) )

		// we filter for the first char of namefilter
		// i.E. "Peter" =>  searches for all files >=p and <q 

		
		vb := []byte(namefilter)
		rl_1, rsize := utf8.DecodeRune(vb)
		_ = rsize
		rl_2 := rl_1 +1
		//c.Infof("next utf8/ascii code of first char %v  (%c) - bytes: %v\n", rl_2 ,rl_2 , rsize)
		q = datastore.NewQuery("__BlobInfo__").Filter("filename>=",fmt.Sprintf("%c",rl_1))
		q = q.Filter("filename<=", fmt.Sprintf("%c",rl_2) )

		/*
		// this would be the complementary query for upper case letters
		// >=P and <Q - but we convert all to lower case 
		ru_1 :=  rl_1 - 'A' + 'a'
		ru_2  := ru_1 + 1
		q = q.Filter("filename>=", fmt.Sprintf("%c",ru_1) )
		q = q.Filter("filename<=", fmt.Sprintf("%c",ru_2)  )

		*/
		
	}  
	for t := q.Run(c); ; {
		var bi BlobInfo
		dsKey, err := t.Next(&bi)
		
		if err == datastore.Done {
			c.Infof("   No Results (any more) blob-list %v", err)
			break
		}
		// other err
		if err != nil {
			util_err.Err_log(err)
			return 
		}

		
		//s1 = fmt.Sprintf("key %v %v %v %v %v %v<br>\n", dsKey.AppID(),dsKey.Namespace() , dsKey.Parent(), dsKey.Kind(), dsKey.StringID(), dsKey.IntID())
		//b1.WriteString( s1 )
		

		//s1 = fmt.Sprintf("blobinfo: %v %v<br>\n", bi.Filename, bi.Size)
		//b1.WriteString( s1 )


		ext  := path.Ext( bi.Filename)
		base := path.Base(bi.Filename)
		base = base[:len(base)-len(ext)]
		
		//b1.WriteString( fmt.Sprintf("-%v-  -%v-",base, ext) )

		base = strings.Replace(base ,"_", " ", -1)
		base = strings.Title(base)
		ext  = strings.ToLower(ext)

		titledFilename := base + ext

		s1 = fmt.Sprintf("<a class='ib' style='width:280px;margin-right:20px' target='_view' href='/blob/serve?blobkey=%v'>%v</a> &nbsp; &nbsp; \n", dsKey.StringID(), titledFilename)
		b1.WriteString( s1 )

		if bi.Content_type == "image/png"  || bi.Content_type == "image/jpeg"   {
			s1 = fmt.Sprintf("<img class='ib' style='width:40px;' src='/_ah/img/%v%v' />\n", 
				dsKey.StringID(),"=s200-c")
			b1.WriteString( s1 )			

			s1 = fmt.Sprintf("<a class='ib' target='_view' href='/_ah/img/%v%v'>Thumb</a>\n", 
				dsKey.StringID(),"=s200-c")
			b1.WriteString( s1 )

		} else {
			s1 = fmt.Sprintf("<span class='ib' style='width:145px;'> &nbsp; no thb</span>")
			b1.WriteString( s1 )			
			
		}

		
		

		s1 = fmt.Sprintf("<a class='ib' target='_rename_delete' href='/blob/rename-delete?action=delete&blobkey=%v'>Delete</a>\n", 
			dsKey.StringID())
		b1.WriteString( s1 )


		s1 = fmt.Sprintf(`
			<span class='ib' style='width:450px; border: 1px solid #aaa'>
				<form target='_rename_delete' action='/blob/rename-delete' >
					<input name='blobkey'  value='%v'     type='hidden'/>
					<input name='action'   value='rename' type='hidden'/>
					<input name='filename' value='%v' size='42'/>
					<input type='submit'   value='Rename' />
				</form>
			</span>
			`, dsKey.StringID(), bi.Filename)
		b1.WriteString( s1 )

		b1.WriteString( "<br><br>\n\n" )
		
	}

	b1.WriteString("<a accesskey='u' href='/blob/upload' >Upload</a>")
	
	w.Header().Set("Content-Type", "text/html")
	w.Write( b1.Bytes() )
	
}



func blobRenameDelete(w http.ResponseWriter, r *http.Request) {

	b1 := new(bytes.Buffer)
	s1 := ""

	defer func(){
		w.Header().Set("Content-Type", "text/html")
		w.Write( b1.Bytes() )		
	}()

	c := appengine.NewContext(r)

	bk := r.FormValue("blobkey")
	if bk == "" {  
		b1.WriteString( "No blob key given<br>" )		
		return
	} else {
		s1 = fmt.Sprintf("Blob key given %q<br>", bk)
		b1.WriteString( s1 )				
	}
	
	
		

	dsKey := datastore.NewKey(c, "__BlobInfo__", bk, 0, nil)
	

	q := datastore.NewQuery("__BlobInfo__").Filter("__key__=",dsKey)

	var bi BlobInfo 
	var found bool

	for t := q.Run(c); ; {
		_, err := t.Next(&bi)
		
		if err == datastore.Done {
			c.Infof("   No Results (any more), blob-rename-delete %v", err)
			break
		}
		// other err
		if err != nil {
			util_err.Err_log(err)
			return 
		}

		found = true
		break
		
	}

	if found {

		ac := r.FormValue("action")
		
		if ac == "delete" {
			b1.WriteString( "deletion  " )				


			// first the binary data
			keyBlob, err := blobstore.BlobKeyForFile(c , bi.Filename )
			util_err.Err_log(err)
			
			if err != nil {
				b1.WriteString(   fmt.Sprintf(" ... failed (1) %v", err ) )				
			} else {
				err = blobstore.Delete(c , keyBlob) 
				util_err.Err_log(err)
	
				if err != nil {
					b1.WriteString(   fmt.Sprintf(" ... failed (2) %v", err ) )				
				} else {
					c.Infof("got a strange blobstore key %#q, %T",keyBlob,keyBlob)

					// now the datastore record		
					err = datastore.Delete(c , dsKey) 
					util_err.Err_log(err)

					if err != nil {
						b1.WriteString(   fmt.Sprintf(" ... failed (3) %v", err ) )				
					} else {
						b1.WriteString(   " ... succeeded<br>") 
					}
					
				}
			}
		}

		if ac == "rename" {
			b1.WriteString( "renaming " )				
			
			nfn := r.FormValue("filename")
			if nfn == "" || len(nfn) <4 {
				b1.WriteString(   " ... failed - at LEAST 4 chars required<br>") 
				return
			}
			nfn  = strings.ToLower(nfn)
			
			bi.Filename = nfn
			_, err := datastore.Put(c , dsKey,  &bi) 
			util_err.Err_log(err)
			if err != nil {
				b1.WriteString(   fmt.Sprintf(" ... failed. %v", err ) )				
			} else {
				b1.WriteString(   " ... succeeded<br>") 
			}
		}
		
		
	} else {
		b1.WriteString( "no blob found for given blobkey<br>" )				
	}

	
}







func blobHome(w http.ResponseWriter, r *http.Request) {
	

	b1 := new(bytes.Buffer)

	
	b1.WriteString( "<a target='_blob' href='/blob/upload'>Upload new Blob</a><br>\n")
	b1.WriteString( "<a target='_blob' href='/blob/list'  >Blob List</a><br>\n")
	
	
	
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")	
	w.Write( b1.Bytes() ) 
	
}


func init() {
	http.HandleFunc("/blob"  , blobHome)
	http.HandleFunc("/blob/"  , blobHome)
	http.HandleFunc("/blob/upload"  , blobUpload)
	http.HandleFunc("/blob/server-process"  , blobServerProcess)
	http.HandleFunc("/blob/serve", blobServe)
	http.HandleFunc("/blob/rename-delete", blobRenameDelete)
	http.HandleFunc("/blob/thumb-serve", blobThumbServe)
	http.HandleFunc("/blob/list", blobList)
}
