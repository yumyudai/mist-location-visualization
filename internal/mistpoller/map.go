package mistpoller

import (
	"encoding/json"
	"log"

	"mist-location-visualization/internal/mistdatafmt"
	"mist-location-visualization/internal/models"
)

func (s *PollAgent) updateDbEntryMap(mapData *mistdatafmt.ApiDataMapEntry) {
	// inject data to db
	ppm, _ := mapData.PPM.Float64()
	mapEntry := &models.Map{
		Name:   mapData.Name,
		Url:    mapData.Url,
		Id:     mapData.Id,
		SiteId: mapData.SiteId,
		Ppm:    ppm,
	}

	w, err := mapData.Width.Int64()
	if err != nil {
		log.Printf("map_engine: failed to convert width %v to int64 (%v)", mapData.Width, err)
	} else {
		// only update if successful convert
		mapEntry.Width = w
	}

	h, err := mapData.Height.Int64()
	if err != nil {
		log.Printf("map_engine: failed to convert height %v to int64 (%v)", mapData.Height, err)
	} else {
		// only update if successful convert
		mapEntry.Height = h
	}

	if s.Debug {
		s.DbConn.Debug().Save(mapEntry)
	} else {
		s.DbConn.Save(mapEntry)
	}

	return
}

func (s *PollAgent) processDataMap(data string) {
	// Get API response
	apiEntries := make([]*mistdatafmt.ApiDataMapEntry, 0)
	err := json.Unmarshal([]byte(data), &apiEntries)
	if err != nil {
		log.Printf("agent#%d: failed to parse JSON (%v)", s.Id, err)
		return
	}

	if s.Debug {
		log.Printf("agent#%d: got %d entries", s.Id, len(apiEntries))
	}

	// Get current data from DB
	delKeys := make(map[string]bool)
	dbEntries := make([]models.Map, 0)
	r := s.DbConn.Find(&dbEntries)
	if r.Error != nil {
		log.Printf("agent#%d: failed to fetch map data in DB (%v)", s.Id, r.Error)
		return
	}

	for _, dbEntry := range(dbEntries) {
		delKeys[dbEntry.Id] = true
	}

	// Check diff and update if necessary
	for _, apiEntry := range(apiEntries) {
		s.updateDbEntryMap(apiEntry)

		id, _ := apiEntry.GetJsonKeyValueAsStr("id")
		log.Printf("agent#%d: map id = %s has been updated", s.Id, id)
		delKeys[id] = false
	}

	// Delete keys if necesary
	for key, flag := range delKeys {
		if flag {
			err := s.DbConn.Delete(&models.Map{}, key)
			if err != nil {
				log.Printf("agent%d: failed to delete key %s", s.Id, key)
			} else if s.Debug {
				log.Printf("agent%d: deleted key %s", s.Id, key)
			}
		}
	}
}

