package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Mshivam2409/AnC-Courses/models"
	"github.com/Mshivam2409/AnC-Courses/services"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateFile sdfs
func CreateFile(c *fiber.Ctx) error {
	d := services.GetService()
	n, err := strconv.Atoi(c.FormValue("length"))
	cno := c.FormValue("course")
	if err != nil {
		return &fiber.Error{Code: 422, Message: "Length of files must be a number"}
	}
	files := []string{}
	for i := 0; i < n; i++ {
		f, err := c.FormFile("file" + fmt.Sprint(n))
		if err != nil {
			return &fiber.Error{Code: 422, Message: "Unable to process file!"}
		}
		r, err := d.CreateFile("", f)
		if err != nil {
			return &fiber.Error{Code: 422, Message: "Unable to upload file!"}
		}
		out, err := json.Marshal(r)
		if err != nil {
			return &fiber.Error{Code: 422, Message: "Unable to get file ID!"}
		}
		files = append(files, string(out))
	}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "driveFiles", Value: bson.A{"$each", files}},
		}},
	}
	err = services.MongoClient.Courses.Collection("courses").FindOneAndUpdate(context.TODO(), bson.D{{Key: "number", Value: cno}}, update).Decode(&models.MGMCourse{})
	return c.Status(200).JSON(fiber.Map{"message": "success"})
}

func DownloadFile(ctx *fiber.Ctx) error {
	d := services.GetService()
	s, err := d.GetFile(ctx.Params("fid"))
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Unable to download file!"}
	}
	return ctx.Status(200).SendStream(s)
}

func Register(ctx *fiber.Ctx) error {
	u := strings.ReplaceAll(ctx.FormValue("username"), " ", "")
	log.Println(u)
	err := services.OryClient.CreateUser(u)
	if err != nil {
		return fiber.NewError(401, fmt.Sprint(err))
	}
	return ctx.Status(200).JSON(fiber.Map{})
}

func Hydrate(ctx *fiber.Ctx) error {
	log.Print(ctx)
	a := &models.AuthenticationSession{}
	err := ctx.BodyParser(a)
	if err != nil {
		log.Print(err)
		return err
	}
	a.Header.Add("role", "user")
	log.Print(a.Header)
	// gqlparser.LoadS
	n := *a
	return ctx.Status(200).JSON(n)
}
