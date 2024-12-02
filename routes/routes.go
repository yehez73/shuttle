package routes

import (
	"shuttle/controllers"
	"shuttle/middleware"

	"github.com/gofiber/fiber/v2"
)

func Route(r *fiber.App) {
	// FOR PUBLIC
	r.Post("/login", controllers.Login)
	r.Post("/refresh-token", controllers.RefreshToken)
	r.Static("/assets", "./assets")

	// FOR AUTHENTICATED
	protected := r.Group("/api")
	protected.Use(middleware.AuthenticationMiddleware())
	protected.Use(middleware.AuthorizationMiddleware([]string{"SA", "AS", "D", "P"}))

	protected.Get("/my/profile", controllers.GetMyProfile)
	protected.Post("/logout", controllers.Logout)

	protectedSuperAdmin := protected.Group("/superadmin")
	protectedSuperAdmin.Use(middleware.AuthorizationMiddleware([]string{"SA"}))

	protectedSchoolAdmin := protected.Group("/school")
	protectedSchoolAdmin.Use(middleware.AuthorizationMiddleware([]string{"AS"}))
	protectedSchoolAdmin.Use(middleware.SchoolAdminMiddleware())

	// USER
	protectedSuperAdmin.Get("/user/sa/all", controllers.GetAllSuperAdmin)
	protectedSuperAdmin.Get("/user/sa/:id", controllers.GetSpecSuperAdmin)

	protectedSuperAdmin.Get("/user/as/all", controllers.GetAllSchoolAdmin)
	protectedSuperAdmin.Get("/user/as/:id", controllers.GetSpecSchoolAdmin)

	protectedSuperAdmin.Get("/user/driver/all", controllers.GetAllPermittedDriver)
	protectedSuperAdmin.Get("/user/driver/:id", controllers.GetSpecPermittedDriver)
	
	protectedSchoolAdmin.Get("/user/driver/all", controllers.GetAllPermittedDriver)
	protectedSchoolAdmin.Get("/user/driver/:id", controllers.GetSpecPermittedDriver)

	protectedSuperAdmin.Post("/user/add", controllers.AddUser)
	protectedSchoolAdmin.Post("/user/driver/add", controllers.AddSchoolDriver)
	
	protectedSuperAdmin.Put("/user/update/:id", controllers.UpdateUser)
	protectedSchoolAdmin.Put("/user/driver/update/:id", controllers.UpdateSchoolDriver)

	protectedSuperAdmin.Delete("/user/delete/:id", controllers.DeleteUser)
	protectedSchoolAdmin.Delete("/user/driver/delete/:id", controllers.DeleteSchoolDriver)

	// SCHOOL
	protectedSuperAdmin.Get("/school/all", controllers.GetAllSchools)
	protectedSuperAdmin.Get("/school/:id", controllers.GetSpecSchool)
	protectedSuperAdmin.Post("/school/add", controllers.AddSchool)
	protectedSuperAdmin.Put("/school/update/:id", controllers.UpdateSchool)
	protectedSuperAdmin.Delete("/school/delete/:id", controllers.DeleteSchool)

	protectedSuperAdmin.Get("/vehicle/all", controllers.GetAllVehicles)
	protectedSuperAdmin.Get("/vehicle/:id", controllers.GetSpecVehicle)
	protectedSuperAdmin.Post("/vehicle/add", controllers.AddVehicle)
	protectedSuperAdmin.Put("/vehicle/update/:id", controllers.UpdateVehicle)
	protectedSuperAdmin.Delete("/vehicle/delete/:id", controllers.DeleteVehicle)

	protectedSchoolAdmin.Get("/student/all", controllers.GetAllStudentWithParents)
	protectedSchoolAdmin.Post("/student/add", controllers.AddSchoolStudentWithParents)
	protectedSchoolAdmin.Put("/student/update/:id", controllers.UpdateSchoolStudentWithParents)
	protectedSchoolAdmin.Delete("/student/delete/:id", controllers.DeleteSchoolStudentWithParents)

	protectedSchoolAdmin.Get("/route/all", controllers.GetAllRoutes)
	protectedSchoolAdmin.Get("/route/:id", controllers.GetSpecRoute)
	protectedSchoolAdmin.Post("/route/add", controllers.AddRoute)
}