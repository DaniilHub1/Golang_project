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

func GetCurrentUser(c *gin.Context) {
	username, err := c.Cookie("username")
	if err != nil || username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": user.ID, "username": user.Username})
}
func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	username, err := c.Cookie("username")
	if err != nil || username == "" {
		return 0, false
	}

	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return 0, false
	}

	return user.ID, true
}


func DeleteOwnMessage(c *gin.Context) {
    id := c.Param("id")
    userID, _ := GetUserIDFromContext(c)
    
    var msg models.Message
    if err := DB.First(&msg, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Сообщение не найдено"})
        return
    }
    
    if msg.SenderID != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "Можно удалять только свои сообщения"})
        return
    }
    
    DB.Delete(&msg)
    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}


func GetMessages(c *gin.Context) {
    userIDStr := c.Query("user_id")
    friendIDStr := c.Query("friend_id")

    userID, _ := strconv.Atoi(userIDStr)
    friendID, _ := strconv.Atoi(friendIDStr)

    var messages []models.Message
    database.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
        userID, friendID, friendID, userID).
        Order("created_at ASC"). // Важно: сортируем по времени создания
        Find(&messages)

    c.JSON(http.StatusOK, messages)
}

func GetUsers(c *gin.Context) {
    search := c.Query("search")
    var users []models.User
    query := DB
    if search != "" {
        query = query.Where("username LIKE ?", "%"+search+"%")
    }
    query.Find(&users)
    c.JSON(http.StatusOK, users)
}

