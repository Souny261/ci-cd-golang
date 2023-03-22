package main

import (
	"ci-cd-golang/config"
	"ci-cd-golang/database"
	"ci-cd-golang/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	app := fiber.New()
	sostgresConnection, err := database.PostgresConnection()
	if err != nil {
		fmt.Printf("Connect database error %v ", err)
		return
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": true,
			"msg":    "Hello World",
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

	app.Post("/user", func(c *fiber.Ctx) error {
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
		// err = sostgresConnection.Create(&payload).Error

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
