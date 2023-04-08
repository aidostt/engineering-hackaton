package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-playground/form"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	config config
	//models        data.Models
	formDecoder *form.Decoder
	//templateCache map[string]*template.Template
}

type config struct {
	port int
	env  string
	db   struct {
		dsn         string
		maxOpenConn uint64
		maxIdleTime time.Duration
	}
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development | staging | production")
	flag.StringVar(&cfg.db.dsn, "mongo_url", os.Getenv("MONGO_DB_DSN"), "db connection string")
	flag.Uint64Var(&cfg.db.maxOpenConn, "maxOpenConn", uint64(100), "maximum number of open connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "maxIdleTime", time.Duration(10), "maximum idle time of one connection")
	flag.Parse()
	client, err := openDB(cfg)
	if err != nil {
		log.Fatal(err, nil)
	}
	log.Println("database connection pool established", nil)
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err, nil)
		}
	}()
	decoder := form.NewDecoder()
	app := application{
		config:      cfg,
		formDecoder: decoder,
	}
	srv := http.Server{
		Addr:    fmt.Sprintf("localhost:%v", cfg.port),
		Handler: app.Router(),
		//ErrorLog: log.New(logger, "", 0), //nigga
	}

	log.Printf("starting server on address %v, environment %v", srv.Addr, cfg.env)

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func openDB(cfg config) (*mongo.Client, error) {
	ClientOptions := options.ClientOptions{
		MaxPoolSize:     &cfg.db.maxOpenConn,
		MaxConnIdleTime: &cfg.db.maxIdleTime,
	}

	client, err := mongo.Connect(context.TODO(), ClientOptions.ApplyURI(cfg.db.dsn))
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}
