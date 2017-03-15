package dashboard

import (
	"html/template"
	"log"
	"mime"
	"net/http"
	"testing"

	"path"

	"io/ioutil"

	"github.com/julienschmidt/httprouter"
)

var templatePath string
var templateFiles []string

func init() {
	templatePath = "template/asset"
	templateFiles = []string{
		"template/public_html/index.html",
		"template/public_html/common/header.html",
		"template/public_html/common/footer.html",
		"template/public_html/common/left.html",
	}

}

func TestTemplate(t *testing.T) {
	router := httprouter.New()
	router.GET("/", IndexHandler)
	router.GET("/asset/:path/:asset", staticHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func staticHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	assetName := ps.ByName("asset")
	pathName := ps.ByName("path")
	assetFile := path.Join(templatePath, pathName, assetName)
	log.Println(assetFile)
	data, _ := ioutil.ReadFile(assetFile)

	ext := path.Ext(assetName)
	ct := mime.TypeByExtension(ext)
	if ct != "" {
		w.Header().Set("Content-Type", ct)
	}

	w.Write(data)
}

func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := map[string]string{
		"Name": "Mike",
	}

	tmpl := template.Must(template.ParseFiles("template/public_html/index.html", "template/public_html/common/header.html", "template/public_html/common/footer.html", "template/public_html/common/left.html"))
	w.Header().Set("Content-Type", "text/html")

	err := tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}
