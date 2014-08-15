package blobstore_image_resize

import (
	"html/template"
	"net/http"

	"appengine"
	"appengine/blobstore"
	
	"fmt"
	
	"appengine/image"
	
	"github.com/pbberlin/tools/parsetools"

	"github.com/pbberlin/tools/util_err"


)







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
	uploadURL, err := blobstore.UploadURL(c, "/blob-server-process", nil)
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
		w.Write(  []byte("<a href='/blob-upload' >Fehler beim Parsing</a>")  )
		//util_err.ServeError(c, w, err)
		return
	} 
	
	
	
	file   := blobs["file"]
	lFile  := len(file)      // this always yields (int)1 
	slFile := fmt.Sprintf("-%v-  %v",lFile, file[:])
	if lFile == 0 {
		c.Errorf("no file uploaded")
		//http.Redirect(w, r, "/blob-upload", http.StatusFound)
		w.Write(  []byte("<a href='/blob-upload' >Keine oder leere Datei, Try again</a>")  )
		return
	}
	urlSuccessX := "/blob-respond/?blobKey="+string(file[0].BlobKey)
	urlSuccessThumb := "/blob-image-thumb/?blobKey="+string(file[0].BlobKey)
	w.Write(  []byte("<a href='"+urlSuccessX+"' >Erfolg ("+slFile+"Bytes)- view it</a><br>\n")  )
	w.Write(  []byte("<a href='"+urlSuccessThumb+"' >Erfolg ("+slFile+"Bytes)- view Thumbnail</a>")  )
	
	//http.Redirect(w, r, urlSuccess, http.StatusFound)
}


func blobRespond(w http.ResponseWriter, r *http.Request) {
	blobstore.Send(w, appengine.BlobKey(r.FormValue("blobKey")))
}

func blobImageThumb(w http.ResponseWriter, r *http.Request) {


		c := appengine.NewContext(r)
		k := appengine.BlobKey(r.FormValue("blobKey"))
		
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

func blobList(w http.ResponseWriter, r *http.Request) {

	// NOT WORKING - opposed to Python
	//blobs := blobstore.BlobInfo.gql("ORDER BY creation DESC")

	//for blob in blobs 
}

func init() {
	http.HandleFunc("/blob-upload"  , blobUpload)
	http.HandleFunc("/blob-server-process"  , blobServerProcess)
	http.HandleFunc("/blob-respond/", blobRespond)
	http.HandleFunc("/blob-image-thumb/", blobImageThumb)
	http.HandleFunc("/blob-list/", blobList)
}

