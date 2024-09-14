package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// Book struct to hold book data
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book

func main() {
	app := fiber.New()

	if err := godotenv.Load(); err != nil {
		log.Fatal("load .env error")
	}

	// Initialize in-memory data
	books = append(books, Book{ID: 1, Title: "1984", Author: "George Orwell"})
	books = append(books, Book{ID: 2, Title: "The Great Gatsby", Author: "F. Scott Fitzgerald"})

	app.Get("/books", getBooks)
	app.Get("/books/:id", getBook)
	app.Post("/books", createBook)
	app.Put("/books/:id", updateBook)
	app.Delete("/books/:id", deleteBook)

	app.Post("/upload", uploadFile)

	app.Get("/config", getEnv)

	app.Listen(":8080")
}

func uploadFile(c *fiber.Ctx) error {
	// #1 - Get file from request in form of FormData
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// #2 - Save file into local and set path
	err = c.SaveFile(file, "./uploads/"+file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// #3 - Make response
	return c.SendString("File upload complete!")
}

func getEnv(c *fiber.Ctx) error {
	// After the main function has "godotenv.Load()", the Getenv can used through os.
	// #1 - get environment variable name "SECRET"
	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "default secret"
	}

	// #2 - Make response
	return c.JSON(fiber.Map{
		"SECRET": secret,
	})
}
