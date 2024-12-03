package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// model

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}
type TodoRequest struct {
	Title string `json:"title"`
}

type TodoUpdateRequest struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var todos []Todo

func main() {
	app := fiber.New()

	apiV1 := app.Group("/api/v1")

	// MIDDLEWARE
	app.Use(func(c *fiber.Ctx) error {
		t := time.Now().Format("15:04:05 AM, 02 Jan 2006")
		fmt.Println("Hello From Middleware", t)
		return c.Next()
	})

	app.Static("/", "./public") // static

	apiV1.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(&fiber.Map{
			"success": true,
			"message": "Server is Healthy 🌴",
		})
	})

	// get all todos
	apiV1.Get("/todos", authCheck, func(c *fiber.Ctx) error {
		if todos == nil {
			return c.Status(200).JSON(&fiber.Map{
				"success": true,
				"message": "No Todos Found",
				"todos":   []Todo{},
			})
		}
		return c.Status(200).JSON(&fiber.Map{
			"success": true,
			"message": "Todo Fetched Successfully ✅",
			"todos":   todos,
		})
	})

	// add todo
	apiV1.Post("/todos", authCheck, func(c *fiber.Ctx) error {

		var todoReq TodoRequest
		var todo Todo

		if err := c.BodyParser(&todoReq); err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": "Invalid Request",
			})
		}

		todo.ID = len(todos) + 1
		todo.Title = todoReq.Title
		todo.Done = false

		todos = append(todos, todo)

		return c.Status(201).JSON(&fiber.Map{
			"success": true,
			"message": "Todo Added Successfully ✅",
			"todo":    todo,
		})
	})

	// get todo by id
	apiV1.Get("/todos/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")

		if err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": "Invalid Request",
			})
		}

		for _, todo := range todos {
			if todo.ID == id {
				return c.Status(200).JSON(&fiber.Map{
					"success": true,
					"message": "Todo Fetched Successfully ✅",
					"todo":    todo,
				})
			}
		}
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": "Todo Not Found",
		})
	})

	// update todo
	apiV1.Patch("/todos/:id", authCheck, func(c *fiber.Ctx) error {

		var todoReq TodoUpdateRequest

		if err := c.BodyParser(&todoReq); err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": "Invalid Request",
			})
		}

		id, err := c.ParamsInt("id")

		if err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": "Invalid Request",
			})
		}

		for i, todo := range todos {
			if todo.ID == id {
				if todoReq.Title != "" {
					todos[i].Title = todoReq.Title
				}
				todos[i].Done = todoReq.Done

				return c.Status(200).JSON(&fiber.Map{
					"success": true,
					"message": "Todo Updated Successfully ✅",
					"todo":    todos[i],
				})
			}
		}
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": "Todo Not Found",
		})
	})

	// delete todo
	apiV1.Delete("/todos/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")

		if err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": "Invalid Request",
			})
		}

		for i, todo := range todos {
			if todo.ID == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(&fiber.Map{
					"success": true,
					"message": "Todo Deleted Successfully ✅",
				})
			}
		}
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": "Todo Not Found",
		})
	})

	log.Fatal(app.Listen(":4000"))
}

func authCheck(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	authType := strings.Split(auth, " ")[0]

	if authType != "Bearer" {
		return c.Status(401).JSON(&fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}

	token := strings.Split(auth, " ")[1]

	if strings.TrimSpace(token) == "" {
		return c.Status(401).JSON(&fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}

	c.Request().Header.Add("token", token)

	return c.Next()
}
