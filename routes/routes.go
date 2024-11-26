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

	// FOR AUTHENTICATED
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())

	protected.Get("/my/profile", controllers.GetMyProfile)

	protected.Post("/logout", controllers.Logout)

	// FOR SUPERADMIN
	protectedSA := r.Group("/superadmin")
	protectedSA.Use(middleware.SuperAdminMiddleware())

	protectedSA.Get("/user/all", controllers.GetAllUser)

	protectedSA.Get("/user/:id", controllers.GetSpecUser)
	protectedSA.Post("/user/add", controllers.AddUser)
	protectedSA.Put("/user/update/:id", controllers.UpdateUser)
	protectedSA.Delete("/user/delete/:id", controllers.DeleteUser)
	
	protectedSA.Get("/school/all", controllers.GetAllSchools)
	protectedSA.Get("/school/:id", controllers.GetSpecSchool)
	protectedSA.Post("/school/add", controllers.AddSchool)
	protectedSA.Put("/school/update/:id", controllers.UpdateSchool)
	protectedSA.Delete("/school/delete/:id", controllers.DeleteSchool)

	protectedSA.Get("/vehicle/all", controllers.GetAllVehicles)
	protectedSA.Get("/vehicle/:id", controllers.GetSpecVehicle)
	protectedSA.Post("/vehicle/add", controllers.AddVehicle)
	protectedSA.Put("/vehicle/update/:id", controllers.UpdateVehicle)
	protectedSA.Delete("/vehicle/delete/:id", controllers.DeleteVehicle)

	// FOR SCHOOL ADMIN
	protectedSchAdmin := r.Group("/admin")
	protectedSchAdmin.Use(middleware.SchoolAdminMiddleware())

	protectedSchAdmin.Get("/school/student/all", controllers.GetAllStudentWithParents)
	protectedSchAdmin.Post("/school/student/add", controllers.AddSchoolStudentWithParents)
	protectedSchAdmin.Put("/school/student/update/:id", controllers.UpdateSchoolStudentWithParents)
	protectedSchAdmin.Delete("/school/student/delete/:id", controllers.DeleteSchoolStudentWithParents)

	protectedSchAdmin.Post("/school/route/add", controllers.AddRoadRoute)
}

// tes jawa