package controllers

import (
	"ci-cd-golang/config"
	minioUpload "ci-cd-golang/platform/minio"
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

func UploadFile(c *fiber.Ctx) error {
	payload := struct {
		Name       string `json:"name"`
		FileUpload string `json:"fileUpload"`
	}{}

	c.BodyParser(&payload)
	fmt.Println(payload.Name)

	ctx := context.Background()
	bucketName := config.GetEnv("minio.bucket", "")
	file, err := c.FormFile("fileUpload")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get Buffer from file
	buffer, err := file.Open()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	defer buffer.Close()

	// Create minio connection.
	minioClient, err := minioUpload.MinioConnection()
	if err != nil {
		// Return status 500 and minio connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	objectName := file.Filename
	fileBuffer := buffer
	contentType := file.Header["Content-Type"][0]
	fileSize := file.Size

	// Upload the zip file with PutObject
	info, err := minioClient.PutObject(ctx, bucketName, objectName, fileBuffer, fileSize, minio.PutObjectOptions{ContentType: contentType})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"info":  info,
		"name":  payload.Name,
	})
}
func GetFile(c *fiber.Ctx) error {
	fileName := c.Params("name")

	// Create minio connection.
	minioClient, err := minioUpload.MinioConnection()
	bucketName := config.GetEnv("minio.bucket", "")
	reqParams := make(url.Values)
	presignedURL, err := minioClient.PresignedGetObject(context.Background(), bucketName, fileName, time.Duration(1000)*time.Second, reqParams)
	if err != nil {
		return c.JSON(fiber.Map{
			"error": false,
			"msg":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"error": true,
		"msg":   "",
		"url":   presignedURL.String(),
		"name":  fileName,
	})
}
