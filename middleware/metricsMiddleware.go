package middleware

import (
	"time"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/stats"
	"github.com/gofiber/fiber/v2"
)

func MetricsMiddleware(c *fiber.Ctx) error {
    
    start := time.Now()

    err := c.Next() 

    endpoint := c.Path()

    stats.CountAPICall(endpoint)

    stats.TimeAPICall(endpoint, start)

    return err
}