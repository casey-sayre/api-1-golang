package controllers

import (
	"encoding/json"
	"example/golang-api/models"
	"net/http"

	snsRepo "example/golang-api/repositories/sns"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AlbumController struct {
	repo models.AlbumRepositoryInterface
  albumUpdatesPublisher *snsRepo.AlbumUpdatesPublisher
  slog *zap.SugaredLogger
}

func NewAlbumController(repo models.AlbumRepositoryInterface, albumUpdatesPublisher *snsRepo.AlbumUpdatesPublisher, slogger *zap.SugaredLogger) *AlbumController {
	controller := AlbumController{
		repo: repo,
    slog: slogger,
    albumUpdatesPublisher: albumUpdatesPublisher,
	}
	return &controller
}

func (controller AlbumController) GetAlbums(c *gin.Context) {

  controller.slog.Info("controllers: getAlbums")

	var albums []models.Album = controller.repo.GetAlbums()

	c.IndentedJSON(http.StatusOK, albums)
}

func (controller AlbumController) PatchAlbum(c *gin.Context) {

	controller.slog.Info("controllers: patchAlbum")

	var requestAlbum models.Album

	if err := c.BindJSON(&requestAlbum); err != nil {
		controller.slog.Warn("patch album unable to bind to request")
		c.AbortWithStatus(400)
		return
	}

	updatedAlbum := controller.repo.PatchAlbum(requestAlbum)

  updatedAlbumJsonData, _ := json.Marshal(updatedAlbum)

  controller.albumUpdatesPublisher.PublishUpdatedAlbum(string(updatedAlbumJsonData))

	c.IndentedJSON(http.StatusOK, updatedAlbum)
}
