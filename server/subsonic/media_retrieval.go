package subsonic

import (
	"io"
	"net/http"
	"os"

	"github.com/deluan/navidrome/engine"
	"github.com/deluan/navidrome/log"
	"github.com/deluan/navidrome/model"
	"github.com/deluan/navidrome/server/subsonic/responses"
)

type MediaRetrievalController struct {
	cover engine.Cover
}

func NewMediaRetrievalController(cover engine.Cover) *MediaRetrievalController {
	return &MediaRetrievalController{cover: cover}
}

func (c *MediaRetrievalController) GetAvatar(w http.ResponseWriter, r *http.Request) (*responses.Subsonic, error) {
	var f *os.File
	f, err := os.Open("static/itunes.png")
	if err != nil {
		log.Error(r, "Image not found", err)
		return nil, NewError(responses.ErrorDataNotFound, "Avatar image not found")
	}
	defer f.Close()
	io.Copy(w, f)

	return nil, nil
}

func (c *MediaRetrievalController) GetCoverArt(w http.ResponseWriter, r *http.Request) (*responses.Subsonic, error) {
	id, err := RequiredParamString(r, "id", "id parameter required")
	if err != nil {
		return nil, err
	}
	size := ParamInt(r, "size", 0)

	err = c.cover.Get(r.Context(), id, size, w)

	switch {
	case err == model.ErrNotFound:
		log.Error(r, "Couldn't find coverArt", "id", id, err)
		return nil, NewError(responses.ErrorDataNotFound, "Cover not found")
	case err != nil:
		log.Error(r, "Error retrieving coverArt", "id", id, err)
		return nil, NewError(responses.ErrorGeneric, "Internal Error")
	}

	return nil, nil
}