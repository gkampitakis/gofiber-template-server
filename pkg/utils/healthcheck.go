package utils

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gkampitakis/gofiber-template-server/pkg/configs"
	"github.com/gofiber/fiber/v2"
)

type HealthcheckMap map[string]func() bool

func RegisterHealthchecks(app *fiber.App, controls ...HealthcheckMap) {
	if len(controls) > 1 {
		log.Println("[Warning] only the 1st element is used")
	}

	app.Get("/health", registerHealthRoute(controls[0], configs.NewHealthcheckConfig()))
}

func registerHealthRoute(controls HealthcheckMap, config *configs.HealthcheckConfig) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		closeChannel := make(chan bool)
		wg := sync.WaitGroup{}
		controlsLength := len(controls)
		response, status := initializeResponse(controls, config)

		wg.Add(controlsLength)

		if controlsLength == 0 {
			response["status"] = "healthy"
		}

		for label, control := range controls {
			go func(label string, control func() bool) {
				defer wg.Done()
				res := control()
				if res {
					response[label] = "healthy"
					return
				}

				status = 500
				response[label] = "unhealthy"
			}(label, control)
		}

		go func() {
			defer close(closeChannel)
			wg.Wait()
		}()

		timeout(config, &status, closeChannel)

		return c.Status(status).JSON(response)
	}
}

func initializeResponse(checks HealthcheckMap, config *configs.HealthcheckConfig) (m map[string]string, status int) {
	m = make(map[string]string, len(checks))
	status = 200

	if !config.TimeoutEnabled {
		return m, status
	}

	for label := range checks {
		m[label] = fmt.Sprintf("Timeout after %d seconds", config.TimeoutPeriod)
	}

	return m, status
}

func timeout(config *configs.HealthcheckConfig, status *int, c <-chan bool) {
	if config.TimeoutEnabled {
		select {
		case <-time.After(time.Second * config.TimeoutPeriod):
			*status = 500
		case <-c:
		}

		return
	}
	<-c
}
