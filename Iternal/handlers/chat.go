package handlers

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"mini_site/database"
	"mini_site/models"
)

func SendMessage(c *gin.Context) {
	senderIDStr := c.PostForm("sender_id")
	receiverIDStr := c.PostForm("receiver_id")
	content := c.PostForm("content")

	senderID, _ := strconv.Atoi(senderIDStr)
	receiverID, _ := strconv.Atoi(receiverIDStr)

	message := models.Message{
		SenderID:   uint(senderID),
		ReceiverID: uint(receiverID),
		Content:    content,
	}

	if err := database.DB.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось отправить сообщение"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "сообщение отправлено"})
}

func GetMessages(c *gin.Context) {
	userIDStr := c.Query("user_id")
	friendIDStr := c.Query("friend_id")

	userID, _ := strconv.Atoi(userIDStr)
	friendID, _ := strconv.Atoi(friendIDStr)

	var messages []models.Message
	database.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		userID, friendID, friendID, userID).Find(&messages)

	c.JSON(http.StatusOK, messages)
}

func GetUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}