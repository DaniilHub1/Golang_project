package admin

import (
	"mini_site/database"
	"mini_site/models"
	"fmt"
)

type AdminStats struct {
	UserCount    int64
	PostCount    int64
	MessageCount int64
}

func GetAdminStats() (AdminStats, error) {
	var stats AdminStats

	if err := database.DB.Model(&models.User{}).Count(&stats.UserCount).Error; err != nil {
		return stats, fmt.Errorf("error fetching user count: %v", err)
	}

	if err := database.DB.Model(&models.Post{}).Count(&stats.PostCount).Error; err != nil {
		return stats, fmt.Errorf("error fetching post count: %v", err)
	}

	if err := database.DB.Model(&models.Message{}).Count(&stats.MessageCount).Error; err != nil {
		return stats, fmt.Errorf("error fetching message count: %v", err)
	}

	return stats, nil
}
