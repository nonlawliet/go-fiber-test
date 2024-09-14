package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func getBooks(c *fiber.Ctx) error {
	return c.JSON(books)
}

func getBook(c *fiber.Ctx) error {
	// #1 - Get bookId from request
	bookId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// #2 - Query book
	for _, book := range books {
		if book.ID == bookId {
			return c.JSON(book)
		}
	}

	return c.SendStatus(fiber.StatusNotFound)
}

func createBook(c *fiber.Ctx) error {
	// #1 - Create book instance, input body parser into it
	book := new(Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// #2 - Append book into books (analogy like update to db)
	books = append(books, *book)

	// #3 - Make response
	return c.JSON(book)
}

func updateBook(c *fiber.Ctx) error {
	// #1 - Get bookId from request
	bookId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// #2 - Create bookUpdate instance
	bookUpdate := new(Book)
	if err := c.BodyParser(bookUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// #3 - Query book
	for i, book := range books {
		if book.ID == bookId {
			books[i].Title = bookUpdate.Title
			books[i].Author = bookUpdate.Author

			// #4 - Make response
			return c.JSON(books[i])
		}
	}

	return c.SendStatus(fiber.StatusNotFound)
}

func deleteBook(c *fiber.Ctx) error {
	// #1 - Get bookId from request
	bookId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// #2 - Query book
	for i, book := range books {
		if book.ID == bookId {

			// #3 - Get book (which want to delete) out of slice (anaoly to delete from db)
			books = append(books[:i], books[i+1:]...)
			// ex. we have slice [1, 2, 3, 4, 5]
			// [1, 2] + [4, 5] = [1, 2, 4, 5]

			// #4 - Make response
			return c.SendStatus(fiber.StatusNoContent)
		}
	}

	return c.SendStatus(fiber.StatusNotFound)
}
