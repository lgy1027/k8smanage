package httpserver

import (
	log "github.com/cihub/seelog"
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	_ "relaper.com/kubemanage/docs"
	"relaper.com/kubemanage/pkg/cluster"
	"relaper.com/kubemanage/pkg/deploy"
)

func Mux(open bool, address string) *http.ServeMux {
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
