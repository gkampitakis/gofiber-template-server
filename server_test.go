package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	t.Run("setup server", func(t *testing.T) {
		t.Run("Development", func(t *testing.T) {
			t.Run("profilling route should be enabled", func(t *testing.T) {
				app := SetupServer(true)

				req, err := http.NewRequest(
					"GET",
					"/debug/pprof/",
					nil,
				)
				if err != nil {
					t.Fatal(err)
				}

				res, err := app.Test(req, 100)
				if err != nil {
					t.Error(err)
				}

				assert.NotEqual(t, 404, res.StatusCode)
				assert.Equal(t, 200, res.StatusCode)
				assert.Equal(t, res.Header["Content-Type"], []string{"text/html; charset=utf-8"})
			})

			t.Run("swagger route should be enabled", func(t *testing.T) {
				app := SetupServer(true)

				req, err := http.NewRequest(
					"GET",
					"/swagger/index.html",
					nil,
				)
				if err != nil {
					t.Fatal(err)
				}

				res, err := app.Test(req, 100)
				if err != nil {
					t.Error(err)
				}

				assert.NotEqual(t, 404, res.StatusCode)
				assert.Equal(t, 200, res.StatusCode)
				assert.Equal(t, res.Header["Content-Type"], []string{"text/html"})
			})
		})

		t.Run("Production", func(t *testing.T) {
			t.Run("profilling route should be disabled", func(t *testing.T) {
				app := SetupServer(false)

				req, err := http.NewRequest(
					"GET",
					"/debug/pprof/",
					nil,
				)
				if err != nil {
					t.Fatal(err)
				}

				res, err := app.Test(req, 100)
				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, 404, res.StatusCode)
			})

			t.Run("swagger route should be disabled", func(t *testing.T) {
				app := SetupServer(false)

				req, err := http.NewRequest(
					"GET",
					"/swagger/index.html",
					nil,
				)
				if err != nil {
					t.Fatal(err)
				}

				res, err := app.Test(req, 100)
				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, 404, res.StatusCode)
			})
		})
	})

	t.Run("routes", func(t *testing.T) {
		t.Run("[/] should greet", func(t *testing.T) {
			app := SetupServer(false)

			req, err := http.NewRequest(
				"GET",
				"/",
				nil,
			)
			if err != nil {
				t.Fatal(err)
			}

			res, err := app.Test(req, 100)
			if err != nil {
				t.Error(err)
			}

			body, _ := ioutil.ReadAll(res.Body)
			assert.Equal(t, 200, res.StatusCode)
			assert.Equal(t, []byte("Hello World from template"), body)
		})

		t.Run("[/hello/:name] should greet with name", func(t *testing.T) {
			app := SetupServer(false)

			req, err := http.NewRequest(
				"GET",
				"/hello/test",
				nil,
			)
			if err != nil {
				t.Fatal(err)
			}

			res, err := app.Test(req, 100)
			if err != nil {
				t.Error(err)
			}

			body, _ := ioutil.ReadAll(res.Body)
			assert.Equal(t, 200, res.StatusCode)
			assert.Equal(t, []byte("Hello test 👋"), body)
		})

		t.Run("[/health] should respond healthy", func(t *testing.T) {
			bodyMap := make(map[string]string)
			app := SetupServer(false)

			req, err := http.NewRequest(
				"GET",
				"/health",
				nil,
			)
			if err != nil {
				t.Fatal(err)
			}

			res, err := app.Test(req, 100)
			if err != nil {
				t.Error(err)
			}

			body, _ := ioutil.ReadAll(res.Body)
			err = json.Unmarshal(body, &bodyMap)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, 200, res.StatusCode)
			assert.Equal(t, "healthy", bodyMap["status"])
		})
	})
}
