package inputParams

import (
	"bank-demo-app/internal/mongodb"
	"flag"
)

type AppConfig struct {
	InMemory  bool
	MongoConf mongodb.MongoConfig
}

func ParseInputParams() (*AppConfig, error) {
	config := new(AppConfig)
	flag.BoolVar(&config.InMemory, "in-memory", true, "Run the application in memory (no database)")
	config.MongoConf.AddFlagsParams()
	flag.Parse()

	// Add more params check at production systems. For this test, DBName is enough.
	if !config.InMemory && config.MongoConf.DbName == "" {
		return &AppConfig{}, nil
	}

	return config, nil
}
