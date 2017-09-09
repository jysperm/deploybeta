package config

import "os"

var Port = ":7000"

var EtcdEndpoints = []string{"http://127.0.0.1:2379"}

var DefaultRegistry string

func init() {
	if os.Getenv("DEFAULT_REGISTRY") != "" {
		DefaultRegistry = os.Getenv("DEFAULT_REGISTRY")
	} else {
		DefaultRegistry = "localhost:5000"
	}
}
