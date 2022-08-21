package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"youtube-downloader-rest/config"
	// Swagger docs
	_ "youtube-downloader-rest/docs"
	"youtube-downloader-rest/pb"
	"youtube-downloader-rest/pkg/logger"
)

type Server struct {
	storageClient pb.StorageClient
	logger        logger.Logger
	cfg           *config.Config
}

func NewServer(storageClient pb.StorageClient, l logger.Logger, cfg *config.Config) *Server {
	return &Server{
		storageClient: storageClient,
		logger:        l,
		cfg:           cfg,
	}
}

func (s *Server) Run() {
	h := echo.New()

	// Options
	h.Use(middleware.Logger())
	h.Use(middleware.Recover())

	// Swagger
	h.GET("/swagger/*", echoSwagger.WrapHandler)

	// Routes
	gh := h.Group("/api/v1")
	{
		gh.POST("/generate_image", s.doGetInfo)
	}

	httpServer := &http.Server{
		Handler: h,
		Addr:    fmt.Sprintf(":%s", s.cfg.Http.Port),
	}

	httpServer.ListenAndServe()
}

type GetInfoRequest struct {
	ImageName string `json:"image_name"`
}

// @Summary     Generate random image
// @Description Creates new random image with provided name
// @Param 		image_name body string true "image name"
// @ID          generate-img
// @Accept      json
// @Produce     json
// @Router      /api/v1/generate_image [post]
func (s *Server) doGetInfo(c echo.Context) error {
	var payload *GetInfoRequest

	if err := (&echo.DefaultBinder{}).BindBody(c, &payload); err != nil {
		return err
	}
	f, err := s.storageClient.CreateTask(context.Background(), &pb.Task{ImageName: payload.ImageName})
	if err != nil {
		return fmt.Errorf("failed to create task: %v", err)
	}

	return c.JSON(http.StatusOK, f.Name)
}
