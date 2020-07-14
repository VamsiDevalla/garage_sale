package main

import (
	"flag"
	"log"
	"os"

	"github.com/VamsiDevalla/garage_sale/internal/platform/conf"
	"github.com/VamsiDevalla/garage_sale/internal/platform/database"
	"github.com/VamsiDevalla/garage_sale/internal/schema"
)

func main() {

	var cfg struct {
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:postgres,noprint"`
			Host       string `conf:"default:localhost"`
			Path       string `conf:"default:postgres"`
			DisableSSL bool   `conf:"default:true"`
		}
	}

	if err := conf.Parse(os.Args[1:], "SALES", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				log.Fatalf("error : generating config usage : %v", err)
			}
			log.Println(usage)
			return
		}
		log.Fatalf("error: parsing config: %s", err)
	}

	// =========================================================================
	// Start Database
	db, err := database.Open(database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Path:       cfg.DB.Path,
		DisableSSL: cfg.DB.DisableSSL,
	})
	if err != nil {
		log.Fatalf("error: connecting to db: %s", err)
	}
	defer db.Close()

	flag.Parse()

	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			log.Println("error applying migrations", err)
			os.Exit(1)
		}
		log.Println("Migrations complete")
		return

	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Println("error seeding database", err)
			os.Exit(1)
		}
		log.Println("Seed data complete")
		return
	}
}
