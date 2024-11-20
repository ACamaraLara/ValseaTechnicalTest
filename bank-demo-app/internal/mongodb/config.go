package mongodb

import (
	"flag"
)

type MongoConfig struct {
	Host   string
	Port   string
	DbName string
	User   string
	Passwd string
}

func (cfg *MongoConfig) AddFlagsParams() {
	flag.StringVar(&cfg.Host, "mongo-host", "localhost", "MongoDB server host.")
	flag.StringVar(&cfg.Port, "mongo-port", "27017", "MongoDB server port.")
	flag.StringVar(&cfg.DbName, "mongo-db", "", "MongoDB database name (MONGODB_DATABASE).")
	flag.StringVar(&cfg.User, "mongo-user", "", "MongoDB username (MONGODB_USER).")
	flag.StringVar(&cfg.Passwd, "mongo-passwd", "", "MongoDB password (MONGODB_PASSWD).")
}

// Returns a url with the necessary format to connect to MongoDB.
func (cfg *MongoConfig) GetURL() string {
	if cfg.User != "" && cfg.Passwd != "" {
		return "mongodb://" + cfg.User + ":" + cfg.Passwd + "@" + cfg.Host + ":" + cfg.Port + "/" + cfg.DbName
	}
	return "mongodb://" + cfg.Host + ":" + cfg.Port + "/" + cfg.DbName
}
