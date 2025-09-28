package student

import (
	"net/http"

	"github.com/HPNV/growlink-backend/model/dto"
	"github.com/HPNV/growlink-backend/service"
	"github.com/gin-gonic/gin"
)

type IStudent interface {
	Create(c *gin.Context)
	GetByUUID(c *gin.Context)
	GetByUserUUID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetAll(c *gin.Context)
	AddSkill(c *gin.Context)
	RemoveSkill(c *gin.Context)
	GetSkills(c *gin.Context)
}

type Student struct {
	service service.IRegistry
}

func NewStudent(service service.IRegistry) IStudent {
	return &Student{
		service: service,
	}
}

func (s *Student) Create(c *gin.Context) {
	var req dto.StudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user UUID from context (would be set by auth middleware)
	userUUID, exists := c.Get("user_uuid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	student, err := s.service.GetStudent().Create(userUUID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, student)
}

func (s *Student) GetByUUID(c *gin.Context) {
	uuid := c.Param("uuid")

	student, err := s.service.GetStudent().GetByUUID(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, student)
}

func (s *Student) GetByUserUUID(c *gin.Context) {
	userUUID := c.Param("userUuid")

	student, err := s.service.GetStudent().GetByUserUUID(userUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, student)
}

func (s *Student) Update(c *gin.Context) {
	uuid := c.Param("uuid")
	var req dto.StudentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student, err := s.service.GetStudent().Update(uuid, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, student)
}

func (s *Student) Delete(c *gin.Context) {
	uuid := c.Param("uuid")

	err := s.service.GetStudent().Delete(uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
}

func (s *Student) GetAll(c *gin.Context) {
	students, err := s.service.GetStudent().GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, students)
}

func (s *Student) AddSkill(c *gin.Context) {
	studentUUID := c.Param("uuid")
	skillUUID := c.Param("skillUuid")

	err := s.service.GetStudent().AddSkill(studentUUID, skillUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Skill added successfully"})
}

func (s *Student) RemoveSkill(c *gin.Context) {
	studentUUID := c.Param("uuid")
	skillUUID := c.Param("skillUuid")

	err := s.service.GetStudent().RemoveSkill(studentUUID, skillUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Skill removed successfully"})
}

func (s *Student) GetSkills(c *gin.Context) {
	studentUUID := c.Param("uuid")

	skills, err := s.service.GetStudent().GetSkills(studentUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, skills)
}
