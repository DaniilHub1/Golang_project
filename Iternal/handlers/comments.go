package handlers

import (
    "database/sql"
    "mini_site/models" 

	)

func GetCommentsByPostID(db *sql.DB, postID int) ([]models.Comment, error) {
    rows, err := db.Query("SELECT id, post_id, user_id, content, created_at FROM comments WHERE post_id=$1 ORDER BY created_at", postID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var comments []models.Comment
    for rows.Next() {
        var c models.Comment
        if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt); err != nil {
            return nil, err
        }
        comments = append(comments, c)
    }
    return comments, nil
}


