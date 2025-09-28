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

	// Serve static files
	r.engine.Static("/static", "./static")

	v := r.engine.Group("/v1")

	r.userRoute(v)
	r.businessRoute(v)
	r.studentRoute(v)
	r.skillRoute(v)
	r.projectRoute(v)
	r.fileRoute(v)

	r.engine.Run(":" + r.cfg.Port)
}

func (r *Route) userRoute(g *gin.RouterGroup) {
	user := r.delivery.GetUser()
	u := g.Group("/user")
	u.POST("/login", user.Login)
	u.POST("/register", user.Register)
	u.GET("", user.GetAll)
}

func (r *Route) businessRoute(g *gin.RouterGroup) {
	business := r.delivery.GetBusiness()
	b := g.Group("/business")

	b.POST("", business.Create)
	b.GET("", business.GetAll)
	b.GET("/:uuid", business.GetByUUID)
	b.GET("/user/:userUuid", business.GetByUserUUID)
	b.PUT("/:uuid", business.Update)
	b.DELETE("/:uuid", business.Delete)
}

func (r *Route) studentRoute(g *gin.RouterGroup) {
	student := r.delivery.GetStudent()
	s := g.Group("/student")

	s.POST("", student.Create)
	s.GET("", student.GetAll)
	s.GET("/:uuid", student.GetByUUID)
	s.GET("/user/:userUuid", student.GetByUserUUID)
	s.PUT("/:uuid", student.Update)
	s.DELETE("/:uuid", student.Delete)

	// Student skills management
	s.POST("/:uuid/skills/:skillUuid", student.AddSkill)
	s.DELETE("/:uuid/skills/:skillUuid", student.RemoveSkill)
	s.GET("/:uuid/skills", student.GetSkills)
}

func (r *Route) skillRoute(g *gin.RouterGroup) {
	skill := r.delivery.GetSkill()
	sk := g.Group("/skill")

	sk.POST("", skill.Create)
	sk.GET("", skill.GetAll)
	sk.GET("/:uuid", skill.GetByUUID)
	sk.PUT("/:uuid", skill.Update)
	sk.DELETE("/:uuid", skill.Delete)
}

func (r *Route) projectRoute(g *gin.RouterGroup) {
	project := r.delivery.GetProject()
	p := g.Group("/project")

	p.POST("/business/:businessUuid", project.Create)
	p.GET("", project.GetAll)
	p.GET("/:uuid", project.GetByUUID)
	p.GET("/business/:businessUuid", project.GetByBusinessUUID)
	p.PUT("/:uuid", project.Update)
	p.DELETE("/:uuid", project.Delete)

	// Project skills management
	p.POST("/:uuid/skills/:skillUuid", project.AddSkill)
	p.DELETE("/:uuid/skills/:skillUuid", project.RemoveSkill)
	p.GET("/:uuid/skills", project.GetSkills)

	// Project students management
	p.POST("/:uuid/students/:studentUuid", project.AddStudent)
	p.DELETE("/:uuid/students/:studentUuid", project.RemoveStudent)
	p.GET("/:uuid/students", project.GetStudents)
}

func (r *Route) fileRoute(g *gin.RouterGroup) {
	file := r.delivery.GetFile()
	f := g.Group("/file")

	f.POST("/upload", file.UploadImage)
	f.GET("/:uuid", file.GetByUUID)
	f.DELETE("/:uuid", file.Delete)
	f.GET("/user/:uploadedBy", file.GetByUploadedBy)
}
