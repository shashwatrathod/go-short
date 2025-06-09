package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type DBConfig struct {
	Host     string
	Port     int
	DBName   string
	DBUser   string
	Password string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
}

type Config struct {
	DBConfigs   []DBConfig
	RedisConfig RedisConfig
}

// Load reads database configuration from environment variables and returns a Config instance.
// It parses comma-separated values from multiple environment variables to configure database connections.
//
// The function will panic if:
//   - Environment variables contain invalid data
//   - Configuration values are inconsistent or missing
//
// Returns a pointer to a Config struct containing the parsed database configurations.
func Load() (*Config, error) {
	dbConfigs, err := loadDBConfigs()

	if err != nil {
		return nil, err
	}

	redisConfig, err := loadRedisConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		DBConfigs:   dbConfigs,
		RedisConfig: *redisConfig,
	}, nil
}

func loadRedisConfig() (*RedisConfig, error) {
	host := strings.TrimSpace(os.Getenv("REDIS_HOST"))
	portString := strings.TrimSpace(os.Getenv("REDIS_PORT"))
	password := strings.TrimSpace(os.Getenv("REDIS_PASSWORD"))

	if host == "" {
		host = "localhost"
	}

	if portString == "" {
		portString = "6379"
	}

	portInt, err := strconv.Atoi(portString)
	if err != nil {
		return nil, err
	}

	return &RedisConfig{
		Host:     host,
		Port:     portInt,
		Password: password,
	}, nil
}

func loadDBConfigs() ([]DBConfig, error) {
	envHosts := os.Getenv("DB_HOST_LIST")
	envPorts := os.Getenv("DB_PORT_LIST")
	envNames := os.Getenv("DB_NAME_LIST")
	envUsers := os.Getenv("DB_USER_LIST")
	envPasswords := os.Getenv("DB_PASSWORD_LIST")

	hosts := strings.Split(envHosts, ",")
	portStrings := strings.Split(envPorts, ",")
	ports := make([]int, len(portStrings))
	for i, port := range portStrings {
		portInt, err := strconv.Atoi(strings.TrimSpace(port))
		if err != nil {
			return make([]DBConfig, 0), err
		}
		ports[i] = portInt
	}

	names := strings.Split(envNames, ",")
	users := strings.Split(envUsers, ",")
	passwords := strings.Split(envPasswords, ",")

	if len(hosts) != len(ports) || len(hosts) != len(names) || len(hosts) != len(users) || len(hosts) != len(passwords) {
		return make([]DBConfig, 0), fmt.Errorf("Environment variables for database configuration are not consistent in length. Please ensure DB_HOST_LIST, DB_PORT_LIST, DB_NAME_LIST, DB_USER_LIST, and DB_PASSWORD_LIST are set correctly.")
	}

	dbConfigs := make([]DBConfig, len(hosts))
	for i := range hosts {
		dbConfigs[i] = DBConfig{
			Host:     strings.TrimSpace(hosts[i]),
			Port:     ports[i],
			DBName:   strings.TrimSpace(names[i]),
			DBUser:   strings.TrimSpace(users[i]),
			Password: strings.TrimSpace(passwords[i]),
		}
	}

	return dbConfigs, nil
}
