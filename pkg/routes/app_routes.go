package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gkampitakis/gofiber-template-server/app/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/xeipuuv/gojsonschema"
)

func loadValidatorSchema(path string) (*gojsonschema.Schema, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// NOTE: here we should cache the schema
	loader := gojsonschema.NewReferenceLoader("file://" + filepath.Join(wd, path))
	return gojsonschema.NewSchema(loader)
}

func AppRoutes(app *fiber.App) {
	route := app.Group("")

	route.Get("/", controllers.Default)
	route.Get("/hello/:name", controllers.Hello)

	route.Post("/test/validation", myMiddleware(ValidatorOptions{
		body: "./pkg/routes/test.json", // TODO: one file with all objects
	}), func(c *fiber.Ctx) error {
		body := string(c.Body())

		fmt.Println(body)

		c.Response().SetStatusCode(200)

		return c.Send([]byte("Hello World"))
	})
}

type ValidatorOptions struct { // NOTE: add support for json literal or file path
	headers string
	body    string
	params  string
	query   string
}

// NOTE: expose errors ??
// Logger ?
// 500 msg
// bad request message

func myMiddleware(v ValidatorOptions) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		if v.body != "" {
			s, err := loadValidatorSchema(v.body)
			if err != nil {
				log.Println(err)
				return c.Status(http.StatusInternalServerError).Send([]byte("Internal server error"))
			}
			bodySchema := gojsonschema.NewStringLoader(string(c.Body()))

			err, validationErrors := validateObject(s, bodySchema)
			if err != nil {
				log.Println(err)
				return c.Status(http.StatusInternalServerError).Send([]byte("Internal server error"))
			}

			if validationErrors != nil && len(validationErrors) > 0 {
				log.Println(validationErrors)

				return c.Status(http.StatusBadRequest).Send([]byte("Invalid request"))
			}
		}

		if v.params != "" {
		}

		if v.query != "" {
			q, err := loadValidatorSchema(v.query)
			if err != nil {
				log.Println(err)
				return c.Status(http.StatusInternalServerError).Send([]byte("Internal server error"))
			}

			queryJSON, err := queryToJSON(c.Context().QueryArgs().String())
			if err != nil {
				log.Println(err)
				return c.Status(http.StatusInternalServerError).Send([]byte("Internal server error"))
			}
			querySchema := gojsonschema.NewStringLoader(string(queryJSON))

			err, validationErrors := validateObject(q, querySchema)
			if err != nil {
				log.Println(err)
				return c.Status(http.StatusInternalServerError).Send([]byte("Internal server error"))
			}

			if validationErrors != nil && len(validationErrors) > 0 {
				log.Println(validationErrors)

				return c.Status(http.StatusBadRequest).Send([]byte("Invalid request"))
			}
		}

		if v.headers != "" {
		}

		return c.Next()
	}
}

func validateObject(s *gojsonschema.Schema, o gojsonschema.JSONLoader) (error, []gojsonschema.ResultError) {
	res, err := s.Validate(o)
	if err != nil {
		return err, nil
	}

	if res.Valid() {
		return nil, nil
	}

	return nil, res.Errors()
}

func queryToJSON(query string) ([]byte, error) {
	values, err := url.ParseQuery(query)
	if err != nil {
		return []byte(""), err
	}

	return json.Marshal(values)
}

// NOTE: we will need support for caching the schemas
// params validation, query validation, body validation, headers
