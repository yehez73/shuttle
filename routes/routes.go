package routes

import (
	"shuttle/handler"
	"shuttle/middleware"
	"shuttle/repositories"
	"shuttle/services"
	"shuttle/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/contrib/websocket"
	"github.com/jmoiron/sqlx"
)

func Route(r *fiber.App, db *sqlx.DB) {
	authRepository := repositories.NewAuthRepository(db)
	userRepository := repositories.NewUserRepository(db)
	schoolRepository := repositories.NewSchoolRepository(db)
	vehicleRepository := repositories.NewVehicleRepository(db)
	studentRepository := repositories.NewStudentRepository(db)
	routeRepository := repositories.NewRouteRepository(db)
	childernRepository := repositories.NewChildernRepository(db)
	shuttleRepository := repositories.NewShuttleRepository(db)
	
	userService := services.NewUserService(userRepository)
	authService := services.NewAuthService(authRepository, userRepository)
	schoolService := services.NewSchoolService(schoolRepository, userRepository)
	vehicleService := services.NewVehicleService(vehicleRepository)
	studentService := services.NewStudentService(studentRepository, &userService, userRepository)
	routeService := services.NewRouteService(routeRepository)
	childernService := services.NewChildernService(childernRepository)
	shuttleService := services.NewShuttleService(shuttleRepository)
	
	authHandler := handler.NewAuthHttpHandler(authService, userService)
	userHandler := handler.NewUserHttpHandler(userService, schoolService, vehicleService)
	schoolHandler := handler.NewSchoolHttpHandler(schoolService)
	vehicleHandler := handler.NewVehicleHttpHandler(vehicleService)
	studentHandler := handler.NewStudentHttpHandler(studentService)
	routeHandler := handler.NewRouteHttpHandler(routeService)
	childernHandler := handler.NewChildernHandler(childernService)
	shuttleHandler := handler.NewShuttleHandler(shuttleService)

	wsService := utils.NewWebSocketService(userRepository, authRepository, shuttleRepository)
	
	////////////////////////////////////// PUBLIC //////////////////////////////////////

	r.Post("login", authHandler.Login)
	r.Post("/refresh-token", authHandler.IssueNewAccessToken)
	r.Static("/assets", "./assets")

	////////////////////////////////////// AUTHENTICATED //////////////////////////////////////

	protected := r.Group("/api")
	protected.Use(middleware.AuthenticationMiddleware())
	protected.Use(middleware.AuthorizationMiddleware([]string{"SA", "AS", "D", "P"}))

	protected.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	protected.Get("/ws/:id", websocket.New(wsService.HandleWebSocketConnection))

	protected.Get("/my/profile", authHandler.GetMyProfile)
	protected.Put("/my/profile/update", authHandler.UpdateMyProfile)
	protected.Patch("/my/profile/picture", authHandler.ChangeProfilePicture)

	protected.Patch("/change-password", authHandler.ChangePassword)

	protected.Post("/logout", authHandler.Logout)
	protected.Post("/device-token", authHandler.AddDeviceToken)

	////////////////////////////////////// SUPER ADMIN //////////////////////////////////////
	
	protectedSuperAdmin := protected.Group("/superadmin")
	protectedSuperAdmin.Use(middleware.AuthorizationMiddleware([]string{"SA"}))
	protectedDriver := protected.Group("/driver")
	protectedDriver.Use(middleware.AuthorizationMiddleware([]string{"D"}))
	
	protectedParent := protected.Group("/parent")
	protectedParent.Use(middleware.AuthorizationMiddleware([]string{"P"}))

	// USER FOR SUPERADMIN
	protectedSuperAdmin.Get("/user/sa/all", userHandler.GetAllSuperAdmin)
	protectedSuperAdmin.Get("/user/as/all", userHandler.GetAllSchoolAdmin)
	protectedSuperAdmin.Get("/user/driver/all", userHandler.GetAllPermittedDriver)
	protectedSuperAdmin.Get("/user/sa/:id", userHandler.GetSpecSuperAdmin)
	protectedSuperAdmin.Get("/user/as/:id", userHandler.GetSpecSchoolAdmin)
	protectedSuperAdmin.Get("/user/driver/:id", userHandler.GetSpecPermittedDriver)
	protectedSuperAdmin.Post("/user/add", userHandler.AddUser)
	protectedSuperAdmin.Put("/user/update/:id", userHandler.UpdateUser)
	protectedSuperAdmin.Delete("/user/sa/delete/:id", userHandler.DeleteSuperAdmin)
	protectedSuperAdmin.Delete("/user/as/delete/:id", userHandler.DeleteSchoolAdmin)
	protectedSuperAdmin.Delete("/user/driver/delete/:id", userHandler.DeleteDriver)

	// SCHOOL FOR SUPERADMIN
	protectedSuperAdmin.Get("/school/all", schoolHandler.GetAllSchools)
	protectedSuperAdmin.Get("/school/:id", schoolHandler.GetSpecSchool)
	protectedSuperAdmin.Post("/school/add", schoolHandler.AddSchool)
	protectedSuperAdmin.Put("/school/update/:id", schoolHandler.UpdateSchool)
	protectedSuperAdmin.Delete("/school/delete/:id", schoolHandler.DeleteSchool)
	
	// VEHICLE FOR SUPERADMIN
	protectedSuperAdmin.Get("/vehicle/all", vehicleHandler.GetAllVehicles)
	protectedSuperAdmin.Get("/vehicle/free/all", vehicleHandler.GetAvailableVehicles)
	protectedSuperAdmin.Get("/vehicle/:id", vehicleHandler.GetSpecVehicle)
	protectedSuperAdmin.Post("/vehicle/add", vehicleHandler.AddVehicle)
	protectedSuperAdmin.Put("/vehicle/update/:id", vehicleHandler.UpdateVehicle)
	protectedSuperAdmin.Delete("/vehicle/delete/:id", vehicleHandler.DeleteVehicle)

	
	protectedSuperAdmin.Get("/shuttle/summary", shuttleHandler.GetShuttleSummary)
	protectedSuperAdmin.Get("/student/growth", studentHandler.GetStudentCountByMonth)


	////////////////////////////////////// SCHOOL ADMIN //////////////////////////////////////

	protectedSchoolAdmin := protected.Group("/school")
	protectedSchoolAdmin.Use(middleware.AuthorizationMiddleware([]string{"AS"}))
	protectedSchoolAdmin.Use(middleware.SchoolAdminMiddleware(userService))

	// STUDENT FOR SCHOOL ADMIN
	protectedSchoolAdmin.Get("/student/all", studentHandler.GetAllStudentWithParents)
	protectedSchoolAdmin.Get("/student/:id", studentHandler.GetSpecStudentWithParents)
	protectedSchoolAdmin.Post("/student/add", studentHandler.AddSchoolStudentWithParents)
	protectedSchoolAdmin.Put("/student/update/:id", studentHandler.UpdateSchoolStudentWithParents)
	protectedSchoolAdmin.Delete("/student/delete/:id", studentHandler.DeleteSchoolStudentWithParentsIfNeccessary)

	protectedSchoolAdmin.Get("/user/driver/all", userHandler.GetAllPermittedDriver)
	protectedSchoolAdmin.Get("/user/driver/:id", userHandler.GetSpecPermittedDriver)
	protectedSchoolAdmin.Post("/user/driver/add", userHandler.AddSchoolDriver)
	protectedSchoolAdmin.Put("/user/driver/update/:id", userHandler.UpdateSchoolDriver)
	protectedSchoolAdmin.Delete("/user/driver/delete/:id", userHandler.DeleteSchoolDriver)
	
	protectedSchoolAdmin.Get("/vehicle/all", vehicleHandler.GetAllVehiclesForPermittedSchool)
	protectedSchoolAdmin.Get("/vehicle/:id", vehicleHandler.GetSpecVehicleForPermittedSchool)
	protectedSchoolAdmin.Post("/vehicle/add", vehicleHandler.AddVehicleForPermittedSchool)
	protectedSchoolAdmin.Put("/vehicle/update/:id", vehicleHandler.UpdateVehicle)
	protectedSchoolAdmin.Delete("/vehicle/delete/:id", vehicleHandler.DeleteVehicle)

	// ROUTE FOR SCHOOL ADMIN
	protectedSchoolAdmin.Get("/route/all", routeHandler.GetAllRoutesByAS)
	protectedSchoolAdmin.Get("/route/:id", routeHandler.GetSpecRouteByAS)
	protectedSchoolAdmin.Post("/route/add", routeHandler.AddRoute)
	protectedSchoolAdmin.Put("/route/update/:id", routeHandler.UpdateRoute)
	protectedSchoolAdmin.Delete("/route/delete/:id", routeHandler.DeleteRoute)

	//ROUTE FOR DRIVER
	protectedDriver.Get("/route/all", routeHandler.GetAllRoutesByDriver)

	protectedParent.Get("/my/childern/track", shuttleHandler.GetShuttleTrackByParent) //buat menu track
	protectedParent.Get("/my/childern/all", childernHandler.GetAllChilderns) //buat menu apalah
	protectedParent.Get("/my/childern/shuttle/:id", shuttleHandler.GetSpecShuttle) //buat menu opo jeneng e lali😂 (spec shutle)
	protectedParent.Get("/my/childern/recap", shuttleHandler.GetAllShuttleByParent) //buat menu recap
	protectedParent.Get("/my/childern/:id", childernHandler.GetSpecChildern) //nih katanya butuh spec
	protectedParent.Put("/my/childern/update/:id", childernHandler.UpdateChildern) //menu update nih tampling
	protectedParent.Put("/my/childern/status/update/:id", childernHandler.UpdateChildernStatus) //menu update nih tampling

	protectedDriver.Get("/shuttle/all", shuttleHandler.GetAllShuttleByDriver)
	protectedDriver.Post("/shuttle/add", shuttleHandler.AddShuttle)
	protectedDriver.Get("/shuttle/:id", shuttleHandler.GetSpecShuttle)
	protectedDriver.Put("/shuttle/update/:id", shuttleHandler.EditShuttle) 
}