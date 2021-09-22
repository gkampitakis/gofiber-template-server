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

// FIXME: there are errors with -race detection in tests
// go test -run HealthcheckRoute ./... -v -count=1 -race -count=1 for not caching results

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
		responseObject := make(map[string]string)

		body, statusCode := healthRequest(t, app, 100)
		err := json.Unmarshal(body, &responseObject)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 200, statusCode)
		assert.Equal(t, 1, len(responseObject))
		assert.Equal(t, "healthy", responseObject["status"])
	})

	t.Run("should enlist all checks healthy", func(t *testing.T) {
		app := fiber.New()
		checks := HealthcheckMap{
			"check1": func() bool { return true },
			"check2": func() bool { return true },
			"check3": func() bool { return true },
		}

		RegisterHealthchecks(app, configs.NewHealthcheckConfig(), checks)
		responseObject := make(map[string]string)

		body, statusCode := healthRequest(t, app, 100)
		err := json.Unmarshal(body, &responseObject)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 200, statusCode)
		assert.Equal(t, 3, len(responseObject))
		for label, value := range responseObject {
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
		responseObject := make(map[string]string)

		body, statusCode := healthRequest(t, app, 100)
		err := json.Unmarshal(body, &responseObject)

		if err != nil {
			t.Fatal(err)
		}

		log.Println(string(body), statusCode)

		assert.Equal(t, 500, statusCode)
		assert.Equal(t, 3, len(responseObject))
		assert.Equal(t, "healthy", responseObject["check1"])
		assert.Equal(t, "Paniced with error: boo 👻", responseObject["check2"])
		assert.Equal(t, "healthy", responseObject["check3"])
	})

	t.Run("should return timed out checks and 500 statusCode", func(t *testing.T) {
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
		responseObject := make(map[string]string)

		body, statusCode := healthRequest(t, app, 3000)
		err := json.Unmarshal(body, &responseObject)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 500, statusCode)
		assert.Equal(t, 3, len(responseObject))
		assert.Equal(t, "unhealthy", responseObject["check1"])
		assert.Equal(t, "Timeout after 1 seconds", responseObject["check2"])
		assert.Equal(t, "healthy", responseObject["check3"])
	})

	t.Run("should ignore not time out healthchecks", func(t *testing.T) {
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
			TimeoutEnabled: false,
			TimeoutPeriod:  1,
		}, checks)
		responseObject := make(map[string]string)

		body, statusCode := healthRequest(t, app, 3000)
		err := json.Unmarshal(body, &responseObject)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 500, statusCode)
		assert.Equal(t, 3, len(responseObject))
		assert.Equal(t, "unhealthy", responseObject["check1"])
		assert.Equal(t, "healthy", responseObject["check2"])
		assert.Equal(t, "healthy", responseObject["check3"])
	})
}
