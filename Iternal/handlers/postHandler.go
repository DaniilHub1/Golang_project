package handlers

import (
	"mini_site/database"
	"mini_site/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var posts []models.Post
var DB *gorm.DB

func CreatePost(c *gin.Context) {
	username, err := c.Cookie("username")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Вы не авторизованы"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, "username = ?", username).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	var newPost models.Post
	if err := c.ShouldBindJSON(&newPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newPost.UserID = user.ID

	if err := database.DB.Create(&newPost).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать пост"})
		return
	}

	c.JSON(http.StatusCreated, newPost)
}

func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	var updatedData models.Post
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var post models.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пост не найден"})
		return
	}

	post.Title = updatedData.Title
	post.Content = updatedData.Content

	if err := database.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func DeletePost(c *gin.Context) {
	username, err := c.Cookie("username")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Вы не авторизованы"})
		return
	}

	role, _ := c.Cookie("user_role")

	id := c.Param("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID поста"})
		return
	}

	var post models.Post
	if err := database.DB.Preload("User").First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пост не найден"})
		return
	}

	// Проверка, является ли пользователь владельцем поста или администратором
	if post.User.Username != username && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "У вас нет прав на удаление этого поста"})
		return
	}

	if err := database.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении поста"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пост удалён"})
}

func HomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func AboutPage(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", nil)
}

func CreatePostFromForm(c *gin.Context) {
	username, err := c.Cookie("username")
	if err != nil {
		c.Redirect(http.StatusFound, "/login?error=Unauthorized")
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.Redirect(http.StatusFound, "/login?error=Unauthorized")
		return
	}

	content := c.PostForm("content")
	if content == "" {
		c.Redirect(http.StatusFound, "/account?error=Content+required")
		return
	}

	post := models.Post{
		Content: content,
		UserID:  user.ID,
	}
	database.DB.Create(&post)

	c.Redirect(http.StatusFound, "/posts_page")
}

func GetPosts(c *gin.Context) {
	var posts []models.Post
	if err := database.DB.Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func RenderPostsPage(c *gin.Context) {
	username, err := c.Cookie("username")
	if err != nil {
		c.Redirect(http.StatusFound, "/login?error=Unauthorized")
		return
	}

	role, _ := c.Cookie("user_role")

	var posts []models.Post
	database.DB.Order("created_at desc").Preload("User").Find(&posts)

	// Структура для постов с дополнительной информацией о пользователе
	type PostWithUser struct {
		Content   string
		CreatedAt string
		Username  string
		IsAdmin   bool // Флаг для проверки, является ли пользователь администратором
	}

	var data []PostWithUser
	for _, post := range posts {
		data = append(data, PostWithUser{
			Content:   post.Content,
			CreatedAt: post.CreatedAt.Format("02.01.2006 15:04"),
			Username:  post.User.Username,
			IsAdmin:   (role == "admin"), // Проверка на роль администратора
		})
	}

	c.HTML(http.StatusOK, "posts.html", gin.H{
		"Posts": data,
		"User":  username,
	})
}

type PostWithUser struct {
	Content   string
	CreatedAt string
	Username  string
}

func RenderPostsFeed(c *gin.Context) {
	username, err := c.Cookie("username")
	if err != nil || username == "" {
		c.Redirect(http.StatusSeeOther, "/login?error=Авторизуйтесь для просмотра постов")
		return
	}

	var posts []models.Post
	if err := database.DB.Order("created_at desc").Find(&posts).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": "Не удалось получить посты"})
		return
	}

	var data []PostWithUser
	for _, post := range posts {
		var user models.User
		database.DB.First(&user, post.UserID)
		data = append(data, PostWithUser{
			Content:   post.Content,
			CreatedAt: post.CreatedAt.Format("02.01.2006 15:04"),
			Username:  user.Username,
		})
	}

	c.HTML(http.StatusOK, "posts.html", gin.H{
		"Posts": data,
		"User":  username,
	})
}
