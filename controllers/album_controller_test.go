package controllers_test

import (
	"encoding/json"
	"example/golang-api/controllers"
	"example/golang-api/models"
	"example/golang-api/web"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func getSlogger() *zap.SugaredLogger {
  zapLogger, _ := zap.NewDevelopment()
  return zapLogger.Sugar()
}

func getAlbum2() models.Album {
  return models.Album{
    ID:     2,
    Title:  "Mock Title 2",
    Artist: "Mock Rock Star",
    Price:  19.95,
  }
}

func getAlbum4() models.Album {
  return models.Album{
    ID:     4,
    Title:  "Mock Title 4",
    Artist: "Mock Rock Star",
    Price:  24.95,
}
}

func getAlbumsData() []models.Album {
	var albums = []models.Album{getAlbum2(), getAlbum4()}
	return albums
}

func getEmptyAlbumsData() []models.Album {
	var albums = []models.Album{}
	return albums
}

type MockAlbumRepo struct {
	mock.Mock
}

func (m *MockAlbumRepo) GetAlbums() []models.Album {
  args := m.Called()
	return args.Get(0).([]models.Album)
}

func (m *MockAlbumRepo) PatchAlbum(models.Album) models.Album {
  args := m.Called()
	return args.Get(0).(models.Album)
}

func TestGetAlbumsShouldRespondStatus200(t *testing.T) {

  mockRepo := new(MockAlbumRepo)
  mockRepo.On("GetAlbums").Return(getAlbumsData())
	controller := controllers.NewAlbumController(mockRepo, getSlogger())

  router := web.NewRouter(getSlogger())
  web.RegisterAlbumController(router, controller)

	req, err := http.NewRequest("GET", "/albums", strings.NewReader(`{}`))
  require.Nil(t, err)

	w := httptest.NewRecorder()
  router.Engine.ServeHTTP(w, req)

  mockRepo.AssertNumberOfCalls(t, "GetAlbums", 1)
	assert.Equal(t, 200, w.Code, "expect status 200")
}

func TestGetAlbumsShouldRespondWithRepoData(t *testing.T) {

  mockRepo := new(MockAlbumRepo)
  mockRepo.On("GetAlbums").Return(getAlbumsData())
	controller := controllers.NewAlbumController(mockRepo, getSlogger())

  router := web.NewRouter(getSlogger())
  web.RegisterAlbumController(router, controller)

	req, err := http.NewRequest("GET", "/albums", strings.NewReader(`{}`))
  require.Nil(t, err)

  w := httptest.NewRecorder()
  router.Engine.ServeHTTP(w, req)
  require.Equal(t, 200, w.Code, "require status 200")

	var body []models.Album
	json.NewDecoder(w.Result().Body).Decode(&body)

  mockRepo.AssertNumberOfCalls(t, "GetAlbums", 1)
	assert.Equal(t, getAlbumsData(), body, "expect mock repo data")
}

func TestGetAlbumsShouldRespondWithRepoDataIfEmpty(t *testing.T) {

  mockRepo := new(MockAlbumRepo)
  mockRepo.On("GetAlbums").Return(getEmptyAlbumsData())
	controller := controllers.NewAlbumController(mockRepo, getSlogger())

  router := web.NewRouter(getSlogger())
  web.RegisterAlbumController(router, controller)

	req, err := http.NewRequest("GET", "/albums", strings.NewReader(`{}`))
  require.Nil(t, err)

  w := httptest.NewRecorder()
  router.Engine.ServeHTTP(w, req)
  require.Equal(t, 200, w.Code, "require status 200")

	var body []models.Album
	json.NewDecoder(w.Result().Body).Decode(&body)

  mockRepo.AssertNumberOfCalls(t, "GetAlbums", 1)
	assert.Equal(t, getEmptyAlbumsData(), body, "expect empty mock repo data")
}

func TestPatchAlbumsShouldRespondStatus200(t *testing.T) {

  mockRepo := new(MockAlbumRepo)
  mockRepo.On("PatchAlbum").Return(getAlbum2())
	controller := controllers.NewAlbumController(mockRepo, getSlogger())

  router := web.NewRouter(getSlogger())
  web.RegisterAlbumController(router, controller)

	req, err := http.NewRequest("PATCH", "/albums/2", strings.NewReader(`{"id": 2, "price": 44.44}`))
  require.Nil(t, err)

	w := httptest.NewRecorder()
  router.Engine.ServeHTTP(w, req)

  mockRepo.AssertNumberOfCalls(t, "PatchAlbum", 1)
	assert.Equal(t, 200, w.Code, "expect status 200")
}

func TestPatchAlbumsShouldRespondWithRepoData(t *testing.T) {

  mockRepo := new(MockAlbumRepo)
  mockRepo.On("PatchAlbum").Return(getAlbum2())
	controller := controllers.NewAlbumController(mockRepo, getSlogger())

  router := web.NewRouter(getSlogger())
  web.RegisterAlbumController(router, controller)

	req, err := http.NewRequest("PATCH", "/albums/2", strings.NewReader(`{"id": 2, "price": 44.44}`))
  require.Nil(t, err)

	w := httptest.NewRecorder()
  router.Engine.ServeHTTP(w, req)

	var body models.Album
	json.NewDecoder(w.Result().Body).Decode(&body)

  mockRepo.AssertNumberOfCalls(t, "PatchAlbum", 1)
	assert.Equal(t, getAlbum2(), body, "expect mock repo data")
}
