package repositories

import (
	"example/golang-api/config"
	"example/golang-api/models"
	"fmt"
	"reflect"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresRepo struct {
	DB   *gorm.DB
	slog *zap.SugaredLogger
}

func ProvideAlbumRepository(config *config.Config, slogger *zap.SugaredLogger) models.AlbumRepositoryInterface {
	return NewPostgresRepo(config, slogger)
}

func NewPostgresRepo(config *config.Config, slogger *zap.SugaredLogger) *PostgresRepo {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/New_York", config.Db.Host, config.Db.User, config.Db.Password, config.Db.DbName, config.Db.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	return &PostgresRepo{DB: db, slog: slogger}
}

func (repo PostgresRepo) GetAlbums() []models.Album {

	repo.slog.Info("postgres GetAlbums")

	var albums []models.Album

	repo.DB.Order("artist, title").Find(&albums)

	return albums
}

func (repo PostgresRepo) PatchAlbum(album models.Album) models.Album {

	repo.slog.Info("postgres PatchAlbums")

	var dbAlbum models.Album
	repo.DB.First(&dbAlbum, album.ID)

	if !reflect.ValueOf(album.Title).IsZero() {
		dbAlbum.Title = album.Title
	}

	if !reflect.ValueOf(album.Artist).IsZero() {
		dbAlbum.Artist = album.Artist
	}

	if !reflect.ValueOf(album.Price).IsZero() {
		dbAlbum.Price = album.Price
	}

	repo.DB.Save(&dbAlbum)

	return dbAlbum
}
