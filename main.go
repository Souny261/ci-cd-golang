package main

import (
	"ci-cd-golang/config"
	"ci-cd-golang/controllers"
	"ci-cd-golang/database"
	"ci-cd-golang/models"
	"encoding/json"
	"fmt"

	minioUpload "ci-cd-golang/platform/minio"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	minioClient, err := minioUpload.MinioConnection()
	_ = minioClient
	app := fiber.New(
		fiber.Config{
			JSONEncoder: json.Marshal,
			JSONDecoder: json.Unmarshal,
		})
	app.Use(logger.New())
	app.Use(cors.New())
	sostgresConnection, err := database.PostgresConnection()
	if err != nil {
		fmt.Printf("Connect database error %v ", err)
		return
	}

	app.Get("/minio", func(c *fiber.Ctx) error {
		minioClient, err := minioUpload.MinioConnection()
		_ = minioClient
		if err != nil {
			return c.JSON(fiber.Map{
				"status": false,
				"msg":    err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"status": true,
			"msg":    "Connect to Minio Successfully",
		})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println("GET: Hello World")
		return c.JSON(fiber.Map{
			"status": true,
			"msg":    "GET: Hello World2023",
		})
	})

	app.Post("/create", func(c *fiber.Ctx) error {
		createUserRequest := CreateUserRequest{}
		err := c.BodyParser(&createUserRequest)
		if err != nil {
			return c.JSON(fiber.Map{
				"status": false,
				"msg":    err.Error(),
			})
		}
		password, _ := bcrypt.GenerateFromPassword([]byte(createUserRequest.Password), 14)
		payload := models.User{
			Name:     createUserRequest.Name,
			Email:    createUserRequest.Email,
			Password: password,
		}
		err = sostgresConnection.Create(&payload).Error

		if err != nil {
			return c.JSON(fiber.Map{
				"status": false,
				"msg":    err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"status": true,
			"msg":    "Successfully",
			"data":   payload,
		})
	})

	app.Get("/user/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.JSON(fiber.Map{
				"status": false,
				"msg":    err.Error(),
			})
		}
		var user models.User
		err = sostgresConnection.Where("id = ?", id).First(&user).Error
		if err != nil {
			return c.JSON(fiber.Map{
				"status": false,
				"msg":    err.Error(),
			})
		}
		res := UserResponse{
			Name:  user.Name,
			Id:    user.Id,
			Email: user.Email,
		}
		return c.JSON(fiber.Map{
			"status": true,
			"msg":    "Successfully",
			"data":   res,
		})
	})

	app.Get("/users", func(c *fiber.Ctx) error {
		var users []models.User
		sostgresConnection.Order("id desc").Find(&users)
		return c.JSON(fiber.Map{
			"status": true,
			"msg":    "Successfully",
			"data":   users,
		})
	})
	app.Post("/upload", controllers.UploadFile)
	app.Get("/read/:name?", controllers.GetFile)
	MYPORT := config.GetEnv("app.port", "3000")
	SERVER_RUNNING := fmt.Sprintf(":%v", MYPORT)
	app.Listen(SERVER_RUNNING)
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserResponse struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
