package routes

import (
	"backend/internal/config"
	"backend/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type Handler struct {
	Services *service.Service
	cfg      *config.Config
}

func NewHandler(service *service.Service, cfg *config.Config) *Handler {
	return &Handler{service, cfg}
}

func (h *Handler) Init(cfg *config.Config) *echo.Echo {
	// Init echo handler
	router := echo.New()
	// Init middleware
	router.Use(
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}, bytes_in=${bytes_in}, bytes_out=${bytes_out}\n",
			Output: router.Logger.Output()}),
		middleware.Recover())

	// Init log level
	router.Debug = cfg.ServerMode != config.Dev

	//docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	//if cfg.ServerMode != config.Dev {
	//	docs.SwaggerInfo.Host = cfg.Host
	//}
	//
	//if cfg.ServerMode != config.Prod {
	//	router.GET("/swagger/*", echoSwagger.WrapHandler)
	//}

	// Init router
	router.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	router.GET("/categories", h.Categories)

	return router
}

func (h *Handler) Categories(c echo.Context) error {
	resp, err := h.Services.DB.GetAllCategories()
	if err != nil {
		return c.JSON(500, err.Error())
	}

	return c.JSON(200, resp)
}
