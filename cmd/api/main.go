package main

import (
    "mini_site/handlers"
    "mini_site/database"
    "github.com/gin-gonic/gin"
    "net/http"
)

func main() {
    db, err := database.InitDB()
    if err != nil {
        panic("failed to connect to the database")
    }

    handlers.DB = db

    router := gin.Default()    

    router.Static("/static", "./static")
    router.LoadHTMLGlob("templates/*.html")

    router.GET("/posts", handlers.GetPosts)
    router.POST("/posts", handlers.CreatePostFromForm)
    router.PUT("/posts/:id", handlers.UpdatePost)
    router.DELETE("/posts/:id", handlers.DeletePost)
    router.GET("/posts_page", handlers.RenderPostsPage)

    router.GET("/", handlers.HomePage)
    router.GET("/about", handlers.AboutPage)
    router.GET("/account", handlers.RenderPage2)
    router.GET("/index.html", handlers.RenderPage3)

    router.GET("/login", func(c *gin.Context) {
        c.HTML(http.StatusOK, "login.html", nil)
    })
    router.POST("/login", handlers.LoginHandler) 

    router.GET("/register", func(c *gin.Context) {
        c.HTML(http.StatusOK, "register.html", nil)
    })
    router.POST("/register", handlers.RegisterHandler) 

    router.GET("/settings", handlers.SettingsPage)
    router.POST("/settings", handlers.SaveSettings)
    router.POST("/profile/save", handlers.SaveProfileSettings)

    router.Run("217.194.148.61:8080")
}
