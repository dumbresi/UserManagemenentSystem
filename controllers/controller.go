package controllers

import (
	"slices"

	"github.com/gofiber/fiber/v2"
)

func CheckHealth(ctx *fiber.Ctx) error {

	ctx.Set("cache-control","no-cache")

	if ctx.Method() != fiber.MethodGet {
		ctx.Status(fiber.StatusMethodNotAllowed)
		return nil
	}

	

	if len(ctx.Body()) > 0 {
		ctx.Status(fiber.StatusBadRequest)
		return nil
	}

	allowedHeaders := []string{"User-Agent", "Accept", "Host", "Connection", "Content-Length", "Date", "Content-Type"}

	ctx.Request().Header.VisitAll(func(key []byte, value []byte) {
		header := string(key)

		if !slices.Contains(allowedHeaders, header) {
			ctx.Status(fiber.StatusBadRequest)
		}
	})

	if len(ctx.Queries()) > 0 {
		ctx.Status(fiber.StatusBadRequest)
		return nil
	}
	ctx.Status(fiber.StatusOK)
	return nil
}

func ErrorHealthCheck(ctx *fiber.Ctx)error {
	ctx.Set("cache-control","no-cache")
	ctx.Status(fiber.StatusNotFound)
	return nil
}