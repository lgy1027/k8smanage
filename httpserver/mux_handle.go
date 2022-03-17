package httpserver

import (
	log "github.com/cihub/seelog"
	"github.com/go-chi/chi"
	_ "github.com/lgy1027/kubemanage/docs"
	"github.com/lgy1027/kubemanage/pkg/cluster"
	"github.com/lgy1027/kubemanage/pkg/deploy"
	"github.com/lgy1027/kubemanage/pkg/webshell"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func Mux(open bool) *http.ServeMux {
	log.Info("action", "start init service")
	r := http.NewServeMux()
	var (
		clusterService     = cluster.NewService()
		clusterEndpoints   = cluster.NewEndpoints(clusterService)
		clusterHttpHandler = cluster.NewHTTPHandler(clusterEndpoints)

		deployService     = deploy.NewService()
		deployEndpoints   = deploy.NewEndpoints(deployService)
		deployHttpHandler = deploy.NewHTTPHandler(deployEndpoints)
	)

	r.Handle("/cluster/", http.StripPrefix("/cluster", clusterHttpHandler))
	r.Handle("/resource/", http.StripPrefix("/resource", deployHttpHandler))
	r.HandleFunc("/ws/{namespace}/{pod}/{container}/log", webshell.ServeWsLogs)
	r.HandleFunc("/ws/{namespace}/{pod}/{container}/shell", webshell.ServeWsTerminal)
	r.HandleFunc("/v1/pod/log/", webshell.LogHandle)
	log.Info("action", "success init service")
	if open {
		go func() {
			r := chi.NewRouter()
			r.Get("/swagger/*", httpSwagger.Handler(
				httpSwagger.URL("http://localhost:7475/swagger/doc.json")))
			_ = http.ListenAndServe(":7475", r)
		}()
	}
	return r
}
