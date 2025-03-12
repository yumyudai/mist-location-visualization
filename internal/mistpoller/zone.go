package mistpoller

import (
	"encoding/json"
	"log"

	"mist-location-visualization/internal/mistdatafmt"
	"mist-location-visualization/internal/models"
)

func (s *PollAgent) updateDbEntryZone(zoneData *mistdatafmt.ApiDataZoneEntry) error {
	// inject data to db
	dbEntry := &models.Zone {
			Id:		zoneData.Id,
			MapId:		zoneData.MapId,
			SiteId:		zoneData.SiteId,
			Name:		zoneData.Name,
	}

	if s.Debug {
		s.DbConn.Debug().Save(dbEntry)
	} else {
		s.DbConn.Save(dbEntry)
	}

	return nil
}
func (s *PollAgent) processDataZone(data string) {
	// Get API response
	apiEntries := make([]*mistdatafmt.ApiDataZoneEntry, 0)
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
	dbEntries := make([]models.Zone, 0)
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
		s.updateDbEntryZone(apiEntry)

		id, _ := apiEntry.GetJsonKeyValueAsStr("id")
		log.Printf("agent#%d: map id = %s has been updated", s.Id, id)
		delKeys[id] = false
	}

	// Delete keys if necesary
	for key, flag := range delKeys {
		if flag {
			err := s.DbConn.Delete(&models.Zone{}, key)
			if err != nil {
				log.Printf("agent%d: failed to delete key %s", s.Id, key)
			} else if s.Debug {
				log.Printf("agent#%d: deleted key %s", s.Id, key)
			}
		}
	}
}

