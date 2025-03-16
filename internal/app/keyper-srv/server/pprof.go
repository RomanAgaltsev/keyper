package server

import (
	"net/http"
	"net/http/pprof"
)

func NewPprofServer(Addr string) *http.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/debug/pprof/", pprof.Index)
	handler.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	handler.HandleFunc("/debug/pprof/profile", pprof.Profile)
	handler.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	handler.HandleFunc("/debug/pprof/trace", pprof.Trace)

	server := &http.Server{
		Addr:    Addr,
		Handler: handler,
	}

	return server
}
