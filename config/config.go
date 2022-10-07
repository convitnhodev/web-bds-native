package config

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	*viper.Viper

	errorHandler fiber.ErrorHandler
	fiber        *fiber.Config
}

var defaultErrorHandler = func(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Set error message
	message := err.Error()

	// Check if it's a fiber.Error type
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// TODO: Check return type for the client, JSON, HTML, YAML or any other (API vs web)

	// Return HTTP response
	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	c.Status(code)

	// Render default error view
	err = c.Render("errors/"+strconv.Itoa(code), fiber.Map{"message": message})
	if err != nil {
		return c.SendString(message)
	}
	return err
}

func New() *Config {
	config := &Config{
		Viper: viper.New(),
	}

	// Set default configurations
	config.setDefaults()

	// Select the .env file
	config.SetConfigName(".env")
	config.SetConfigType("dotenv")
	config.AddConfigPath(".")

	// Automatically refresh environment variables
	config.AutomaticEnv()

	// Read configuration
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("failed to read configuration:", err.Error())
			os.Exit(1)
		}
	}

	config.SetErrorHandler(defaultErrorHandler)

	// TODO: Logger (Maybe a different zap object)

	// TODO: Add APP_KEY generation

	// TODO: Write changes to configuration file

	// Set Fiber configurations
	config.setFiberConfig()

	return config
}

func (config *Config) SetErrorHandler(errorHandler fiber.ErrorHandler) {
	config.errorHandler = errorHandler
}

func (config *Config) setDefaults() {
	// Set default App configuration
	config.SetDefault("APP_ADDR", ":3000")
	config.SetDefault("APP_ENV", "local")
	// Set default database configuration
	config.SetDefault("DB_URI", "postgresql://postgres:postgres@localhost/postgres")

	// Set default Fiber configuration
	config.SetDefault("FIBER_PREFORK", false)
	config.SetDefault("FIBER_SERVERHEADER", "")
	config.SetDefault("FIBER_STRICTROUTING", false)
	config.SetDefault("FIBER_CASESENSITIVE", false)
	config.SetDefault("FIBER_IMMUTABLE", false)
	config.SetDefault("FIBER_UNESCAPEPATH", false)
	config.SetDefault("FIBER_ETAG", false)
	config.SetDefault("FIBER_BODYLIMIT", 4194304)
	config.SetDefault("FIBER_CONCURRENCY", 262144)
	config.SetDefault("FIBER_VIEWS", "html")
	config.SetDefault("FIBER_VIEWS_DIRECTORY", "resources/views")
	config.SetDefault("FIBER_VIEWS_RELOAD", false)
	config.SetDefault("FIBER_VIEWS_DEBUG", false)
	config.SetDefault("FIBER_VIEWS_LAYOUT", "embed")
	config.SetDefault("FIBER_VIEWS_DELIMS_L", "{{")
	config.SetDefault("FIBER_VIEWS_DELIMS_R", "}}")
	config.SetDefault("FIBER_READTIMEOUT", 0)
	config.SetDefault("FIBER_WRITETIMEOUT", 0)
	config.SetDefault("FIBER_IDLETIMEOUT", 0)
	config.SetDefault("FIBER_READBUFFERSIZE", 4096)
	config.SetDefault("FIBER_WRITEBUFFERSIZE", 4096)
	config.SetDefault("FIBER_COMPRESSEDFILESUFFIX", ".fiber.gz")
	config.SetDefault("FIBER_PROXYHEADER", "")
	config.SetDefault("FIBER_GETONLY", false)
	config.SetDefault("FIBER_DISABLEKEEPALIVE", false)
	config.SetDefault("FIBER_DISABLEDEFAULTDATE", false)
	config.SetDefault("FIBER_DISABLEDEFAULTCONTENTTYPE", false)
	config.SetDefault("FIBER_DISABLEHEADERNORMALIZING", false)
	config.SetDefault("FIBER_DISABLESTARTUPMESSAGE", false)
	config.SetDefault("FIBER_REDUCEMEMORYUSAGE", false)
}

func (config *Config) getFiberViewsEngine() fiber.Views {
	var viewsEngine fiber.Views
	config.Set("FIBER_VIEWS_EXTENSION", ".html")
	engine := html.NewFileSystem(http.Dir(config.GetString("FIBER_VIEWS_DIRECTORY")), config.GetString("FIBER_VIEWS_EXTENSION"))
	engine.Reload(config.GetBool("FIBER_VIEWS_RELOAD")).
		Debug(config.GetBool("FIBER_VIEWS_DEBUG")).
		Layout(config.GetString("FIBER_VIEWS_LAYOUT")).
		Delims(config.GetString("FIBER_VIEWS_DELIMS_L"), config.GetString("FIBER_VIEWS_DELIMS_R"))
	viewsEngine = engine
	return viewsEngine
}

func (config *Config) setFiberConfig() {
	config.fiber = &fiber.Config{
		Prefork:                   config.GetBool("FIBER_PREFORK"),
		ServerHeader:              config.GetString("FIBER_SERVERHEADER"),
		StrictRouting:             config.GetBool("FIBER_STRICTROUTING"),
		CaseSensitive:             config.GetBool("FIBER_CASESENSITIVE"),
		Immutable:                 config.GetBool("FIBER_IMMUTABLE"),
		UnescapePath:              config.GetBool("FIBER_UNESCAPEPATH"),
		ETag:                      config.GetBool("FIBER_ETAG"),
		BodyLimit:                 config.GetInt("FIBER_BODYLIMIT"),
		Concurrency:               config.GetInt("FIBER_CONCURRENCY"),
		Views:                     config.getFiberViewsEngine(),
		ReadTimeout:               config.GetDuration("FIBER_READTIMEOUT"),
		WriteTimeout:              config.GetDuration("FIBER_WRITETIMEOUT"),
		IdleTimeout:               config.GetDuration("FIBER_IDLETIMEOUT"),
		ReadBufferSize:            config.GetInt("FIBER_READBUFFERSIZE"),
		WriteBufferSize:           config.GetInt("FIBER_WRITEBUFFERSIZE"),
		CompressedFileSuffix:      config.GetString("FIBER_COMPRESSEDFILESUFFIX"),
		ProxyHeader:               config.GetString("FIBER_PROXYHEADER"),
		GETOnly:                   config.GetBool("FIBER_GETONLY"),
		ErrorHandler:              config.errorHandler,
		DisableKeepalive:          config.GetBool("FIBER_DISABLEKEEPALIVE"),
		DisableDefaultDate:        config.GetBool("FIBER_DISABLEDEFAULTDATE"),
		DisableDefaultContentType: config.GetBool("FIBER_DISABLEDEFAULTCONTENTTYPE"),
		DisableHeaderNormalizing:  config.GetBool("FIBER_DISABLEHEADERNORMALIZING"),
		DisableStartupMessage:     config.GetBool("FIBER_DISABLESTARTUPMESSAGE"),
		ReduceMemoryUsage:         config.GetBool("FIBER_REDUCEMEMORYUSAGE"),
	}
}

func (config *Config) GetFiberConfig() *fiber.Config {
	return config.fiber
}
