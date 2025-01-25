package main

import (
	"eventix/config"
	"eventix/controller"
	_ "eventix/docs" // Import Swagger docs
	"eventix/middleware"
	"eventix/repository"
	"eventix/service"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// Set zona waktu default
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}
	time.Local = loc

	// Inisialisasi database
	db := config.DBInit()

	// Dependency Injection
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	eventRepo := repository.NewEventRepository(db)
	ticketRepo := repository.NewTicketRepository(db)
	eventService := service.NewEventService(eventRepo, ticketRepo)
	eventController := controller.NewEventController(eventService)

	ticketService := service.NewTicketService(ticketRepo, eventRepo)
	ticketController := controller.NewTicketController(ticketService)

	reportService := service.NewReportService(ticketRepo)
	reportController := controller.NewReportController(reportService)

	blacklistRepo := repository.NewTokenBlacklistRepository(db)
	authService := service.NewAuthService(userRepo, blacklistRepo)
	authController := controller.NewAuthController(authService)

	exportController := controller.NewExportController(reportService)

	// Setup Router
	r := gin.Default()

	// Rute Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Route publik tanpa middleware
	r.POST("/users/register", userController.RegisterUser)

	// Middleware untuk autentikasi
	r.Use(middleware.AuthenticationMiddleware("your_secret_key"))

	// Route untuk logout
	r.POST("/login", authController.Login)
	r.POST("/logout", authController.Logout)

	// Routes untuk pengguna umum (User)
	r.GET("/events", middleware.AuthorizeRole("User"), eventController.GetAllEvents)
	r.GET("/events/:id", middleware.AuthorizeRole("User"), eventController.GetEventByID)
	r.GET("/tickets", middleware.AuthorizeRole("User"), ticketController.GetTickets)
	r.POST("/tickets", middleware.AuthorizeRole("User"), ticketController.CreateTicket)
	r.PATCH("/tickets/:id", middleware.AuthorizeRole("User"), ticketController.CancelTicket)

	// Routes untuk admin
	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.AuthorizeRole("Admin"))
	adminRoutes.POST("/events", eventController.CreateEvent)
	adminRoutes.PUT("/events/:id", eventController.UpdateEvent)
	adminRoutes.DELETE("/events/:id", eventController.DeleteEvent)
	adminRoutes.GET("/reports/summary", reportController.GetSummaryReport)
	adminRoutes.GET("/reports/event/:id", reportController.GetEventReport)

	adminRoutes.GET("/export/reports/summary", exportController.ExportSummaryReport)
	adminRoutes.GET("/export/reports/event/:id", exportController.ExportEventReport)

	// Jalankan server
	r.Run(":8080")
}
