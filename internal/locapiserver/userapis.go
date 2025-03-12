package locapiserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"mist-location-visualization/internal/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// EntityExtView represents the external view of an entity for API responses
type EntityExtView struct {
	Id          string  `json:"id"`
	MapId       string  `json:"map_id"`
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	Lastseen    int64   `json:"last_seen"`
	ZoneName    string  `json:"zone_name"`
	DisplayName string  `json:"display_name"`
	DisplayOrg  string  `json:"display_org"`
}

func (e *EntityExtView) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *LocApiServer) apiEntityRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", s.apiEntityGetAll)

	return r
}

func (s *LocApiServer) apiEntityGetAll(w http.ResponseWriter, r *http.Request) {
	entities := make([]models.Entity, 0)
	ret := s.dbConn.Find(&entities)
	if ret.Error != nil {
		log.Printf("apiEntityGetAll: Failed to query DB (%v)", ret.Error)
		err := fmt.Errorf("failed to get data from backend")
		render.Render(w, r, s.httpErrUnexpected(err))
		return
	}

	outs := []render.Renderer{}
	for _, e := range entities {
		// check timeout
		tNow := time.Now()
		timeoutDuration := time.Duration(s.cfg.Mist.LocationTimeout) * time.Second
		tExpire := time.Unix(int64(e.Lastseen), 0).Add(timeoutDuration)
		if tNow.After(tExpire) && e.X != -1 && e.Y != -1 {
			log.Printf("apiEntityGetAll: Mac %s has timed out", e.Mac)

			e.X = -1
			e.Y = -1
			e.MapId = ""
			e.ZoneId = ""
			e.ZoneName = ""
			s.dbConn.Debug().Save(&e)
		}

		o := &EntityExtView{
			Id:          e.Mac,
			MapId:       e.MapId,
			X:           e.X,
			Y:           e.Y,
			Lastseen:    int64(e.Lastseen),
			ZoneName:    e.ZoneName,
			DisplayName: e.DisplayName,
			DisplayOrg:  e.DisplayOrg,
		}

		outs = append(outs, o)
	}

	render.RenderList(w, r, outs)
	return
}

/* Search Query */
func (s *LocApiServer) apiSearchRouter() chi.Router {
	r := chi.NewRouter()
	return r
}

// MapExtView represents the external view of a map for API responses
type MapExtView struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Width  int64  `json:"width"`
	Height int64  `json:"height"`
}

func (e *MapExtView) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *LocApiServer) apiMapIdCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "mapid")
		if key == "" {
			err := fmt.Errorf("Missing mapid param")
			render.Render(w, r, s.httpErrInvalidRequest(err))
			return
		}

		ctx := context.WithValue(r.Context(), "mapid", key)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *LocApiServer) apiMapRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", s.apiMapGetAll)
	r.Route("/{mapid}", func(r chi.Router) {
		r.Use(s.apiMapIdCtx)
		r.Get("/zone", s.apiMapGetZone)
	})

	return r
}

func (s *LocApiServer) apiMapGetAll(w http.ResponseWriter, r *http.Request) {
	maps := make([]models.Map, 0)
	ret := s.dbConn.Find(&maps)
	if ret.Error != nil {
		log.Printf("apiMapGetAll: Failed to query DB (%v)", ret.Error)
		err := fmt.Errorf("Failed to get data from backend")
		render.Render(w, r, s.httpErrUnexpected(err))
		return
	}

	outs := []render.Renderer{}
	for _, e := range maps {
		o := &MapExtView{
			Id:     e.Id,
			Name:   e.Name,
			Height: e.Height,
			Width:  e.Width,
		}

		outs = append(outs, o)
	}

	render.RenderList(w, r, outs)
	return
}

func (s *LocApiServer) apiMapGetZone(w http.ResponseWriter, r *http.Request) {
	mapId := getCtxValueString(r.Context(), "mapid")
	zones := make([]models.Zone, 0)
	ret := s.dbConn.Where("map_id = ?", mapId).Find(&zones)
	if ret.Error != nil {
		log.Printf("apiMapGetZone: Failed to query DB (%v)", ret.Error)
		err := fmt.Errorf("failed to get data from backend")
		render.Render(w, r, s.httpErrUnexpected(err))
		return
	}

	outs := []render.Renderer{}
	for _, e := range zones {
		var count int64

		result := s.dbConn.Model(&models.Entity{}).Where("zone_id = ?", e.Id).Count(&count)
		if result.Error != nil {
			log.Printf("apiMapGetZone: Failed to query DB on count (%v)", result.Error)
			count = 0
		}

		o := &ZoneExtView{
			Id:    e.Id,
			Name:  e.Name,
			MapId: e.MapId,
			Count: count,
		}

		outs = append(outs, o)
	}

	render.RenderList(w, r, outs)
	return
}

// ZoneExtView represents the external view of a zone for API responses
type ZoneExtView struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	MapId string `json:"map_id"`
	Count int64  `json:"count"`
}

func (e *ZoneExtView) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *LocApiServer) apiZoneRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", s.apiZoneGetAll)

	return r
}

func (s *LocApiServer) apiZoneGetAll(w http.ResponseWriter, r *http.Request) {
	zones := make([]models.Zone, 0)
	ret := s.dbConn.Find(&zones)
	if ret.Error != nil {
		log.Printf("apiZoneGetAll: Failed to query DB (%v)", ret.Error)
		err := fmt.Errorf("failed to get data from backend")
		render.Render(w, r, s.httpErrUnexpected(err))
		return
	}

	outs := []render.Renderer{}
	for _, e := range zones {
		var count int64

		result := s.dbConn.Model(&models.Entity{}).Where("zone_id = ?", e.Id).Count(&count)
		if result.Error != nil {
			log.Printf("apiZoneGetAll: Failed to query DB on count (%v)", result.Error)
			count = 0
		}

		o := &ZoneExtView{
			Id:    e.Id,
			Name:  e.Name,
			MapId: e.MapId,
			Count: count,
		}

		outs = append(outs, o)
	}

	render.RenderList(w, r, outs)
	return
}

