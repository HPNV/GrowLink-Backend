package routing

import (
	"github.com/HPNV/growlink-backend/config"
	"github.com/HPNV/growlink-backend/delivery"
	"github.com/gin-gonic/gin"
)

type Route struct {
	cfg      config.ServerConfig
	engine   *gin.Engine
	delivery delivery.IDelivery
}

func NewRoute(cfg config.ServerConfig, delivery delivery.IDelivery) *Route {
	gin.SetMode(cfg.Mode)

	return &Route{
		cfg:      cfg,
		engine:   gin.Default(),
		delivery: delivery,
	}
}

func (r *Route) SetupRoutes() {
	r.engine = gin.Default()
	v := r.engine.Group("/v1")

	r.dummyRoute(v)
	r.developerRoute(v)

	r.engine.Run(":" + r.cfg.Port)
}

func (r *Route) dummyRoute(g *gin.RouterGroup) {
	g.GET("/dummy")
}

func (r *Route) developerRoute(g *gin.RouterGroup) {
	group := g.Group("/developer")
	group.POST("", r.delivery.GetDeveloper().CreateDeveloper)
}
