package handler

import "github.com/gofiber/fiber/v3"

type HttpServer struct {
	fiberApp *fiber.App
}

func NewHttpHandler() *HttpServer {
	app := fiber.New()
	return &HttpServer{
		fiberApp: app,
	}
}

func (hs *HttpServer) Start(port string) {
	go hs.fiberApp.Listen(port)
}

func (hs *HttpServer) App() *fiber.App {
	return hs.fiberApp
}
