package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"html/template"
	"log"
	"mini_site/models"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"gorm.io/gorm"
	"mini_site/database"
)

var ctx = context.Background()
var rdb *redis.Client
var posts []models.Post
var DB *gorm.DB


func SomeHandler(c *gin.Context) {

}

func RenderPostsPage(c *gin.Context) {
	posts := []models.Post{
		{
			ID:      1,
			Title:   "Post 1",
			Content: "Content of post 1",
			Author:  "Author 1",
			Date:    time.Now(),
		},
		{
			ID:      2,
			Title:   "Post 2",
			Content: "Content of post 2",
			Author:  "Author 2",
			Date:    time.Now(),
		},
	}

	log.Println("Пост написан:", posts)

	c.HTML(http.StatusOK, "posts_page.html", gin.H{
		"posts": posts,
	})
}

func CreatePost(c *gin.Context) {
	var newPost models.Post
	if err := c.ShouldBindJSON(&newPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newPost.ID = len(posts) + 1
	posts = append(posts, newPost)
	c.JSON(http.StatusCreated, newPost)
}

func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var updatedPost models.Post
	if err := c.ShouldBindJSON(&updatedPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, post := range posts {
		if strconv.Itoa(post.ID) == id {
			posts[i].Title = updatedPost.Title
			posts[i].Content = updatedPost.Content
			c.JSON(http.StatusOK, posts[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
}

func DeletePost(c *gin.Context) {
	id := c.Param("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID поста"})
		return
	}

	if err := DeletePostFromRedis(postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении поста из Redis"})
		return
	}

	for i, post := range posts {
		if post.ID == postID {
			posts = append(posts[:i], posts[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Пост удален"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Пост не найден"})
}

func HomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func AboutPage(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", nil)
}

func CreatePostFromForm(c *gin.Context) {
	title := c.DefaultPostForm("title", "")
	content := c.DefaultPostForm("content", "")
	username := c.DefaultPostForm("username", "")

	newPost := models.Post{
		ID:       len(posts) + 1,
		Title:    title,
		Content:  content,
		Username: username,
	}

	posts = append(posts, newPost)

	c.Redirect(http.StatusSeeOther, "/posts_page")
}

func RenderPage2(c *gin.Context) {
	tmpl, err := template.ParseFiles("templates/account.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка загрузки шаблона")
		return
	}

	tmpl.Execute(c.Writer, nil)
}
func RenderPage3(c *gin.Context) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка загрузки шаблона")
		return
	}
	tmpl.Execute(c.Writer, nil)
}

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "admin" && password == "password" {
		c.SetCookie("token", "user-auth-token", 3600, "/", "localhost", false, true)
		c.Redirect(http.StatusFound, "/account")
	} else {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Неверные данные"})
	}
}

func SettingsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "settings.html", gin.H{
		"message": "",
	})
}

func SaveSettings(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")

	file, err := c.FormFile("avatar")
	var avatarPath string
	if err == nil {
		avatarPath = filepath.Join("static", "uploads", file.Filename)

		if _, err := os.Stat("static/uploads"); os.IsNotExist(err) {
			os.MkdirAll("static/uploads", os.ModePerm)
		}

		if err := c.SaveUploadedFile(file, avatarPath); err != nil {
			c.String(http.StatusInternalServerError, "Ошибка при загрузке изображения")
			return
		}
	}

	c.HTML(http.StatusOK, "settings.html", gin.H{
		"username": username,
		"email":    email,
		"avatar":   "/" + avatarPath,
		"message":  "Настройки сохранены!",
	})
}

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", 
		Password: "",               
		DB:       0,                
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}
}

func SaveProfileSettings(c *gin.Context) {
	nickname := c.PostForm("nickname")

	file, err := c.FormFile("profile_photo")
	var photoPath string
	if err == nil {
		photoPath = "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, photoPath); err != nil {
			c.String(http.StatusInternalServerError, "Не удалось сохранить файл")
			return
		}
	} else {
		photoPath = "uploads/default.png"
	}

	if err := rdb.Set(ctx, "profile:nickname", nickname, 0).Err(); err != nil {
		c.String(http.StatusInternalServerError, "Не удалось сохранить никнейм")
		return
	}
	if err := rdb.Set(ctx, "profile:photo", photoPath, 0).Err(); err != nil {
		c.String(http.StatusInternalServerError, "Не удалось сохранить фотографию")
		return
	}

	c.Redirect(http.StatusSeeOther, "/settings")
}

func LoadProfileSettings() (string, string, error) {
	nickname, err := rdb.Get(ctx, "profile:nickname").Result()
	if err != nil && err != redis.Nil {
		return "", "", err
	}
	if nickname == "" {
		nickname = "Guest" 
	}

	photoPath, err := rdb.Get(ctx, "profile:photo").Result()
	if err != nil && err != redis.Nil {
		return "", "", err
	}
	if photoPath == "" {
		photoPath = "uploads/default.png"
	}

	return nickname, photoPath, nil
}

func GetPosts(c *gin.Context) {
    var posts []models.Post
    if err := database.DB.Find(&posts).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, posts)
}
