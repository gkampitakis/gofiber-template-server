module github.com/gkampitakis/gofiber-template-server

go 1.17

require (
	github.com/gofiber/fiber/v2 v2.21.0
	github.com/joho/godotenv v1.4.0
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.19.1
)

require (
	github.com/arsmn/fiber-swagger/v2 v2.20.0
	github.com/gkampitakis/fiber-modules v1.0.0-beta
	github.com/stretchr/testify v1.7.0
	github.com/swaggo/swag v1.7.4
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/andybalholm/brotli v1.0.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/spec v0.20.3 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/qri-io/jsonpointer v0.1.1 // indirect
	github.com/qri-io/jsonschema v0.2.1 // indirect
	github.com/swaggo/files v0.0.0-20210815190702-a29dd2bc99b2 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.31.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d // indirect
	golang.org/x/sys v0.0.0-20211105183446-c75c47738b0c // indirect
	golang.org/x/text v0.3.6 // indirect
	golang.org/x/tools v0.1.5 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/gkampitakis/fiber-modules v1.0.0-beta => ../gofiber-custom-healthcheck
