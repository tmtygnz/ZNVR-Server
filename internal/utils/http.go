package utils

import (
	"sync"

	"github.com/gofiber/fiber/v3"
)

type ResponseBase struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

var responsePool = sync.Pool{
	New: func() any {
		return new(ResponseBase)
	},
}

func CreateMessage(c fiber.Ctx, statusCode int, message string, extraData any) error {
	res := responsePool.Get().(*ResponseBase)

	res.Message = message
	res.Data = extraData

	err := c.Status(statusCode).JSON(res)

	res.Message = ""
	res.Data = nil
	responsePool.Put(res)

	return err
}
