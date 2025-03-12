package locapiserver

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"mist-location-visualization/internal/models"
	"mist-location-visualization/internal/mistdatafmt"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// MistWebhookData represents the main webhook data structure from Mist
type MistWebhookData struct {
	Topic  string            `json:"topic"`
	Events []json.RawMessage `json:"events"`
}

// MistWhDataLocationAsset represents location data for an asset
type MistWhDataLocationAsset struct {
	Mac       string      `json:"mac"`
	SiteId    string      `json:"site_id"`
	MapId     string      `json:"map_id"`
	X         json.Number `json:"x"`
	Y         json.Number `json:"y"`
	Timestamp json.Number `json:"timestamp"`
}

// MistWhDataZone represents zone entry/exit data for an asset
type MistWhDataZone struct {
	Mac     string `json:"mac"`
	AssetId string `json:"asset_id"`
	MapId   string `json:"map_id"`
	Trigger string `json:"trigger"`
	ZoneId  string `json:"zone_id"`
}

func (s *LocApiServer) apiMistRecvRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/", s.apiMistRecvPost)

	return r
}

// buildURL constructs a properly formatted URL with the given endpoint and URI
func buildURL(endpoint string, uri string) string {
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		return "https://" + endpoint + uri
	}
	return endpoint + uri
}

func (s *LocApiServer) doMistAssetSearchCall(siteid string, mac string) (string, error) {
	// build request
	uri := fmt.Sprintf("/api/v1/sites/%s/stats/assets/search?mac=%s", siteid, mac)
	reqURL := buildURL(s.cfg.Mist.Endpoint, uri)
	
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return "", err
	}

	// set authentication header
	tokenStr := fmt.Sprintf("token %s", s.cfg.Mist.Apikey)
	req.Header.Set("Authorization", tokenStr)

	// start requet
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// read response
	log.Printf("doMistAssetSearchCall: GET %s (response %-v)", reqURL, resp)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (s *LocApiServer) fetchAssetData(siteid string, mac string) (*mistdatafmt.ApiDataAssetEntry, error) {
	r, err := s.doMistAssetSearchCall(siteid, mac)
	if err != nil {
		return nil, fmt.Errorf("asset search call failed: %w", err)
	}

	apiResult := mistdatafmt.ApiDataAssetSearchResult{}
	err = json.Unmarshal([]byte(r), &apiResult)
	if err != nil {
		return nil, fmt.Errorf("failed to parse asset data: %w", err)
	}

	resultTotal, _ := apiResult.Total.Int64()
	if resultTotal > 1 {
		log.Printf("fetchAssetData: Warning: More than 1 asset found for mac %s in site %s", mac, siteid)
	}

	if s.cfg.Mist.Debug {
		log.Printf("Asset data: %s", r)
	}
	
	if resultTotal < 1 {
		return nil, fmt.Errorf("mac %s not found", mac)
	} else if len(apiResult.Results) < 1 {
		return nil, fmt.Errorf("invalid result length %d", len(apiResult.Results))
	}

	return &(apiResult.Results[0]), nil
}

func (s *LocApiServer) handleWhInLocationAsset(dataIn MistWhDataLocationAsset) {
	// Make sure we have Map information
	mapEntry := models.Map{}
	ret := s.dbConn.Where(&models.Map{Id: dataIn.MapId}).First(&mapEntry)
	if ret.Error != nil {
		log.Printf("handleWhInLocationAsset: Failed to query DB (%v)", ret.Error)
		return
	}

	// Update or Create?
	dbEntry := models.Entity{}
	s.dbConn.Where(&models.Entity{Mac: dataIn.Mac}).First(&dbEntry)

	x, _ := dataIn.X.Float64()
	y, _ := dataIn.Y.Float64()
	px := x * mapEntry.Ppm
	py := y * mapEntry.Ppm

	dbEntry.Mac = dataIn.Mac
	dbEntry.MapId = dataIn.MapId
	dbEntry.X = px
	dbEntry.Y = py
	dbEntry.Lastseen, _ = dataIn.Timestamp.Float64()

	// Fetch name
	tNow := time.Now()
	refreshDuration := time.Duration(s.cfg.Mist.RefreshTime) * time.Second
	tExpire := dbEntry.LastRefresh.Add(refreshDuration)
	if tNow.After(tExpire) {
		apidata, err := s.fetchAssetData(dataIn.SiteId, dataIn.Mac)
		if err != nil {
			log.Printf("handleWhInLocationAsset: Failed to fetch client name (%v)", err)
		} else {
			// poor man's display name
			regPattern := `\[(?P<org>.+)\] (?P<name>.+)`
			re, err := regexp.Compile(regPattern)
			if err != nil {
				log.Printf("handleWhInLocationAsset: failed to compile regexp (%v)", err)
			} else {
				reMatch := re.FindAllStringSubmatch(apidata.Name, -1)
				if len(reMatch) > 0 {
					dbEntry.DisplayName = reMatch[0][re.SubexpIndex("name")]
					dbEntry.DisplayOrg = reMatch[0][re.SubexpIndex("org")]
				}
			}
			
			dbEntry.Name = apidata.Name
		}
		dbEntry.LastRefresh = tNow
	}

	s.dbConn.Debug().Save(&dbEntry)

	return
}

func (s *LocApiServer) handleWhInZone(dataIn MistWhDataZone) {
	dbEntry := models.Entity{}
	ret := s.dbConn.Where(&models.Entity{Mac: dataIn.Mac}).First(&dbEntry)
	if ret.Error != nil {
		log.Printf("handleWhInZone: Failed to query DB (%v)", ret.Error)
		return
	}

	switch dataIn.Trigger {
	case "enter":
		zone := models.Zone{}
		result := s.dbConn.Where(&models.Zone{Id: dataIn.ZoneId}).First(&zone)
		if result.Error != nil {
			log.Printf("handleWhInZone: Failed to query zone data (%v)", result.Error)
			return
		}
		dbEntry.ZoneName = zone.Name
		dbEntry.ZoneId = dataIn.ZoneId

	case "exit":
		dbEntry.ZoneName = ""
		dbEntry.ZoneId = ""
	}

	s.dbConn.Debug().Save(&dbEntry)

	return
}

func (s *LocApiServer) apiMistRecvAuthenticate(inputSig string, body []byte) bool {
	secret := []byte(s.cfg.Mist.Secret)

	// Create an HMAC hasher, write the body to it, and compute the expected signature
	h := hmac.New(sha256.New, secret)
	h.Write(body)
	expectedSig := hex.EncodeToString(h.Sum(nil))

	if inputSig != expectedSig {
		log.Printf("unexpected signature %s", inputSig)
		return false
	}
	
	return true
}

func (s *LocApiServer) apiMistRecvPost(w http.ResponseWriter, r *http.Request) {
	// get data
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("apiMistRecvPost: Failed to read request body: %v", err)
		render.Render(w, r, s.httpErrUnexpected(err))
		return
	}

	// authenticate
	if s.cfg.Mist.Secret != "" {
		sig := r.Header.Get("x-mist-signature-v2")
		if !s.apiMistRecvAuthenticate(sig, body) {
			err := fmt.Errorf("invalid signature")
			render.Render(w, r, s.httpErrUnauthorized(err))
			return
		}
	}

	// process data
	dataIn := MistWebhookData{}
	err = json.Unmarshal(body, &dataIn)
	if err != nil {
		log.Printf("apiMistRecvPost: Failed to parse webhook data: %v", err)
		render.Render(w, r, s.httpErrInvalidRequest(err))
		return
	}

	switch dataIn.Topic {
	case "location-asset":
		for _, ev := range dataIn.Events {
			evData := MistWhDataLocationAsset{}
			err := json.Unmarshal(ev, &evData)
			if err != nil {
				log.Printf("failed to decode webhook input %s (%v)", string(ev), err)
				continue
			}

			s.handleWhInLocationAsset(evData)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(nil)

	case "zone":
		for _, ev := range dataIn.Events {
			evData := MistWhDataZone{}
			err := json.Unmarshal(ev, &evData)
			if err != nil {
				log.Printf("failed to decode webhook input %s (%v)", string(ev), err)
				continue
			}

			s.handleWhInZone(evData)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(nil)

	default:
		log.Printf("unsupported topic: %s", dataIn.Topic)
		w.WriteHeader(http.StatusOK)
		w.Write(nil)
	}

	return
}

