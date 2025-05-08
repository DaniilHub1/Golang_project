package handlers

import (
	"mini_site/database"
	"mini_site/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateComment(c *gin.Context) {
	username, err := c.Cookie("username")
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login?error=Сначала войдите в аккаунт")
		return
	}

	content := c.PostForm("content")
	postIDStr := c.PostForm("post_id")

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Неверный ID поста")
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.String(http.StatusBadRequest, "Пользователь не найден")
		return
	}

	comment := models.Comment{
		Content: content,
		UserID:  user.ID,
		PostID:  uint(postID),
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.String(http.StatusInternalServerError, "Не удалось сохранить комментарий")
		return
	}

	c.Redirect(http.StatusFound, "/posts_page")
}
