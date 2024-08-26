package main

import (
	"database/sql"
	"log"

	"unap-auth/adapter/http"
	"unap-auth/config"

	"unap-auth/infrastructure/mysql"
	"unap-auth/usecase"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Cargar la configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Establecer la conexión a la base de datos
	db, err := sql.Open("mysql", cfg.DBDSN)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	defer db.Close()

	// Verificar la conexión
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	// Configurar la clave secreta para JWT
	jwtSecret := cfg.JWTSecret

	// Inicializar el repositorio y usecase
	userRepo := &mysql.MySQLUserAccountRepository{DB: db}
	authUsecase := usecase.NewAuthUsecase(userRepo, jwtSecret)

	// Inicializar el manejador HTTP
	authHandler := http.NewAuthHandler(authUsecase)

	// Configurar el servidor y las rutas
	app := fiber.New()
	http.SetupRoutes(app, authHandler)

	// Iniciar el servidor
	if err := app.Listen(":3500"); err != nil {
		log.Fatal(err)
	}
}
