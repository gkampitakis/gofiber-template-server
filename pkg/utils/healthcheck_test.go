package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/gkampitakis/gofiber-template-server/pkg/configs"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func healthRequest(t *testing.T, app *fiber.App, timeout int) ([]byte, int) {
	req, err := http.NewRequest(
		"GET",
		"/health",
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	res, err := app.Test(req, timeout)

	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Fatal(err)
	}

	return body, res.StatusCode
}

func TestHealthcheckRoute(t *testing.T) {
	t.Run("should return just a status", func(t *testing.T) {
		app := fiber.New()
		RegisterHealthchecks(app, configs.NewHealthcheckConfig())
		response := HealthCheckResponse{}

		body, statusCode := healthRequest(t, app, 100)
		err := json.Unmarshal(body, &response)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 200, statusCode)

		assert.Equal(t, "gofiber-template", response.Service)
		assert.Equal(t, 0, len(response.HealthChecks))
	})

	t.Run("should list all checks healthy", func(t *testing.T) {
		app := fiber.New()
		checks := HealthcheckMap{
			"check1": func() bool { return true },
			"check2": func() bool { return true },
			"check3": func() bool { return true },
		}

		RegisterHealthchecks(app, configs.NewHealthcheckConfig(), checks)
		response := HealthCheckResponse{}

		body, statusCode := healthRequest(t, app, 100)
		err := json.Unmarshal(body, &response)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 200, statusCode)
		assert.Equal(t, 3, len(response.HealthChecks))
		for label, value := range response.HealthChecks {
			assert.Equal(t, "healthy", value)
			_, exists := checks[label]
			assert.True(t, exists)
		}
	})

	t.Run("should handle panic inside checks", func(t *testing.T) {
		app := fiber.New()
		checks := HealthcheckMap{
			"check1": func() bool { return true },
			"check2": func() bool {
				panic("boo 👻")
			},
			"check3": func() bool { return true },
		}

		RegisterHealthchecks(app, configs.NewHealthcheckConfig(), checks)
		response := HealthCheckResponse{}

		body, statusCode := healthRequest(t, app, 100)
		err := json.Unmarshal(body, &response)

		if err != nil {
			t.Fatal(err)
		}

		checksResults := response.HealthChecks

		assert.Equal(t, 500, statusCode)
		assert.Equal(t, 3, len(checksResults))
		assert.Equal(t, "healthy", checksResults["check1"])
		assert.Equal(t, "Paniced with error: boo 👻", checksResults["check2"])
		assert.Equal(t, "healthy", checksResults["check3"])
	})

	t.Run("should report slow checks as timedout and 500 statusCode", func(t *testing.T) {
		app := fiber.New()
		checks := HealthcheckMap{
			"check1": func() bool { return false },
			"check2": func() bool {
				time.Sleep(2 * time.Second)
				return true
			},
			"check3": func() bool { return true },
		}

		RegisterHealthchecks(app, &configs.HealthcheckConfig{
			TimeoutEnabled: true,
			TimeoutPeriod:  1,
		}, checks)
		response := HealthCheckResponse{}

		body, statusCode := healthRequest(t, app, 3000)
		err := json.Unmarshal(body, &response)
		if err != nil {
			t.Fatal(err)
		}

		checksResults := response.HealthChecks

		log.Println(response.HealthChecks)

		assert.Equal(t, 500, statusCode)
		assert.Equal(t, 3, len(checksResults))
		assert.Equal(t, "unhealthy", checksResults["check1"])
		assert.Equal(t, "Timeout after 1 seconds", checksResults["check2"])
		assert.Equal(t, "healthy", checksResults["check3"])
	})

	t.Run("should not report slow checks as timedout", func(t *testing.T) {
		app := fiber.New()
		checks := HealthcheckMap{
			"check1": func() bool { return true },
			"check2": func() bool {
				time.Sleep(2 * time.Second)
				return true
			},
			"check3": func() bool { return true },
		}

		RegisterHealthchecks(app, &configs.HealthcheckConfig{
			TimeoutEnabled: false,
			TimeoutPeriod:  1,
		}, checks)
		response := HealthCheckResponse{}

		body, statusCode := healthRequest(t, app, 3000)
		err := json.Unmarshal(body, &response)
		if err != nil {
			t.Fatal(err)
		}

		checksResults := response.HealthChecks

		assert.Equal(t, 200, statusCode)
		assert.Equal(t, 3, len(checksResults))
		assert.Equal(t, "healthy", checksResults["check1"])
		assert.Equal(t, "healthy", checksResults["check2"])
		assert.Equal(t, "healthy", checksResults["check3"])
	})
}
