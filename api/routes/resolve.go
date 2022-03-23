package routes

import (
	"log"
	"shorten-url/database"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")
	log.Println("ResolveURL: ", url)

	r := database.CreateClient(0)
	defer r.Close()

	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "short not found in the database",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot connect to DB",
		})
	}

	r1nr := database.CreateClient(1)
	log.Println("R1NR: ", r1nr)
	defer r1nr.Close()

	_ = r1nr.Incr(database.Ctx, "counter")

	return c.Redirect(value, 301)
}
