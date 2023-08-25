package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/zerolog/log"
)

var ErrorHandling = func(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(err)
}

func Init() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandling,
	})

	app.Use(requestid.New())

	return app
}

func Listen(app *fiber.App, address string) {
	err := app.Listen(address)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start the server")
	}
}

func Shutdown(app *fiber.App, cleanup func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Print("Shutting down the server")
	_ = app.Shutdown()
	if cleanup != nil {
		log.Print("Cleanning up tasks")
		cleanup()
	}
	log.Print("Successfully shutdown the server")
}
