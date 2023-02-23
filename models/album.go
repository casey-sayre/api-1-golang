package models

type AlbumRepositoryInterface interface {
  GetAlbums() []Album
  PatchAlbum(Album) Album
}

type Album struct {
	ID     int  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}
