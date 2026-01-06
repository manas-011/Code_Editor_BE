package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port string
	ReadTimeout int
	WriteTimeout int
}

func Load() *Config {

	readTimeout, err := strconv.Atoi(getEnv("READ_TIMEOUT", "5"))
	if err != nil {
		log.Fatal(err)
	}

	writeTimeout, err := strconv.Atoi(getEnv("WRITE_TIMEOUT", "10"))
	if err != nil {
		log.Fatal(err)
	}

	return &Config {
		Port: getEnv("PORT", "8080"),
		ReadTimeout: readTimeout,
		WriteTimeout: writeTimeout,
	}

}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}