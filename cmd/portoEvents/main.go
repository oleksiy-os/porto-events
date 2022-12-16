package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/oleksiy-os/porto-events/configs"
	"github.com/oleksiy-os/porto-events/internal/store"
	"github.com/oleksiy-os/porto-events/internal/store/boltdb"
	"github.com/oleksiy-os/porto-events/internal/web"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func main() {
	config := configInit()

	var s store.StoreInterface = boltdb.New()

	srv := web.New(config, &s)
	srv.ListenAndServe()
}

func configInit() *configs.Config {
	var (
		config     *configs.Config
		configPath string
		logLevel   string
	)

	flag.StringVar(&configPath, "config-path", "configs/config.toml", "path to config file")
	flag.StringVar(&logLevel, "log-level", "", "log level, int:0-6 (panic=0, fatal=1, error=2, warn=3, info=4, debug=5, trace=6)")
	flag.Parse()

	_, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		log.Fatal(err)
	}

	log.SetLevel(log.Level(config.LogLevel))

	if logLevel != "" {
		lvl, err := strconv.Atoi(logLevel)
		if err != nil && lvl >= 0 && lvl < 7 {
			log.SetLevel(log.Level(lvl))
			log.Println("changed log level:", log.GetLevel())
		}
	}

	if log.GetLevel() > 4 {
		formatter := &log.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		}
		log.SetFormatter(formatter)
	}

	if log.GetLevel() == 6 {
		log.SetReportCaller(true)
	}

	log.WithField("logLevel", log.GetLevel()).Debugln("App info:")

	return config
}
