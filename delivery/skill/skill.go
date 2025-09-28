package skill

import (
	"net/http"

	"github.com/HPNV/growlink-backend/model/dto"
	"github.com/HPNV/growlink-backend/service"
	"github.com/gin-gonic/gin"
)

type ISkill interface {
	Create(c *gin.Context)
	GetByUUID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetAll(c *gin.Context)
}

type Skill struct {
	service service.IRegistry
}

func NewSkill(service service.IRegistry) ISkill {
	return &Skill{
		service: service,
	}
}

func (s *Skill) Create(c *gin.Context) {
	var req dto.SkillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	skill, err := s.service.GetSkill().Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, skill)
}

func (s *Skill) GetByUUID(c *gin.Context) {
	uuid := c.Param("uuid")

	skill, err := s.service.GetSkill().GetByUUID(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}

	c.JSON(http.StatusOK, skill)
}

func (s *Skill) Update(c *gin.Context) {
	uuid := c.Param("uuid")
	var req dto.SkillRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	skill, err := s.service.GetSkill().Update(uuid, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, skill)
}

func (s *Skill) Delete(c *gin.Context) {
	uuid := c.Param("uuid")

	err := s.service.GetSkill().Delete(uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Skill deleted successfully"})
}

func (s *Skill) GetAll(c *gin.Context) {
	skills, err := s.service.GetSkill().GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, skills)
}
