package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"chatserver/logger"
	"chatserver/server"
)

const (
	defaultHost        = "localhost"
	defaultPort        = 3001
	defaultLogLocation = "/tmp/chatserver.log"
)

var (
	configPath = flag.String("config", "", "path to config file")
)

type serverConfig struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	LogLocation string `json:"logLocation"`
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	conf, err := loadConfig(*configPath)
	if err != nil {
		return err
	}

	fileLogger := logger.NewFileLogger(conf.LogLocation)
	if err := fileLogger.Open(); err != nil {
		return err
	}
	defer fileLogger.Close()

	server := server.NewServer(conf.Host, conf.Port, fileLogger)
	return server.Serve()
}

func loadConfig(configPath string) (*serverConfig, error) {
	conf := &serverConfig{
		Host:        defaultHost,
		Port:        defaultPort,
		LogLocation: defaultLogLocation,
	}

	if configPath == "" {
		// you prefer the default config? you have chosen well...
		return conf, nil
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
