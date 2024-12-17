package handlers

import (
    "encoding/json"
    "fmt"
    "log"
    "mini_site/models" 
    "github.com/go-redis/redis/v8"
)

func SavePostToRedis(post models.Post) error {
    postData, err := json.Marshal(post)
    if err != nil {
        return fmt.Errorf("не удалось преобразовать пост в JSON: %v", err)
    }

    key := fmt.Sprintf("post:%d", post.ID) 
    return rdb.Set(ctx, key, postData, 0).Err()
}

func GetPostFromRedis(postID int) (*models.Post, error) {
    key := fmt.Sprintf("post:%d", postID)
    postData, err := rdb.Get(ctx, key).Result()
    if err == redis.Nil {
        return nil, fmt.Errorf("пост с ID %d не найден", postID)
    } else if err != nil {
        return nil, fmt.Errorf("ошибка при получении поста: %v", err)
    }

    var post models.Post
    if err := json.Unmarshal([]byte(postData), &post); err != nil {
        return nil, fmt.Errorf("ошибка при декодировании JSON: %v", err)
    }
    return &post, nil
}

func GetAllPostsFromRedis() ([]models.Post, error) {
    var posts []models.Post

    keys, err := rdb.Keys(ctx, "post:*").Result()
    if err != nil {
        return nil, fmt.Errorf("ошибка при получении ключей постов: %v", err)
    }

    for _, key := range keys {
        postData, err := rdb.Get(ctx, key).Result()
        if err != nil {
            log.Printf("Ошибка получения поста с ключом %s: %v", key, err)
            continue
        }

        var post models.Post
        if err := json.Unmarshal([]byte(postData), &post); err != nil {
            log.Printf("Ошибка декодирования поста с ключом %s: %v", key, err)
            continue
        }
        posts = append(posts, post)
    }

    return posts, nil
}

func DeletePostFromRedis(postID int) error {
    key := fmt.Sprintf("post:%d", postID) 
    return rdb.Del(ctx, key).Err()       
}