package routes

import (
	"hirehound/handlers"
	"hirehound/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/login", handlers.Login)
	auth.Post("/register", handlers.Register)
	auth.Get("/logout", handlers.Logout)
	// auth.Get("/reset", handlers.ResetPassword)

	auth.Get("/valid", middleware.Protected())
}

func ApplicationRoutes(app *fiber.App) {
	applications := app.Group("/applications", middleware.Protected())

	// Get all applications
	applications.Get("/", handlers.GetAllApplications)
	// Get application by id
	applications.Get("/:id", handlers.GetApplicationByID)

	// Create application
	applications.Post("/", handlers.CreateApplication)

	// Update application by id
	applications.Put("/:id", handlers.UpdateApplication)

	// Delete application
	applications.Delete("/:id", handlers.DeleteApplication)
}
