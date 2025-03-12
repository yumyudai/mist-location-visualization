package mistpoller

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"mist-location-visualization/internal/models"
)

type Poller struct {
	cfg	Config

	dbConn *gorm.DB
	agents	[]*PollAgent
	wg	*sync.WaitGroup
}

func getDbConn(cfg Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch cfg.Db.Driver {
	case "mysql":
		if cfg.Db.Mysql.User == "" || cfg.Db.Mysql.Host == "" || cfg.Db.Mysql.Database == "" {
			return nil, fmt.Errorf("missing connection info")
		}

		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Db.Mysql.User, cfg.Db.Mysql.Password, cfg.Db.Mysql.Host, cfg.Db.Mysql.Database)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unknown db driver %s", cfg.Db.Driver)
	}

	if cfg.Db.Debug {
		db.Logger = db.Logger.LogMode(logger.Info)
	}

	return db, err
}

func New(cfg Config) (*Poller, error) {
	var err error

	// Base Initialization
	r := &Poller {
		cfg:	cfg,
		agents:	make([]*PollAgent, 0),
		wg:	&sync.WaitGroup{},
	}

	// DB Conn Initialization
	r.dbConn, err = getDbConn(cfg)
	if err != nil {
		return nil, err
	}

	err = r.dbConn.Debug().AutoMigrate(&models.Map{})
	if err != nil {
		log.Printf("failed to automigrate database %v", err)
		return nil, err
	}

	err = r.dbConn.Debug().AutoMigrate(&models.Zone{})
	if err != nil {
		log.Printf("failed to automigrate database %v", err)
		return nil, err
	}

	err = r.dbConn.Debug().AutoMigrate(&models.Entity{})
	if err != nil {
		log.Printf("failed to automigrate database %v", err)
		return nil, err
	}

	// Poll Agent Initialization
	for id, v := range(cfg.Datasource) {
		agent := &PollAgent {
			Id:		id,
			DbConn:		r.dbConn,
			Endpoint:	cfg.Mist.Endpoint,
			Apikey:		cfg.Mist.Apikey,
			Uri:		v.Uri,
			Layout:		v.Datalayout,
			Interval:	v.Interval,
			Debug:		cfg.Mist.Debug,
		}
	
		r.agents = append(r.agents, agent)
		id++
	}

	return r, nil 
}

func (s *Poller) Run() error {
	var shutdownSigs []chan struct{}
	// Launch
	for _, agent := range(s.agents) {
		agentShutdownSig := make(chan struct{}, 1)
		shutdownSigs = append(shutdownSigs, agentShutdownSig)
		go agent.Run(s.wg, agentShutdownSig)
	}

	// Main thread to wait until we get a kill signal or something go wrong
	killSig := make(chan os.Signal, 1)
	signal.Notify(killSig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-killSig

	log.Printf("Caught kill signal, shutting down")
	for _, sig := range(shutdownSigs) {
		close(sig)
	}
	s.wg.Wait()

	log.Printf("All threads exited")

	return nil
}
