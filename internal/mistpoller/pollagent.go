package mistpoller

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

// buildURL constructs a properly formatted URL with the given endpoint and URI
func buildURL(endpoint string, uri string) string {
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		return "https://" + endpoint + uri
	}
	return endpoint + uri
}

type PollAgent struct {
	Id		int
	DbConn		*gorm.DB
	Endpoint	string
	Apikey		string
	Uri		string
	Layout		string
	Interval	int
	Debug		bool

	intvlTicker	*time.Ticker
	killSig		chan struct{}
	wg		*sync.WaitGroup
}


func (s *PollAgent) runRequest() {
	// build request
	reqURL := buildURL(s.Endpoint, s.Uri)
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		log.Printf("agent#%d: failed to build HTTP request (%v)", s.Id, err)
		return
	}

	// set authentication header
	tokenStr := fmt.Sprintf("token %s", s.Apikey)
	req.Header.Set("Authorization", tokenStr)

	// start request
	if s.Debug {
		log.Printf("agent#%d: start HTTP GET request: url %s", s.Id, reqURL)
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("agent#%d: HTTP request failure (%v)", s.Id, err)
		return
	}
	defer resp.Body.Close()

	// read response
	if s.Debug {
		log.Printf("agent#%d: got HTTP response: %-v", s.Id, resp)
	}

	if resp.StatusCode != 200 {
		log.Printf("agent#%d: HTTP request has returned status code %d", s.Id, resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("agent#%d: Failed to read HTTP response body (%v)", s.Id, err)
		return
	}

	s.processData(string(body))
	return
}


func (s *PollAgent) processData(data string) {
	switch(s.Layout) {
	case "maps":
		s.processDataMap(data)

	case "zones":
		s.processDataZone(data)

	default:
		log.Printf("agent#%d: unknown data layout %s", s.Id, s.Layout)
		return
	}

	return
}

func (s *PollAgent) finish() {
	if s.intvlTicker != nil {
		s.intvlTicker.Stop()
	}

	if s.wg != nil {
		s.wg.Done()
	}

	log.Printf("agent#%d: finished process thread", s.Id)

	return
}

func (s *PollAgent) Run(wg *sync.WaitGroup, killSig chan struct{}) error {
	log.Printf("agent#%d: start poll agent thread (uri %s, interval %d)", s.Id, s.Uri, s.Interval)

	// init
	s.intvlTicker = time.NewTicker(time.Duration(s.Interval) * time.Second)
	s.killSig = killSig
	s.wg = wg

	// start
	wg.Add(1)
	defer s.finish()

	s.runRequest()
	for {
		select {
		case <-killSig:
			return nil
		case <-s.intvlTicker.C:
			s.runRequest()
		}
	}

	return nil
}

