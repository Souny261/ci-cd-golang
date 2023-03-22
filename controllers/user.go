package controllers

// import (
// 	"ci-cd-golang/config"
// 	"context"
// 	"net/url"
// 	"time"

// 	"github.com/gofiber/fiber/v2"
// )

// func GetUsers(c *fiber.Ctx) error {
// 	// Create minio connection.
// 	minioClient, err := minioUpload.MinioConnection()
// 	bucketName := config.GetEnv("minio.bucket", "")
// 	reqParams := make(url.Values)
// 	presignedURL, err := minioClient.PresignedGetObject(context.Background(), bucketName, "B-Green.png", time.Duration(1000)*time.Second, reqParams)
// 	if err != nil {
// 		return c.JSON(fiber.Map{
// 			"error": false,
// 			"msg":   err.Error(),
// 		})
// 	}
// 	return c.JSON(fiber.Map{
// 		"error": false,
// 		"msg":   nil,
// 		"url":   presignedURL.String(),
// 	})
// }
