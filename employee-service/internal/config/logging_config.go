package config

import (
	"log"
)

func ConfigureLogging() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}
