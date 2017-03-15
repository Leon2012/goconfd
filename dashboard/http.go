package dashboard

import (
	"html/template"
	"io/ioutil"
	"mime"
	"net"
	"net/http"
	"path"
	"path/filepath"

	"strings"

	"github.com/Leon2012/goconfd/store/types"
	"github.com/julienschmidt/httprouter"
)

type httpServer struct {
	Path            string
	HtmlPath        string
	AssetPath       string
	HtmlCommonFiles []string
	ctx             *Context
	router          http.Handler
}

func newHttpServer(c *Context, path string) *httpServer {
	router := httprouter.New()
	router.HandleMethodNotAllowed = false
	h := &httpServer{}
	h.Path = path
	h.AssetPath = filepath.Join(h.Path, "asset")
	h.HtmlPath = filepath.Join(h.Path, "public_html")
	h.HtmlCommonFiles = []string{
		filepath.Join(h.HtmlPath, "common/header.html"),
		filepath.Join(h.HtmlPath, "common/footer.html"),
		filepath.Join(h.HtmlPath, "common/left.html"),
		filepath.Join(h.HtmlPath, "common/layout.html"),
	}
	h.ctx = c
	h.router = router
	router.Handle("GET", "/", h.Index)
	router.Handle("GET", "/agent", h.Agent)
	router.Handle("GET", "/heartbeat/:hostname/:keyprefix", h.Heartbeat)
	router.Handle("GET", "/var", h.Var)
	router.Handle("GET", "/add/:key/:value", h.Add)
	router.Handle("GET", "/system", h.System)
	router.GET("/asset/:path/:asset", h.Static)
	return h
}

func (h *httpServer) serve(listener net.Listener) {
	h.ctx.Dashboard.logf("http server listening on %s", listener.Addr())
	server := &http.Server{
		Handler: h,
		//ErrorLog: h.ctx.Agent.opts.Logger,
	}
	err := server.Serve(listener)
	if err != nil {
		h.ctx.Dashboard.logf("ERROR: http.Server() - %s", err)
	}
	h.ctx.Dashboard.logf("http server closing %s", listener.Addr())
}

func (h *httpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if req.Method == "POST" {
		h.DoAdd(w, req)
	} else {
		h.router.ServeHTTP(w, req)
	}
}

func (h *httpServer) Render(templateFile string, data interface{}, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	templateFile += ".html"
	templateFile = filepath.Join(h.HtmlPath, templateFile)
	templateFiles := []string{}
	templateFiles = append(templateFiles, templateFile)
	for _, htmlCommonFile := range h.HtmlCommonFiles {
		templateFiles = append(templateFiles, htmlCommonFile)
	}
	tmpl := template.Must(template.ParseFiles(templateFiles...))
	err := tmpl.Execute(w, data)
	if err != nil {
		h.ctx.Dashboard.logf("FATAL: execute template %s", err.Error())
	}
}

func (h *httpServer) DoAdd(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimSpace(r.FormValue("key"))
	value := strings.TrimSpace(r.FormValue("value"))
	if key != "" && value != "" {
		err := h.ctx.Dashboard.idc.Put(key, value)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	http.Redirect(w, r, "/var", 302)
}

func (h *httpServer) Add(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	key := strings.TrimSpace(ps.ByName("key"))
	value := strings.TrimSpace(ps.ByName("value"))
	if key != "" && value != "" {
		err := h.ctx.Dashboard.idc.Put(key, value)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		w.Write([]byte("OK\n"))
	}
}

func (h *httpServer) Var(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	templateFile := "var"
	h.Render(templateFile, nil, w, r, ps)
}

func (h *httpServer) System(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (h *httpServer) Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	templateFile := "index"
	data := make(map[string]string)
	data["title"] = "goconfd dashboard"
	h.Render(templateFile, data, w, r, ps)
}

func (h *httpServer) Agent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	templateFile := "agent"
	data, err := h.ctx.Dashboard.db.GetAgents()
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		h.Render(templateFile, data, w, r, ps)
	}
}

func (h *httpServer) Heartbeat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var data []*types.Heartbeat
	var err error
	templateFile := "heartbeat"
	hostName := ps.ByName("hostname")
	keyPrefix := ps.ByName("keyprefix")
	if hostName != "" && keyPrefix != "" {
		agent := &types.Agent{}
		agent.HostName = hostName
		agent.KeyPrefix = keyPrefix
		data, err = h.ctx.Dashboard.db.GetHeartbeatsByAgent(agent)
	} else {
		data, err = h.ctx.Dashboard.db.GetHeartbeats()
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		h.Render(templateFile, data, w, r, ps)
	}
}

func (h *httpServer) Static(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	assetName := ps.ByName("asset")
	pathName := ps.ByName("path")
	assetFile := filepath.Join(h.AssetPath, pathName, assetName)
	//log.Println(assetFile)
	data, err := ioutil.ReadFile(assetFile)
	if err != nil {
		http.Error(w, err.Error(), 404)
	} else {
		ext := path.Ext(assetName)
		ct := mime.TypeByExtension(ext)
		if ct != "" {
			w.Header().Set("Content-Type", ct)
		}
		w.Write(data)
	}
}
