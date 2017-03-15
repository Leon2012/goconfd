package agent

import (
	"fmt"
	"net"
	"net/http"
	_ "strings"

	"github.com/Leon2012/goconfd/libs/kv"
	"github.com/julienschmidt/httprouter"
)

type httpServer struct {
	ctx    *Context
	router http.Handler
}

func newHttpServer(ctx *Context) *httpServer {
	router := httprouter.New()
	router.HandleMethodNotAllowed = true
	s := &httpServer{
		ctx:    ctx,
		router: router,
	}
	router.Handle("GET", "/ping", s.Ping)
	router.Handle("GET", "/get/:key", s.DoGet)
	router.Handle("GET", "/info", s.DoInfo)
	return s
}

func (h *httpServer) serve(listener net.Listener) {
	h.ctx.Agent.logf("http server listening on %s", listener.Addr())
	server := &http.Server{
		Handler: h,
		//ErrorLog: h.ctx.Agent.opts.Logger,
	}
	err := server.Serve(listener)
	if err != nil {
		h.ctx.Agent.logf("ERROR: http.Server() - %s", err)
	}
	h.ctx.Agent.logf("http server closing %s", listener.Addr())
}

func (h *httpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.router.ServeHTTP(w, req)
}

func (h *httpServer) Ping(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "Ok\n")
}

func (h *httpServer) DoInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	m := make(map[string]string)
	if h.ctx.Agent.lastHeartbeat != nil {
		m["lastUpdateKey"] = h.ctx.Agent.lastHeartbeat.Kv.Key
		m["lastUpdateValue"] = h.ctx.Agent.lastHeartbeat.Kv.Value
		m["lastUpdateTime"] = h.ctx.Agent.lastHeartbeat.UpdateTime.Format("2006-01-02 15:04:05")
	}
	fmt.Fprint(w, m)
}

func (h *httpServer) DoGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	key := ps.ByName("key")
	k, err := LoadValueByKey(h.ctx, key)
	if err != nil {
		http.Error(w, err.Error(), 404)
	} else {
		data, err := kv.JsonEncode(k)
		if err != nil {
			http.Error(w, err.Error(), 404)
		} else {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, string(data))
		}
	}
}
