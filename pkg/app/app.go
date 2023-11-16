package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger" // gin-swagger middleware

	"github.com/klaxn/klaxn-api/docs"
	"github.com/klaxn/klaxn-api/internal/config"
	"github.com/klaxn/klaxn-api/internal/data"
	"github.com/klaxn/klaxn-api/internal/routes"
	"github.com/klaxn/klaxn-api/pkg/outbound"
)

type App struct {
	Logger  logrus.FieldLogger
	manager *data.Manager
	ob      []outbound.Sender
	conf    *config.Config
}

func New() (*App, error) {
	conf, err := config.New()
	if err != nil {
		return nil, err
	}

	manager, err := data.New(conf.DatabaseConfig)
	if err != nil {
		return nil, err
	}
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	return &App{Logger: logger, manager: manager, ob: outbound.CreateOutbounds(conf.OutboundConfig, logger), conf: conf}, nil
}

func (a *App) Run() error {
	router := routes.New(a.manager, a.ob, a.Logger, a.conf)
	r := gin.Default()

	r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Users
	r.GET("/api/users", router.GetUsers)
	r.GET("/api/users/:id", router.GetUser)
	r.POST("/api/users", router.CreateUser)
	r.PUT("/api/users/:id", router.UpdateUser)
	r.DELETE("/api/users/:id", router.DeleteUser)

	// Teams
	r.GET("/api/teams", router.GetTeams)
	r.GET("/api/teams/:id", router.GetTeam)
	r.POST("/api/teams", router.CreateTeam)
	r.PUT("/api/teams/:id", router.UpdateTeam)
	r.DELETE("/api/teams/:id", router.DeleteTeam)

	// Services
	r.GET("/api/services", router.GetServices)
	r.GET("/api/services/:id", router.GetService)
	r.POST("/api/services", router.CreateService)
	r.PUT("/api/services/:id", router.UpdateService)
	r.DELETE("/api/services/:id", router.DeleteService)

	// Alert
	r.POST("/api/alerts/grafana", router.GrafanaInbound)
	r.GET("/api/alerts", router.GetAlerts)

	// Escalations
	r.GET("/api/escalations", router.GetEscalations)
	r.GET("/api/escalations/:id", router.GetEscalation)
	r.PUT("/api/escalations/:id", router.UpdateEscalation)
	r.POST("/api/escalations", router.CreateEscalation)
	r.DELETE("/api/escalations/:id", router.DeleteEscalation)

	docs.SwaggerInfo.BasePath = "/api"

	// Debug
	r.GET("/api/debug/config", router.GetConfig)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r.Run()
}
