package controllers

import (
	"crud-ukom/config"
	"crud-ukom/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get all packets
func GetPackets(c *gin.Context) {
	var packets []models.Packet
	config.DB.Find(&packets)
	c.JSON(http.StatusOK, packets)
}

// Get packet by ID
func GetPacketByID(c *gin.Context) {
	var packet models.Packet
	id := c.Param("id")
	if err := config.DB.First(&packet, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Packet not found!"})
		return
	}
	c.JSON(http.StatusOK, packet)
}

// Create new packet
func CreatePacket(c *gin.Context) {
	var input models.Packet
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&input)
	c.JSON(http.StatusOK, input)
}

// Update packet by ID
func UpdatePacket(c *gin.Context) {
	var packet models.Packet
	id := c.Param("id")
	if err := config.DB.First(&packet, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Packet not found!"})
		return
	}

	var input models.Packet
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&packet).Updates(input)
	c.JSON(http.StatusOK, packet)
}

// Delete packet by ID
func DeletePacket(c *gin.Context) {
	var packet models.Packet
	id := c.Param("id")
	if err := config.DB.First(&packet, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Packet not found!"})
		return
	}
	config.DB.Delete(&packet)
	c.JSON(http.StatusOK, gin.H{"message": "Packet deleted successfully!"})
}
