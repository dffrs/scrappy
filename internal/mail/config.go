package mail

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const envName = "scrappy.env"

type Config struct {
	From     string
	To       []string
	Password string
	Host     string
	Port     int
}

func createEmptyEnvFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	builder := strings.Builder{}

	builder.WriteString("FROM=\n")
	builder.WriteString("TO=\n")
	builder.WriteString("PASSWORD=\n")
	builder.WriteString("HOST=\n")
	builder.WriteString("PORT=\n")

	_, err = file.Write([]byte(builder.String()))
	if err != nil {
		return err
	}

	return nil
}

func isEnvFileCreated(filePath string) bool {
	if _, err := os.Open(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func loadConfig() (*Config, error) {
	path, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	envPath := fmt.Sprintf("%s%s%s", path, string(os.PathSeparator), envName)

	if !isEnvFileCreated(envPath) {
		fmt.Println("env file not found. Creating empty template...")
		err := createEmptyEnvFile(envPath)
		if err != nil {
			return nil, err
		}

		fmt.Printf("empty template created at %s\n", envPath)
	}

	err = godotenv.Load(envPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load %s file", envPath)
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return nil, err
	}

	to := strings.Split(os.Getenv("TO"), ",")

	return &Config{
		From:     os.Getenv("FROM"),
		To:       to,
		Password: os.Getenv("PASSWORD"),
		Host:     os.Getenv("HOST"),
		Port:     port,
	}, nil
}
