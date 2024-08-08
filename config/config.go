package config

import (
	"flag"
	"fmt"
	"github.com/lunarnuts/inboxsuite-test/pkg/errors"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
	"strings"
)

type (
	Config struct {
		Logger   Logger   `yaml:"logger"`
		RabbitMQ RabbitMQ `yaml:"rabbitmq"`
		DB       DB       `yaml:"db"`
		Worker   int
	}

	Logger struct {
		Level string `yaml:"level"`
		Env   string `yaml:"env"`
	}

	RabbitMQ struct {
		JobQueue           string `yaml:"job_queue"`
		ResultExchange     string `yaml:"result_exchange"`
		StatisticsExchange string `yaml:"statistics_exchange"`
		Host               string `yaml:"host"`
		Port               string `yaml:"port"`
		User               string `yaml:"user"`
		Password           string `yaml:"password"`
	}

	DB struct {
		Name               string `yaml:"name"`
		Host               string `yaml:"host"`
		Port               string `yaml:"port"`
		User               string `yaml:"user"`
		Password           string `yaml:"password"`
		SSL                string `yaml:"ssl"`
		MaxIdleConnections int    `yaml:"maxIdleConnections" default:"10"`
		MaxOpenConnections int    `yaml:"maxOpenConnections" default:"2"`
		LogLevel           string `yaml:"logLevel" default:"info"`
	}
)

func (db *DB) ParseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		db.User, db.Password, db.Host, db.Port, db.Name, db.SSL)
}

func (r *RabbitMQ) ParseURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s", r.User, r.Password, r.Host, r.Port)
}

func Load() (*Config, error) {
	configFile := flag.String("config", "config.example.yaml", "path to config file")
	flag.Parse()

	data, err := preprocess(configFile)
	if err != nil {
		return nil, errors.Wrap(err, "error preprocessing config file")
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing config file")
	}

	n := flag.Int("worker", 10, "number of workers")
	if n != nil {
		cfg.Worker = *n
	}

	return &cfg, nil
}

// preprocess processes config file before parsing it. config file might have env vars, as not to directly expose sensitive
// parameters. It supports it, if env var in ${envVar} format is detected.
// Example:
// db:
//
//	name: ${POSTGRES_DB}
//	host: ${POSTGRES_HOST}
//	port: ${POSTGRES_PORT}
//	user: ${POSTGRES_USER}
//	password: ${POSTGRES_PASSWORD}
func preprocess(configFile *string) ([]byte, error) {
	data, err := os.ReadFile(*configFile)
	if err != nil {
		return nil, err
	}

	content := string(data)

	re := regexp.MustCompile(`\$\{(.+?)\}`)
	replacedContent := re.ReplaceAllStringFunc(content, func(s string) string {
		envVarName := strings.TrimSuffix(strings.TrimPrefix(s, `${`), `}`)
		envVarValue := os.Getenv(envVarName)
		return envVarValue
	})

	return []byte(replacedContent), err
}
