package utils

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gkampitakis/gofiber-template-server/pkg/configs"
	"github.com/gofiber/fiber/v2"
)

type HealthcheckMap map[string]func() bool

func RegisterHealthchecks(app *fiber.App, checks ...HealthcheckMap) {
	if len(checks) > 1 {
		log.Println("[Warning] only the 1st element is used")
	}

	var _checks HealthcheckMap

	if len(checks) == 0 {
		_checks = make(HealthcheckMap)
	} else {
		_checks = checks[0]
	}

	app.Get("/health", registerHealthRoute(configs.NewHealthcheckConfig(), _checks))
}

// FIXME: register this route to swagger
func registerHealthRoute(config *configs.HealthcheckConfig, checks HealthcheckMap) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		checksLength := len(checks)
		response, status := initializeResponse(checks, config)

		// If we don't pass checks we prematurely respond as healthy, nothing to "check"
		if checksLength == 0 {
			response["status"] = "healthy"
			return c.Status(status).JSON(response)
		}

		closeChannel := make(chan bool)
		wg := sync.WaitGroup{}

		wg.Add(checksLength)

		for label, control := range checks {
			go func(label string, control func() bool) {
				defer wg.Done()
				res := control()
				if res {
					response[label] = "healthy"
					return
				}

				status = http.StatusInternalServerError
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
	status = http.StatusOK

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
			*status = http.StatusInternalServerError
		case <-c:
		}

		return
	}
	<-c
}
