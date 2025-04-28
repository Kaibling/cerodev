package api

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/kaibling/cerodev/bootstrap"
	"github.com/kaibling/cerodev/bootstrap/appctx"
)

func proxyHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxyPath := strings.TrimPrefix(r.URL.Path, "/proxy")
		containerID := strings.Split(proxyPath, "/")[1]

		_, l, cfg, err := appctx.GetBaseData(r.Context())
		if err != nil {
			l.Warn("could not read context: %s", err.Error())
			http.Error(w, "Not found", http.StatusNotFound)

			return
		}

		cs, err := bootstrap.NewContainerService(r.Context())
		if err != nil {
			l.Warn("could not build containerservice: %s", err.Error())
			http.Error(w, "Not found", http.StatusNotFound)

			return
		}

		c, err := cs.GetByID(containerID)
		if err != nil {
			l.Warn("could not read container: %s", err.Error())
			http.Error(w, "Not found", http.StatusNotFound)

			return
		}

		targetURL := cfg.PublicURL + ":" + c.UIPort

		target, err := url.Parse(targetURL)
		if err != nil {
			l.Warn("could not build proxy target: %s", err.Error())
			http.Error(w, "Not found", http.StatusNotFound)

			return
		}

		proxy := httputil.NewSingleHostReverseProxy(target)

		// Fix redirects from backend that use Location header
		proxy.ModifyResponse = func(resp *http.Response) error {
			location := resp.Header.Get("Location")
			if location != "" {
				if strings.HasPrefix(location, "/") {
					resp.Header.Set("Location", "/proxy/"+containerID+location)
				} else if strings.HasPrefix(location, "./") {
					newLoc := "/proxy/" + containerID + "/" + strings.TrimPrefix(location, "./")
					resp.Header.Set("Location", newLoc)
				}
			}

			return nil
		}

		newPath := strings.TrimPrefix(proxyPath, "/"+containerID)
		r.URL.Path = newPath
		proxy.ServeHTTP(w, r)
	})
}
