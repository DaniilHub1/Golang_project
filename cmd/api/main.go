package main

import (
	"html/template"
	"mini_site/Iternal/admin"
	"mini_site/Iternal/handlers"
	"mini_site/database"
	"mini_site/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		panic("failed to connect to the database")
	}

	// Автомиграция моделей
	db.AutoMigrate(
		&models.User{},
		&models.Message{},
		&models.Post{},
		&models.Comment{},
	)
	handlers.DB = db

	router := gin.Default()

	// Настройка функций для шаблонов
	router.SetFuncMap(template.FuncMap{
		"currentYear": func() int {
			return time.Now().Year()
		},
	})

	// Статические файлы
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/**/*.html")

	// Основные маршруты
	router.GET("/", handlers.HomePage)
	router.GET("/about", handlers.AboutPage)
	router.GET("/account", handlers.RenderPage2)
	router.GET("/index.html", handlers.RenderPage3)
	router.GET("/logout", handlers.LogoutHandler)

	// Посты
	router.GET("/posts", handlers.GetPosts)
	router.POST("/posts_page", handlers.CreatePostFromForm)
	router.PUT("/posts/:id", handlers.UpdatePost)
	router.DELETE("/posts/:id", handlers.DeletePost)
	router.GET("/posts_page", handlers.RenderPostsPage)

	// Аутентификация
	router.GET("/login", func(c *gin.Context) {
		errorMsg := c.Query("error")
		successMsg := c.Query("success")

		c.HTML(http.StatusOK, "login.html", gin.H{
			"Error":   errorMsg,
			"Success": successMsg,
		})
	})
	router.POST("/login", handlers.LoginHandler)

	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})
	router.POST("/register", handlers.RegisterHandler)

	// Настройки профиля
	router.GET("/settings", handlers.SettingsPage)
	router.POST("/settings", handlers.SaveSettings)
	router.POST("/profile/save", handlers.SaveSettings)

	// Комментарии
	router.POST("/comment", handlers.CreateComment)

	// Лента постов
	router.GET("/feed", handlers.RenderPostsFeed)

 

	// Чат и сообщения (обновленные роуты)
	router.GET("/chat", handlers.AuthRequired(), func(c *gin.Context) {
		username, _ := c.Cookie("username")
		c.HTML(http.StatusOK, "chat.html", gin.H{
			"Username": username,
			"Now":      time.Now(),
		})
	})
	router.POST("/send", handlers.AuthRequired(), handlers.SendMessage)
	router.GET("/messages", handlers.AuthRequired(), handlers.GetMessages)
	router.GET("/users", handlers.AuthRequired(), handlers.GetUsers)
	router.GET("/get_current_user", handlers.GetCurrentUser)
	router.POST("/messages/:id/delete", handlers.AuthRequired(), handlers.DeleteOwnMessage)

	// Админ-панель
	adminGroup := router.Group("/admin")
	adminGroup.Use(handlers.AdminRequired())
	{
		adminGroup.GET("/dashboard", admin.AdminDashboard)
		adminGroup.DELETE("/users/:id", admin.DeleteUserByAdmin)
		adminGroup.GET("/users", admin.AdminUsersPage)
		adminGroup.GET("/adminPosts", admin.AdminPostsPage)
		adminGroup.DELETE("/posts/:id", admin.DeletePostByAdmin)
		adminGroup.DELETE("/messages/:id", admin.DeleteMessageByAdmin)
		adminGroup.GET("/messages", admin.AdminMessagesPage)
	}

	// Обработка методов DELETE через POST
	router.Use(func(c *gin.Context) {
		if c.Request.Method == "POST" && c.PostForm("_method") == "DELETE" {
			c.Request.Method = "DELETE"
		}
		c.Next()
	})
	

	// Запуск сервера
	router.Run("217.194.148.61:8080")
}