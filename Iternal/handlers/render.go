package handlers

import (
	"github.com/gin-gonic/gin"
	"mini_site/database"
	"mini_site/models"
	"net/http"
	"html/template"
)

type AccountData struct {
	Username string
	Email    string
	Nickname string
	Avatar   string
}

func RenderPage2(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil || token == "" {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", token).First(&user).Error; err != nil {
		c.String(http.StatusUnauthorized, "Ошибка: пользователь не найден.")
		return
	}

	accountData := AccountData{
		Username: user.Username,
		Email:    user.Email,
		Nickname: user.Nickname,
		Avatar:   user.PhotoPath,
	}

	tmpl, err := template.ParseFiles("templates/user/account.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка загрузки шаблона: %v", err)
		return
	}

	tmpl.Execute(c.Writer, accountData)
}

func RenderPage3(c *gin.Context) {
	tmpl, err := template.ParseFiles("templates/user/index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка загрузки шаблона: %v", err)
		return
	}

	tmpl.Execute(c.Writer, nil)
}

