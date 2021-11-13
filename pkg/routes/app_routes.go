package routes

import (
	"context"
	"fmt"
	"regexp"

	jptr "github.com/qri-io/jsonpointer"
	"github.com/qri-io/jsonschema"

	"github.com/gkampitakis/fiber-modules/bodyvalidator"
	"github.com/gkampitakis/gofiber-template-server/app/controllers"
	"github.com/gofiber/fiber/v2"
)

var emailRgx, _ = regexp.Compile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

type IsEmail bool

func newIsEmail() jsonschema.Keyword {
	return new(IsEmail)
}

func (e *IsEmail) Register(uri string, registry *jsonschema.SchemaRegistry) {}

func (e *IsEmail) Resolve(pointer jptr.Pointer, uri string) *jsonschema.Schema {
	return nil
}

// ValidateKeyword implements jsonschema.Keyword
func (e IsEmail) ValidateKeyword(ctx context.Context, currentState *jsonschema.ValidationState, data interface{}) {
	if str, ok := data.(string); ok {
		if !emailRgx.Match([]byte(str)) {
			currentState.AddError(data, fmt.Sprintf("should be comform to email. plz make '%s' == type email. plz", str))
		}
	}
}

var validator = bodyvalidator.New(
	// bodyvalidator.ExposeErrors(),
	// bodyvalidator.SetResponse(func(ke []jsonschema.KeyError) interface{} {
	// 	return "you fucked up"
	// }),
	bodyvalidator.RegisterKeywords(bodyvalidator.Keywords{"email": newIsEmail}),
)

var testMiddleware = func() fiber.Handler {
	return validator(bodyvalidator.Config{
		// SchemaLiteral: t,
		SchemaPath: "pkg/routes/test.json",
	})
}

func AppRoutes(app *fiber.App) {
	route := app.Group("")

	route.Get("/", controllers.Default)
	route.Get("/hello/:name", controllers.Hello)

	route.Post("/test/validation",
		testMiddleware(),
		func(c *fiber.Ctx) error {
			return c.Status(200).JSON("Hello World")
		})
}

// validator = routevalidator.New(routevalidator.ExposeErrors(),routevalidator.RegisterKeyWords([]{}))
// validator({
// 	// SchemaPath: "pkg/routes/bigJSONValidator.json",
// 	SchemaLiteral: t,
// })
