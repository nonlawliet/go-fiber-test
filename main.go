package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	// มันคือ app = Express() // ตัวสื่อสารระหว่างแอพกับ Server

	app.Get("/hello", func(c *fiber.Ctx) error { //
		return c.SendString("Hello World!")
	})
	/* คืน resp ปกติเมื่อ function ไม่มี error, แต่ถ้ามี error ออกมา fiber จะจัดการให้
	โดยหยุดการทำงานของฟังก์ชัน และ throw error ออกมา */

	app.Listen(":8080") // ระบุให้ app listen port นี้
}
