package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"pendaftaran/config"
	"pendaftaran/controllers"
	"pendaftaran/database"
	"pendaftaran/middleware"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	database.InitDB()
	defer database.CloseDB()

	app.Post("/login", controllers.Login)
	app.Post("/logout", controllers.Logout)

	app.Static("/uploads", "./uploads")
	app.Use(middleware.NewAuthMiddleware(config.Secret))

	app.Get("/", func(c *fiber.Ctx) error { return c.SendStatus(200) })

	app.Get("/user", controllers.GetUser)
	app.Get("/user/:id", controllers.GetUserByID)
	app.Post("/user", controllers.CreateUser)
	app.Put("/user/:id", controllers.UpdateUser)
	app.Delete("/user/:id", controllers.DeleteUser)

	app.Get("/jurusan", controllers.GetJurusan)
	app.Get("/jurusan/:id", controllers.GetJurusanByID)
	app.Post("/jurusan", controllers.CreateJurusan)
	app.Put("/jurusan/:id", controllers.UpdateJurusan)
	app.Delete("/jurusan/:id", controllers.DeleteJurusan)

	app.Get("/smp", controllers.GetSmp)
	app.Get("/smp/:id", controllers.GetSmpByID)
	app.Post("/smp", controllers.CreateSmp)
	app.Put("/smp/:id", controllers.UpdateSmp)
	app.Delete("/smp/:id", controllers.DeleteSmp)

	app.Get("/tahun", controllers.GetTahun)
	app.Get("/tahun/:id", controllers.GetTahunByID)
	app.Post("/tahun", controllers.CreateTahun)
	app.Put("/tahun/:id", controllers.UpdateTahun)
	app.Delete("/tahun/:id", controllers.DeleteTahun)

	app.Get("/guru", controllers.GetGuru)
	app.Get("/guru/:id", controllers.GetGuruByID)
	app.Post("/guru", controllers.CreateGuru)
	app.Put("/guru/:id", controllers.UpdateGuru)
	app.Delete("/guru/:id", controllers.DeleteGuru)

	app.Get("/siswa", controllers.GetSiswa)
	app.Get("/siswa/:id", controllers.GetSiswaByID)
	app.Post("/siswa", controllers.CreateSiswa)
	app.Put("/siswa/:id", controllers.UpdateSiswa)
	app.Delete("/siswa/:id", controllers.DeleteSiswa)

	app.Get("/pendaftaran", controllers.GetPendaftaran)
	app.Get("/pendaftaran/:id", controllers.GetPendaftaranByID)
	app.Post("/pendaftaran", controllers.CreatePendaftaran)
	app.Put("/pendaftaran/:id", controllers.UpdatePendaftaran)
	app.Delete("/pendaftaran/:id", controllers.DeletePendaftaran)

	app.Listen(":3030")
}
