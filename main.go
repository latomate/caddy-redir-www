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

func (rd *RedirWww) Provision(ctx caddy.Context) error {
	rd.logger = ctx.Logger() // g.logger is a *zap.Logger

	tlsAppIface, err := ctx.App("tls")
	if err != nil {
		return fmt.Errorf("getting tls app: %v", err)
	}
	rd.tlsApp = tlsAppIface.(*caddytls.TLS)

	return nil
}

func (rd RedirWww) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	//response, err := net.LookupTXT("_redirwww." + r.Host)

	//if err != nil || len(response) == 0 {
	//	rd.logger.Info("error", zap.String("host", "_redirwww."+r.Host), zap.Error(err))
	//	return next.ServeHTTP(w, r)
	//}

	rd.logger.Info("redir_www", zap.String("ask", rd.tlsApp.Automation.OnDemand.Ask))

	//rd.logger.Info("redir_www", zap.String("host", r.Host), zap.String("response", response[0]))

	//w.Header().Set("Location", response[0])
	//w.WriteHeader(301)

	return next.ServeHTTP(w, r)
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (rd *RedirWww) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	return nil
}

// parseCaddyfile unmarshals tokens from h into a new RedirWww.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var rd RedirWww
	err := rd.UnmarshalCaddyfile(h.Dispenser)
	return rd, err
}

// Interface guard
var (
	_ caddy.Provisioner = (*RedirWww)(nil)
	// _ caddy.Validator             = (*RedirWww)(nil)
	_ caddyhttp.MiddlewareHandler = (*RedirWww)(nil)
	_ caddyfile.Unmarshaler       = (*RedirWww)(nil)
)
