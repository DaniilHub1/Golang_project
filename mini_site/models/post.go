package models

import "time"

type Post struct {
    ID       int
    Title    string
    Content  string  
    Author   string
    Date     time.Time
    Username string 
}
