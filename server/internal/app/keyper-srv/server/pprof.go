package server

import (
	"net/http"
	"net/http/pprof"

	"github.com/RomanAgaltsev/keyper/server/internal/config"
)

func NewPprofServer(cfg *config.PprofConfig) *http.Server {

	handler := http.NewServeMux()
	handler.HandleFunc("/debug/pprof/", pprof.Index)
	handler.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	handler.HandleFunc("/debug/pprof/profile", pprof.Profile)
	handler.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	handler.HandleFunc("/debug/pprof/trace", pprof.Trace)

	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: handler,
	}

	return srv
}
