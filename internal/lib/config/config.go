package config

import (
	"go-admin/internal/lib/env"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Settings = &Config{}

type Config struct {
	AdminServer *AdminServer `mapstructure:"admin-server"`
	MySQL       *MySQL       `mapstructure:"mysql"`
	Redis       *Redis       `mapstructure:"redis"`
}

// AdminServer
type AdminServer struct {
	Server ServerConfig `mapstructure:"server"`
	Log    LogConfig    `mapstructure:"log"`
}

type ServerConfig struct {
	GinMode      string        `mapstructure:"gin_mode"`
	HttpPort     int           `mapstructure:"http_port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	JwtSecret    string        `mapstructure:"jwt_secret"`
}

type LogConfig struct {
	LogSavePath string `mapstructure:"log_save_path"`
	LogSaveName string `mapstructure:"log_save_name"`
	LogFileExt  string `mapstructure:"log_file_ext"`
}

// MySQL
type MySQL struct {
	GoAdmin MySQLConfig `mapstructure:"go-admin"`
}

type MySQLConfig struct {
	Conn        []string      `mapstructure:"conn"`
	MaxIdle     int           `mapstructure:"max_idle"`
	MaxOpen     int           `mapstructure:"max_open"`
	MaxLifetime time.Duration `mapstructure:"max_lifetime"`
}

// Redis
type Redis struct {
	GoAdmin RedisConfig `mapstructure:"go-admin"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func Init() {
	var rootDir string
	currentPath, _ := os.Getwd()
	rootDir = currentPath
	configDir := filepath.Join(rootDir, "configs", env.Mode())

	files, err := os.ReadDir(configDir)
	if err != nil {
		log.Fatalf("read config files name error: %v", err)
	}

	viperInstances := map[string]*viper.Viper{}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".yaml" {
			fileName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			v := viper.New()
			v.AddConfigPath(configDir)
			v.SetConfigName(fileName)
			v.SetConfigType("yaml")

			if err := v.ReadInConfig(); err != nil {
				log.Fatalf("Error reading config file %s: %v", file.Name(), err)
			}

			viperInstances[fileName] = v

			if err := v.Unmarshal(Settings); err != nil {
				log.Fatalf("Error unmarshaling config file %s: %v", file.Name(), err)
			}

			v.WatchConfig()
			v.OnConfigChange(func(e fsnotify.Event) {
				// fmt.Printf("Config file changed: %s\n", e.Name)
				if err := v.Unmarshal(Settings); err != nil {
					log.Printf("Error unmarshaling updated config file %s: %v", file.Name(), err)
					return
				}
				// fmt.Printf("Updated config [%s]: %+v\n", fileName, Settings)
			})
		}
	}

	// fmt.Println("all_settings", viperConfig.AllSettings())
}
