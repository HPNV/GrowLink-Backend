package business

import (
	"net/http"

	"github.com/HPNV/growlink-backend/model/dto"
	"github.com/HPNV/growlink-backend/service"
	"github.com/gin-gonic/gin"
)

type IBusiness interface {
	Create(c *gin.Context)
	GetByUUID(c *gin.Context)
	GetByUserUUID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetAll(c *gin.Context)
}

type Business struct {
	service service.IRegistry
}

func NewBusiness(service service.IRegistry) IBusiness {
	return &Business{
		service: service,
	}
}

func (b *Business) Create(c *gin.Context) {
	var req dto.BusinessRequest
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

	business, err := b.service.GetBusiness().Create(userUUID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, business)
}

func (b *Business) GetByUUID(c *gin.Context) {
	uuid := c.Param("uuid")

	business, err := b.service.GetBusiness().GetByUUID(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Business not found"})
		return
	}

	c.JSON(http.StatusOK, business)
}

func (b *Business) GetByUserUUID(c *gin.Context) {
	userUUID := c.Param("userUuid")

	business, err := b.service.GetBusiness().GetByUserUUID(userUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Business not found"})
		return
	}

	c.JSON(http.StatusOK, business)
}

func (b *Business) Update(c *gin.Context) {
	uuid := c.Param("uuid")
	var req dto.BusinessRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	business, err := b.service.GetBusiness().Update(uuid, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, business)
}

func (b *Business) Delete(c *gin.Context) {
	uuid := c.Param("uuid")

	err := b.service.GetBusiness().Delete(uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Business deleted successfully"})
}

func (b *Business) GetAll(c *gin.Context) {
	businesses, err := b.service.GetBusiness().GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, businesses)
}
