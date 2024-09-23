package controllers

import "github.com/gofiber/fiber/v2"

func CheckHealth(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}