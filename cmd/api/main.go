package main

import (
	"mini_site/Iternal/admin"
	"mini_site/Iternal/handlers"
	"mini_site/database"
	"mini_site/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		panic("failed to connect to the database")
	}

	db.AutoMigrate(&models.User{}, &models.Message{})
	handlers.DB = db
	

	router := gin.Default()
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/**/*.html")

	// Страницы и формы
	router.GET("/", handlers.HomePage)
	router.GET("/about", handlers.AboutPage)
	router.GET("/account", handlers.RenderPage2)
	router.GET("/index.html", handlers.RenderPage3)

	router.GET("/chat", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", nil)
	})

	router.GET("/posts", handlers.GetPosts)
	router.POST("/posts_page", handlers.CreatePostFromForm)
	router.PUT("/posts/:id", handlers.UpdatePost)
	router.DELETE("/posts/:id", handlers.DeletePost)
	router.GET("/posts_page", handlers.RenderPostsPage)

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

	router.GET("/settings", handlers.SettingsPage)
	router.POST("/settings", handlers.SaveSettings)

	router.POST("/send", handlers.SendMessage)
	router.GET("/messages", handlers.GetMessages)
	router.GET("/users", handlers.GetUsers)
	router.POST("/profile/save", handlers.SaveSettings)

	//  Админ-панель
	adminGroup := router.Group("/admin")
	adminGroup.Use(handlers.AdminRequired())
	{
		adminGroup.GET("/dashboard", admin.AdminDashboard)

		// роут: удалить пользователя по ID
		adminGroup.DELETE("/users/:id", admin.DeleteUserByAdmin)
		adminGroup.GET("/users", admin.AdminUsersPage)
		adminGroup.GET("/adminPosts", admin.AdminPostsPage)
		adminGroup.DELETE("/posts/:id", admin.DeletePostByAdmin)
		adminGroup.DELETE("/messages/:id", admin.DeleteMessageByAdmin)

	}

	// Посты
	router.GET("/feed", handlers.RenderPostsFeed)

	router.Run("217.194.148.61:8080")
}
