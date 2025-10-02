package user

import (
	"net/http"

	modelDTO "github.com/HPNV/growlink-backend/model/dto"
	"github.com/HPNV/growlink-backend/service"
	"github.com/gin-gonic/gin"
)

type IUser interface {
	Login(*gin.Context)
	Register(*gin.Context)
	GetAll(*gin.Context)
	GetDetail(*gin.Context)
	GetStudentList(*gin.Context)
}

type User struct {
	service service.IRegistry
}

func NewUser(service service.IRegistry) *User {
	return &User{
		service: service,
	}
}

func (u *User) Login(c *gin.Context) {
	var req modelDTO.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := u.service.GetUser().Login(c, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (u *User) Register(c *gin.Context) {
	var req modelDTO.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := u.service.GetUser().Register(c, req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (u *User) GetAll(c *gin.Context) {
	users, err := u.service.GetUser().GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (u *User) GetDetail(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID parameter is required"})
		return
	}

	user, err := u.service.GetUser().GetDetail(c, uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (u *User) GetStudentList(c *gin.Context) {
	var req modelDTO.StudentListRequest

	// Parse query parameters
	if name := c.Query("name"); name != "" {
		req.Name = &name
	}

	if university := c.Query("university"); university != "" {
		req.University = &university
	}

	if skill := c.Query("skill"); skill != "" {
		req.Skill = &skill
	}

	// Parse page and limit with defaults
	page := 1
	if pageStr := c.DefaultQuery("page", "1"); pageStr != "" {
		if pageInt, err := parseIntFromString(pageStr); err == nil {
			page = pageInt
		}
	}
	req.Page = page

	limit := 10
	if limitStr := c.DefaultQuery("limit", "10"); limitStr != "" {
		if limitInt, err := parseIntFromString(limitStr); err == nil {
			limit = limitInt
		}
	}
	req.Limit = limit

	response, err := u.service.GetUser().GetStudentList(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Helper function to parse string to int
func parseIntFromString(s string) (int, error) {
	result := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, gin.Error{Err: nil}
		}
		result = result*10 + int(r-'0')
	}
	return result, nil
}
