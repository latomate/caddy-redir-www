package redir_www

import (
	"fmt"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/caddyserver/caddy/v2/modules/caddytls"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(RedirWww{})
	httpcaddyfile.RegisterHandlerDirective("redir_www", parseCaddyfile)
}

// RedirWww is a redirirection for non-www websites.
type RedirWww struct {
	logger *zap.Logger
	tlsApp *caddytls.TLS
}

// CaddyModule returns the Caddy module information.
func (RedirWww) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.redir_www",
		New: func() caddy.Module { return new(RedirWww) },
	}
}

func (rw *RedirWww) Provision(ctx caddy.Context) error {
	rw.logger = ctx.Logger() // g.logger is a *zap.Logger

	tlsAppIface, err := ctx.App("tls")
	if err != nil {
		return fmt.Errorf("getting tls app: %v", err)
	}
	rw.tlsApp = tlsAppIface.(*caddytls.TLS)

	return nil
}

func (rw RedirWww) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	//response, err := net.LookupTXT("_redirwww." + r.Host)

	//if err != nil || len(response) == 0 {
	//    rw.logger.Info("error", zap.String("host", "_redirwww."+r.Host), zap.Error(err))
	//    return next.ServeHTTP(w, r)
	//}

	if rw.tlsApp.Automation != nil {
		if rw.tlsApp.Automation.OnDemand != nil {
			rw.logger.Info("redir_www", zap.String("ask", rw.tlsApp.Automation.OnDemand.Ask))
		}
	}

	//rw.logger.Info("redir_www", zap.String("host", r.Host), zap.String("response", response[0]))

	//w.Header().Set("Location", response[0])
	//w.WriteHeader(301)

	return next.ServeHTTP(w, r)
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (rw *RedirWww) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	return nil
}

// parseCaddyfile unmarshals tokens from h into a new RedirWww.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var rw RedirWww
	err := rw.UnmarshalCaddyfile(h.Dispenser)
	return rw, err
}

// Interface guard
var (
	_ caddy.Provisioner = (*RedirWww)(nil)
	// _ caddy.Validator             = (*RedirWww)(nil)
	_ caddyhttp.MiddlewareHandler = (*RedirWww)(nil)
	_ caddyfile.Unmarshaler       = (*RedirWww)(nil)
)
