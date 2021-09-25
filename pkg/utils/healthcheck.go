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

func RegisterHealthchecks(app *fiber.App, hc_config *configs.HealthcheckConfig, checks ...HealthcheckMap) {
	if len(checks) > 1 {
		log.Println("[Warning] only the 1st element is used")
	}

	var _checks HealthcheckMap

	if len(checks) == 0 {
		_checks = make(HealthcheckMap)
	} else {
		_checks = checks[0]
	}

	app.Get("/health", registerHealthRoute(hc_config, _checks))
}

// @Description Route reporting health of service
// @Summary Healthcheck route
// @Tags health
// @Accept text/plain
// @Product json/application
// @Success 200 {object} map[string]string "This can be dynamic and add more fields in checks"
// @Failure 500 {object} map[string]string "The route can return 500 in case of failed check,timeouts or panic"
// @Router /health [get]
func registerHealthRoute(config *configs.HealthcheckConfig, checks HealthcheckMap) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		checksLength := len(checks)

		// If we don't pass checks we prematurely respond as healthy, nothing to "check"
		if checksLength == 0 {
			return c.Status(http.StatusOK).JSON(map[string]string{"status": "healthy"})
		}

		results := initializeResults(checks, config)
		closeChannel := make(chan struct{})
		wg := sync.WaitGroup{}
		mutex := sync.Mutex{}

		wg.Add(checksLength)

		for label, control := range checks {
			go func(label string, control func() bool) {
				/**
				To future self, deferred function calles are push onto a stack. When function
				returns, its deferred called are executed in LIFO order.
				*/
				defer wg.Done()
				defer handlePanic(results, label)
				res := control()
				if res {
					results.Store(label, "healthy")
					return
				}

				results.Store(label, "unhealthy")
			}(label, control)
		}

		go func() {
			defer close(closeChannel)
			wg.Wait()
		}()

		timeout(config, closeChannel, &mutex)

		responseObject, status := getResponse(results)

		return c.Status(status).JSON(responseObject)
	}
}

func initializeResults(checks HealthcheckMap, config *configs.HealthcheckConfig) *sync.Map {
	var m sync.Map

	if !config.TimeoutEnabled {
		return &m
	}

	for label := range checks {
		m.Store(label, fmt.Sprintf("Timeout after %d seconds", config.TimeoutPeriod))
	}

	return &m
}

func getResponse(object *sync.Map) (map[string]string, int) {
	responseObject := make(map[string]string)
	status := http.StatusOK

	object.Range(func(key, value interface{}) bool {
		responseObject[key.(string)] = value.(string)
		if value.(string) != "healthy" {
			status = http.StatusInternalServerError
		}
		return true
	})

	return responseObject, status
}

func timeout(config *configs.HealthcheckConfig, c <-chan struct{}, l *sync.Mutex) {
	if config.TimeoutEnabled {
		select {
		case <-time.After(time.Second * config.TimeoutPeriod):
		case <-c:
		}

		return
	}
	<-c
}

func handlePanic(response *sync.Map, label string) {
	if e := recover(); e != nil {
		response.Store(label, fmt.Errorf("Paniced with error: %v", e).Error())
	}
}
