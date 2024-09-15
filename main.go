package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	jwtware "github.com/gofiber/jwt/v2"
	"github.com/gofiber/swagger"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/nonlawliet/go-fiber-test/docs" // load generated docs
)

// Book struct to hold book data
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book

// @title Book API
// @description This is a sample server for a book API.
// @version 1.0
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	app := fiber.New()

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	// Init .env data into main function
	if err := godotenv.Load(); err != nil {
		log.Fatal("load .env error")
	}

	// Initialize in-memory data
	books = append(books, Book{ID: 1, Title: "1984", Author: "George Orwell"})
	books = append(books, Book{ID: 2, Title: "The Great Gatsby", Author: "F. Scott Fitzgerald"})

	/* concept of login and check token (ex. of middleware)
	1. check user email and password, create token if correct
	2. loggin the URL, Method, Time
	3. check the signingKey of request, show error "Missing or malformed JWT" if the token isn't correct
	*/
	// Init middle ware (login)
	app.Post("/login", login)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("JWT_SECRET"),
	}))

	// Init middle ware (logging) bf the other handlers function
	app.Use(checkMiddleware)

	/* if the login and check token are passed, go to other api process */
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

func checkMiddleware(c *fiber.Ctx) error {
	start := time.Now()
	fmt.Printf("URL = %s, Method = %s, Time = %s\n", c.OriginalURL(), c.Method(), start)

	// read user context local (that write by jwtware) and convert into token
	user := c.Locals("user").(*jwt.Token)

	// get token data by claim
	claims := user.Claims.(jwt.MapClaims)

	// check role in token
	if claims["role"] != "admin" {
		return fiber.ErrUnauthorized
	}

	return c.Next()
}

// Dummy user for example
type User = struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var memberUser = User{
	Email:    "user@example.com",
	Password: "password123",
}

func login(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if user.Email != memberUser.Email || user.Password != memberUser.Password {
		return fiber.ErrUnauthorized
	}

	// Create token // Set pattern token as a MethodHS256
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims // Set token data
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["role"] = "admin"
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("JWT_SECRET"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "Login success",
		"token":   t,
	})
}
