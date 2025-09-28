package project

import (
	"net/http"

	"github.com/HPNV/growlink-backend/model/dto"
	"github.com/HPNV/growlink-backend/service"
	"github.com/gin-gonic/gin"
)

type IProject interface {
	Create(c *gin.Context)
	GetByUUID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetAll(c *gin.Context)
	GetByBusinessUUID(c *gin.Context)
	AddSkill(c *gin.Context)
	RemoveSkill(c *gin.Context)
	GetSkills(c *gin.Context)
	AddStudent(c *gin.Context)
	RemoveStudent(c *gin.Context)
	GetStudents(c *gin.Context)
}

type Project struct {
	service service.IRegistry
}

func NewProject(service service.IRegistry) IProject {
	return &Project{
		service: service,
	}
}

func (p *Project) Create(c *gin.Context) {
	var req dto.ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get business UUID from context or params
	businessUUID := c.Param("businessUuid")
	if businessUUID == "" {
		// Could also get from authenticated user context
		c.JSON(http.StatusBadRequest, gin.H{"error": "Business UUID required"})
		return
	}

	project, err := p.service.GetProject().Create(businessUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, project)
}

func (p *Project) GetByUUID(c *gin.Context) {
	uuid := c.Param("uuid")

	project, err := p.service.GetProject().GetByUUID(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

func (p *Project) Update(c *gin.Context) {
	uuid := c.Param("uuid")
	var req dto.ProjectUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := p.service.GetProject().Update(uuid, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, project)
}

func (p *Project) Delete(c *gin.Context) {
	uuid := c.Param("uuid")

	err := p.service.GetProject().Delete(uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

func (p *Project) GetAll(c *gin.Context) {
	projects, err := p.service.GetProject().GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}

func (p *Project) GetByBusinessUUID(c *gin.Context) {
	businessUUID := c.Param("businessUuid")

	projects, err := p.service.GetProject().GetByBusinessUUID(businessUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}

func (p *Project) AddSkill(c *gin.Context) {
	projectUUID := c.Param("uuid")
	skillUUID := c.Param("skillUuid")

	err := p.service.GetProject().AddSkill(projectUUID, skillUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Skill added to project successfully"})
}

func (p *Project) RemoveSkill(c *gin.Context) {
	projectUUID := c.Param("uuid")
	skillUUID := c.Param("skillUuid")

	err := p.service.GetProject().RemoveSkill(projectUUID, skillUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Skill removed from project successfully"})
}

func (p *Project) GetSkills(c *gin.Context) {
	projectUUID := c.Param("uuid")

	skills, err := p.service.GetProject().GetSkills(projectUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, skills)
}

func (p *Project) AddStudent(c *gin.Context) {
	projectUUID := c.Param("uuid")
	studentUUID := c.Param("studentUuid")

	err := p.service.GetProject().AddStudent(projectUUID, studentUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student added to project successfully"})
}

func (p *Project) RemoveStudent(c *gin.Context) {
	projectUUID := c.Param("uuid")
	studentUUID := c.Param("studentUuid")

	err := p.service.GetProject().RemoveStudent(projectUUID, studentUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student removed from project successfully"})
}

func (p *Project) GetStudents(c *gin.Context) {
	projectUUID := c.Param("uuid")

	students, err := p.service.GetProject().GetStudents(projectUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, students)
}
