package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mini_site/database"
	"mini_site/models"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SettingsPage(c *gin.Context) {
	nickname, avatar, _ := LoadProfileSettings(c)

	c.HTML(http.StatusOK, "settings.html", gin.H{
		"nickname": nickname,
		"avatar":   avatar,
	})
}

func SaveSettings(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	nickname := c.PostForm("nickname")

	token, err := c.Cookie("token")
	if err != nil || token == "" {
		fmt.Println("[SaveSettings] Ошибка: токен не найден")
		c.HTML(http.StatusUnauthorized, "settings.html", gin.H{
			"error": "Вы не авторизованы. Пожалуйста, войдите в систему.",
		})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", token).First(&user).Error; err != nil {
		fmt.Println("[SaveSettings] Ошибка: пользователь не найден по токену =", token)
		c.HTML(http.StatusUnauthorized, "settings.html", gin.H{
			"error": "Пользователь не найден. Возможно, вы не авторизованы.",
		})
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		fmt.Println("[SaveSettings] Нет файла:", err)
	} else {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.HTML(http.StatusBadRequest, "settings.html", gin.H{
				"error": "Только JPG/PNG файлы разрешены",
			})
			return
		}

		uploadDir := "./static/avatars" 
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.HTML(http.StatusInternalServerError, "settings.html", gin.H{
				"error": "Ошибка сервера при создании папки",
			})
			return
		}

		if user.PhotoPath != "" && !strings.Contains(user.PhotoPath, "default-avatar") {
			_ = os.Remove("." + user.PhotoPath)
		}

		uniqueFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		avatarPath := filepath.Join(uploadDir, uniqueFilename)

		if err := c.SaveUploadedFile(file, avatarPath); err != nil {
			fmt.Println("[SaveSettings] Ошибка загрузки файла:", err)
			c.HTML(http.StatusInternalServerError, "settings.html", gin.H{
				"error": "Ошибка загрузки файла",
			})
			return
		}

		user.PhotoPath = "/static/avatars/" + uniqueFilename
		fmt.Println("[SaveSettings] Аватарка загружена:", user.PhotoPath)
	}

	user.Username = username
	user.Email = email
	user.Nickname = nickname

	if err := database.DB.Save(&user).Error; err != nil {
		fmt.Println("[SaveSettings] Ошибка при сохранении в БД:", err)
		c.HTML(http.StatusInternalServerError, "settings.html", gin.H{
			"error": "Ошибка сохранения профиля",
		})
		return
	}

	fmt.Println("[SaveSettings] Профиль успешно сохранён для пользователя:", user.Username)

	c.HTML(http.StatusOK, "settings.html", gin.H{
		"success":  "Настройки сохранены!",
		"nickname": user.Nickname,
		"avatar":   user.PhotoPath,
		"username": user.Username,
		"email":    user.Email,
	})
}



func LoadProfileSettings(c *gin.Context) (string, string, error) {
	token, err := c.Cookie("token")
	if err != nil || token == "" {
		return "Гость", "/static/default-avatar.png", nil
	}

	var user models.User
	if err := database.DB.Where("username = ?", token).First(&user).Error; err != nil {
		return "Гость", "/static/default-avatar.png", nil
	}

	nickname := user.Nickname
	if nickname == "" {
		nickname = "Гость"
	}

	avatar := user.PhotoPath
	if avatar == "" {
		avatar = "/static/default-avatar.png"
	}

	return nickname, avatar, nil
}
