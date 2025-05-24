package handlers

	import (
		"fmt"
		"mini_site/database"
		"mini_site/models"
		"net/http"
		"github.com/gin-gonic/gin"
		"golang.org/x/crypto/bcrypt"
	)

	func IsAuthenticated(c *gin.Context) bool {
		token, err := c.Cookie("token")
		return err == nil && token != ""
	}

	func RegisterHandler(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user := models.User{Username: username, Password: string(hashedPassword), Role: "user"}
		result := database.DB.Create(&user)
		if result.Error != nil {
			c.Redirect(http.StatusFound, "/register?error=User+already+exists")
			return
		}

		c.Redirect(http.StatusFound, "/login?success=Registration+successful")
	}


	func LoginHandler(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		var user models.User
		result := database.DB.Where("username = ?", username).First(&user)
		if result.Error != nil {
			fmt.Println("❌ Пользователь не найден:", username)
			c.Redirect(http.StatusFound, "/login?error=Invalid+username+or+password")
			return
		}

		fmt.Println(" Пользователь найден:", user.Username)
		fmt.Println(" Хеш из БД:", user.Password)
		fmt.Println(" Проверяем пароль:", password)

		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			fmt.Println("❌ Пароль не совпадает")
			c.Redirect(http.StatusFound, "/login?error=Invalid+username+or+password")
			return
		}

		c.SetCookie("token", user.Username, 3600, "/", "", false, true)
		c.SetCookie("username", user.Username, 3600, "/", "", false, true)
		c.SetCookie("user_role", user.Role, 3600, "/", "", false, true)

		if user.Role == "admin" {
			c.Redirect(http.StatusFound, "/admin/dashboard")
		} else {
			c.Redirect(http.StatusFound, "/account")
		}
	}

	func LogoutHandler(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.SetCookie("username", "", -1, "/", "", false, true)
	c.SetCookie("user_role", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login?success=Вы+вышли+из+аккаунта")
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token") 
		if err != nil || token == "" {
			// проверяем заголовок Authorization, как сейчас
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		var user models.User
		if err := database.DB.Where("username = ?", token).First(&user).Error; err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// сохраняем в контекст
		c.Set("user", user)
		c.Set("userID", user.ID)
		c.Set("userRole", user.Role)
		c.Set("username", user.Username)

		c.Next()
	}
}



	// Показ формы логина с алертами
	func LoginPage(c *gin.Context) {
		errorMsg := c.Query("error")
		successMsg := c.Query("success")

		c.HTML(http.StatusOK, "login.html", gin.H{
			"Error":   errorMsg,
			"Success": successMsg,
		})
	}

	func InitAdminUser() {
		var admin models.User
		database.DB.First(&admin, "username = ?", "admin")
		if admin.ID == 0 {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
			if err != nil {
				panic("Ошибка хеширования пароля админа")
			}
			admin = models.User{
				Username: "admin",
				Password: string(hashedPassword),
				Role:     "admin",
			}
			database.DB.Create(&admin)
			fmt.Println("✅ Админ создан: admin / admin123")
			fmt.Println("Хеш пароля:", admin.Password)
		} else {
			fmt.Println("⚠️ Админ уже существует")
		}
	}