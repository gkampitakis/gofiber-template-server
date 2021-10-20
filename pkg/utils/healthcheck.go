package utils

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gkampitakis/gofiber-template-server/pkg/configs"
	"github.com/gofiber/fiber/v2"
)

var upTime = time.Now()

type HealthCheckResponse struct {
	Service      string            `json:"service"`
	Uptime       string            `json:"uptime"`
	Memory       MemoryMetrics     `json:"memory"`
	GoRoutines   int               `json:"go_routines"`
	HealthChecks map[string]string `json:"health_checks,omitempty"`
}

type MemoryMetrics struct {
	ResidentSetSize uint64 `json:"rss"`
	TotalAlloc      uint64 `json:"total_alloc"`
	HeapAlloc       uint64 `json:"heap_alloc"`
}

type C struct {
	msg   string
	label string
}

type HealthcheckMap map[string]func() bool

func RegisterHealthchecks(app *fiber.App, hc_config *configs.HealthcheckConfig, checks ...HealthcheckMap) {
	if len(checks) > 1 {
		log.Println("[Warning] only the 1st element is used")
	}

	if checks == nil {
		app.Get("/health", registerHealthRoute(hc_config, nil))
	} else {
		app.Get("/health", registerHealthRoute(hc_config, checks[0]))
	}
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
	return func(ctx *fiber.Ctx) error {
		// If we don't pass checks we prematurely respond as healthy, nothing to "check"
		if checks == nil {
			return ctx.Status(http.StatusOK).JSON(prepareResponse(config.Service, nil))
		}

		status := http.StatusOK
		response := prepareResponse(config.Service, map[string]string{})
		c := make(chan C)
		checksLength := len(checks)

		for label, control := range checks {
			go check(
				label,
				control,
				&status,
				c,
				config.TimeoutEnabled,
				config.TimeoutPeriod,
			)
		}

		for i := 0; i < checksLength; i++ {
			checkResponse := <-c
			response.HealthChecks[checkResponse.label] = checkResponse.msg
		}

		return ctx.Status(status).JSON(response)
	}
}

func prepareResponse(service string, healthChecks map[string]string) *HealthCheckResponse {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	return &HealthCheckResponse{
		Uptime:  time.Since(upTime).String(),
		Service: service,
		Memory: MemoryMetrics{
			ResidentSetSize: mem.HeapSys,
			TotalAlloc:      mem.TotalAlloc,
			HeapAlloc:       mem.HeapAlloc,
		},
		GoRoutines:   runtime.NumGoroutine(),
		HealthChecks: healthChecks,
	}
}

func check(
	label string,
	control func() bool,
	status *int, c chan<- C,
	timeoutEnabled bool,
	timeoutPeriod time.Duration,
) {
	/**
	To future self, this channel needs to be buffered.
	The reason is when the timeout is enabled and  the `res := control()` doesn't respond on time,
	later when it tries to push to `localC` there is no one to read the message so the goroutine
	is left there "stuck" and ending up leaking goroutines.
	*/
	localC := make(chan C, 1)
	// FIXME: go routine here is problematic
	go func() {
		/**
		To future self, deferred function calles are push onto a stack. When function
		returns, its deferred called are executed in LIFO order.
		*/
		defer func() {
			if e := recover(); e != nil {
				localC <- C{msg: fmt.Errorf("Paniced with error: %v", e).Error(), label: label}
				if *status == http.StatusOK {
					*status = http.StatusInternalServerError
				}
			}
		}()

		res := control()
		if res {
			localC <- C{msg: "healthy", label: label}
			return
		}

		if *status == http.StatusOK {
			*status = http.StatusInternalServerError
		}

		localC <- C{msg: "unhealthy", label: label}
	}()

	if timeoutEnabled {
		select {
		case tmp := <-localC:
			c <- tmp
		case <-time.After(time.Second * timeoutPeriod):
			c <- C{msg: fmt.Sprintf("Timeout after %d seconds", timeoutPeriod), label: label}
		}

		return
	}

	res := <-localC
	c <- res
}
