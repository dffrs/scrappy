package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const envName = "scrappy.env"

type config struct {
	from     string
	to       []string
	password string
	host     string
	port     int
}

func LoadEnv() (*config, error) {
	path, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	envPath := fmt.Sprintf("%s%s%s", path, string(os.PathSeparator), envName)

	// TODO: if file does not exist, create it
	err = godotenv.Load(envPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load %s file. Create it, if it does not exist", envPath)
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return nil, err
	}

	to := strings.Split(os.Getenv("TO"), ",")

	return &config{
		from:     os.Getenv("FROM"),
		to:       to,
		password: os.Getenv("PASSWORD"),
		host:     os.Getenv("HOST"),
		port:     port,
	}, nil
}
