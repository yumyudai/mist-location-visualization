package locapiserver

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"mist-location-visualization/internal/models"
)

type LocApiServer struct {
	cfg    Config
	dbConn *gorm.DB
}

/* Main */
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

func New(cfg Config) (*LocApiServer, error) {
	var err error

	// Base Initialization
	r := &LocApiServer{
		cfg: cfg,
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


	return r, nil
}

func (s *LocApiServer) Run() error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	if s.cfg.Http.BasicAuth {
		userdb := make(map[string]string)
		for _, v := range s.cfg.Http.Users {
			userdb[v.User] = v.Password
		}
		r.Use(middleware.BasicAuth(s.cfg.Http.ServerName, userdb))
	}

	r.Route("/entity", func(r chi.Router) {
		r.Mount("/", s.apiEntityRouter())
	})

	r.Route("/zone", func(r chi.Router) {
		r.Mount("/", s.apiZoneRouter())
	})

	r.Route("/search", func(r chi.Router) {
		r.Mount("/", s.apiSearchRouter())
	})

	r.Route("/map", func(r chi.Router) {
		r.Mount("/", s.apiMapRouter())
	})

	r.Route("/mistrecv", func(r chi.Router) {
		r.Mount("/", s.apiMistRecvRouter())
	})

	// Start HTTP Handler
	err := http.ListenAndServe(s.cfg.Http.Listen, r)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
