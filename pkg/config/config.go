package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	Postgres PostgresConfig
	Web      WebConfig
}

type PostgresConfig struct {
	Host        string
	Port        int
	DBName      string
	User        string
	Password    string
	SSLMode     string
	ConnTimeout int
}

type WebConfig struct {
	Host string
	Port int
}

func (p PostgresConfig) ConnectionURL() (string, error) {
	host := p.Host
	v := p.Port
	if v < 1 && v > 65536 {
		return "", fmt.Errorf("PSQL_PORT invalid")
	}
	host = host + ":" + strconv.Itoa(p.Port)

	u := &url.URL{
		Scheme: "postgres",
		Host:   host,
		Path:   p.DBName,
	}

	if p.User == "" || p.Password == "" {
		return "", fmt.Errorf("PSQL_USER or PSQL_PASSWORD invalid")
	}
	u.User = url.UserPassword(p.User, p.Password)

	q := u.Query()
	connTimeout := p.ConnTimeout
	if connTimeout < 1 {
		return "", fmt.Errorf("PSQL_CONN_TIMEOUT invalid")
	}
	q.Add("connect_timeout", strconv.Itoa(p.ConnTimeout))

	if p.SSLMode != "disable" && p.SSLMode != "enable" {
		return "", fmt.Errorf("PSQL_SSL_MODE invalid")
	}
	q.Add("sslmode", p.SSLMode)

	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (w WebConfig) ConnectionURL() string {
	return fmt.Sprintf("%s:%d", w.Host, w.Port)
}

func LoadConfig(filenames ...string) error {
	return godotenv.Load(filenames...)
}

func New() Config {
	return Config{Postgres: PostgresConfig{
		Host:        getEnvAsString("PSQL_HOST", "localhost"),
		Port:        getEnvAsInt("PSQL_PORT", 5432),
		DBName:      getEnvAsString("PSQL_DB_NAME", "postgres"),
		User:        getEnvAsString("PSQL_USER", "postgres"),
		Password:    getEnvAsString("PSQL_PASSWORD", "postgres"),
		SSLMode:     getEnvAsString("PSQL_SSL_MODE", "disable"),
		ConnTimeout: getEnvAsInt("PSQL_CONN_TIMEOUT", 60),
	},
		Web: WebConfig{
			Host: getEnvAsString("WEB_HOST", "localhost"),
			Port: getEnvAsInt("WEB_PORT", 80),
		}}
}

func getEnvAsString(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnvAsString(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultValue
}
