package runtime

import (
	"flag"
	"log"

	"github.com/Unknwon/goconfig"
)

var (
	Port       int
	Migrate    bool
	Debug      bool
	LokiServer string
	Cfg        *goconfig.ConfigFile
)

func init() {
	var err error
	Cfg, err = goconfig.LoadConfigFile("backend/dagger.ini")
	if err != nil {
		Cfg, err = goconfig.LoadConfigFile("/etc/dagger/dagger.ini")
		if err != nil {
			log.Panicf("loading setting conf: dagger.ini fail %s", err)
		}
	}

	flag.IntVar(&Port, "port", 8000, "port")
	flag.BoolVar(&Migrate, "migrate", true, "migrate db")
	flag.BoolVar(&Debug, "debug", false, "debug mode")
	flag.StringVar(&LokiServer, "loki-server", "", "loki server address, ex: http://127.0.0.1:3100")
	flag.Parse()
}
