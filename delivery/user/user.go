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
