package developer

import (
	modelDB "github.com/HPNV/growlink-backend/model/db"
	"github.com/HPNV/growlink-backend/service"
	"github.com/gin-gonic/gin"
)

type IDeveloper interface {
	CreateDeveloper(*gin.Context)
}

type Developer struct {
	service service.IRegistry
}

func NewDeveloperService(service service.IRegistry) *Developer {
	return &Developer{
		service: service,
	}
}

func (s *Developer) CreateDeveloper(ctx *gin.Context) {
	var developer modelDB.Developer
	if err := ctx.ShouldBindJSON(&developer); err != nil {
		if err := ctx.ShouldBindJSON(&developer); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
	}

	err := s.service.GetDeveloper().CreateDeveloper(&developer)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"message": "Developer created successfully"})
}
