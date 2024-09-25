package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/xuri/excelize/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db          *gorm.DB
	redisClient *redis.Client
	ctx         = context.Background()
)

type User struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	City        string `json:"city"`
	County      string `json:"county"`
	Postal      string `json:"postal"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Web         string `json:"web"`
}

func initDB() {
	var err error
	dsn := "root:ebax5010IM@1@tcp(127.0.0.1:3306)/my_new_database?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	db.AutoMigrate(&User{})
}

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func uploadExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	if err := c.SaveUploadedFile(file, file.Filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	go processExcelFile(file.Filename)
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}

func processExcelFile(filename string) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		log.Printf("failed to open file: %v", err)
		return
	}

	rows, err := f.GetRows("uk-500")
	if err != nil {
		log.Printf("failed to get rows: %v", err)
		return
	}

	for _, row := range rows[1:] {
		user := User{
			FirstName:   row[0],
			LastName:    row[1],
			CompanyName: row[2],
			Address:     row[3],
			City:        row[4],
			County:      row[5],
			Postal:      row[6],
			Phone:       row[7],
			Email:       row[8],
			Web:         row[9],
		}

		db.Create(&user)
		redisClient.Set(ctx, fmt.Sprintf("user:%d", user.ID), user, 5*time.Minute)
	}
}

func getUsers(c *gin.Context) {
	var users []User
	result, err := redisClient.Get(ctx, "users").Result()
	if err == nil {
		json.Unmarshal([]byte(result), &users)
		c.JSON(http.StatusOK, users)
		return
	}

	db.Find(&users)
	redisClient.Set(ctx, "users", users, 5*time.Minute)
	c.JSON(http.StatusOK, users)
}

func updateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32) // Convert string to uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user.ID = uint(id)

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	redisClient.Set(ctx, fmt.Sprintf("user:%d", user.ID), user, 5*time.Minute)
	c.JSON(http.StatusOK, user)
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	db.Delete(&User{}, id)
	redisClient.Del(ctx, fmt.Sprintf("user:%s", id))
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func main() {
	initDB()
	initRedis()

	r := gin.Default()
	r.POST("/upload", uploadExcel)
	r.GET("/users", getUsers)
	r.PUT("/users/:id", updateUser)
	r.DELETE("/users/:id", deleteUser)

	r.Run(":8080")
}
