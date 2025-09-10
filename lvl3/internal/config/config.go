package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Config содержит конфигурацию приложения
type Config struct {
	HTTP     HTTPConfig     `mapstructure:"http"`
	DBConfig DBConfig       `mapstructure:"postgres"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Telegram TelegramConfig `mapstructure:"telegram"`
	Email    EmailConfig    `mapstructure:"email"`
	Worker   WorkerConfig   `mapstructure:"worker"`
	Retry    RetryConfig    `mapstructure:"retry"`
}

// HTTPConfig содержит конфигурацию HTTP сервера
type HTTPConfig struct {
	Port         int           `mapstructure:"port" envconfig:"HTTP_SERVER_PORT" default:"8080"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" envconfig:"HTTP_READ_TIMEOUT" default:"15s"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" envconfig:"HTTP_WRITE_TIMEOUT" default:"15s"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout" envconfig:"HTTP_IDLE_TIMEOUT" default:"60s"`
}

// DBConfig содержит конфигурацию базы данных
type DBConfig struct {
	Host            string        `mapstructure:"host" envconfig:"POSTGRES_HOST" default:"localhost"`
	Port            int           `mapstructure:"port" envconfig:"POSTGRES_PORT" default:"5432"`
	Username        string        `mapstructure:"username" envconfig:"POSTGRES_USERNAME" default:"postgres"`
	Password        string        `mapstructure:"password" envconfig:"POSTGRES_PASSWORD" default:"postgres"`
	Database        string        `mapstructure:"database" envconfig:"POSTGRES_DB" default:"notifications"`
	SSLMode         string        `mapstructure:"ssl_mode" envconfig:"POSTGRES_SSLMODE" default:"disable"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" envconfig:"POSTGRES_MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" envconfig:"POSTGRES_MAX_IDLE_CONNS" default:"5"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" envconfig:"POSTGRES_CONN_MAX_LIFETIME" default:"5m"`
}

// RabbitMQConfig содержит конфигурацию RabbitMQ
type RabbitMQConfig struct {
	URL        string        `mapstructure:"url" envconfig:"RABBITMQ_URL" default:"amqp://guest:guest@localhost:5672/"`
	QueueName  string        `mapstructure:"queue_name" envconfig:"RABBITMQ_QUEUE" default:"notifications"`
	Exchange   string        `mapstructure:"exchange" envconfig:"RABBITMQ_EXCHANGE" default:"notifications_exchange"`
	MaxRetries int           `mapstructure:"max_retries" envconfig:"RABBITMQ_MAX_RETRIES" default:"5"`
	RetryDelay time.Duration `mapstructure:"retry_delay" envconfig:"RABBITMQ_RETRY_DELAY" default:"5s"`
}

// RedisConfig содержит конфигурацию Redis
type RedisConfig struct {
	URL             string        `mapstructure:"url" envconfig:"REDIS_URL" default:"localhost:6379"`
	Password        string        `mapstructure:"password" envconfig:"REDIS_PASSWORD" default:""`
	DB              int           `mapstructure:"db" envconfig:"REDIS_DB" default:"0"`
	Timeout         time.Duration `mapstructure:"timeout" envconfig:"REDIS_TIMEOUT" default:"5s"`
	NotificationTTL time.Duration `mapstructure:"notification_ttl" envconfig:"REDIS_NOTIFICATION_TTL" default:"24h"`
	PoolSize        int           `mapstructure:"pool_size" envconfig:"REDIS_POOL_SIZE" default:"10"`
	MinIdleConns    int           `mapstructure:"min_idle_conns" envconfig:"REDIS_MIN_IDLE_CONNS" default:"5"`
	MaxConnAge      time.Duration `mapstructure:"max_conn_age" envconfig:"REDIS_MAX_CONN_AGE" default:"30m"`
	PoolTimeout     time.Duration `mapstructure:"pool_timeout" envconfig:"REDIS_POOL_TIMEOUT" default:"30s"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout" envconfig:"REDIS_IDLE_TIMEOUT" default:"5m"`
	IdleCheckFreq   time.Duration `mapstructure:"idle_check_freq" envconfig:"REDIS_IDLE_CHECK_FREQ" default:"1m"`
}

// LoggingConfig содержит конфигурацию логирования
type LoggingConfig struct {
	Level  string `mapstructure:"level" envconfig:"LOGGING_LEVEL" default:"info"`
	Format string `mapstructure:"format" envconfig:"LOGGING_FORMAT" default:"json"`
}

// TelegramConfig содержит конфигурацию Telegram
type TelegramConfig struct {
	BotToken string `mapstructure:"bot_token" envconfig:"TELEGRAM_BOT_TOKEN" default:""`
	ChatID   int64  `mapstructure:"chat_id" envconfig:"TELEGRAM_CHAT_ID" default:"0"`
}

// EmailConfig содержит конфигурацию email
type EmailConfig struct {
	SMTPHost  string `mapstructure:"smtp_host" envconfig:"EMAIL_SMTP_HOST" default:"smtp.gmail.com"`
	SMTPPort  string `mapstructure:"smtp_port" envconfig:"EMAIL_SMTP_PORT" default:"587"`
	Username  string `mapstructure:"username" envconfig:"EMAIL_USERNAME" default:""`
	Password  string `mapstructure:"password" envconfig:"EMAIL_PASSWORD" default:""`
	FromEmail string `mapstructure:"from_email" envconfig:"EMAIL_FROM_EMAIL" default:""`
	FromName  string `mapstructure:"from_name" envconfig:"EMAIL_FROM_NAME" default:"Notification Service"`
}

// WorkerConfig содержит конфигурацию воркеров
type WorkerConfig struct {
	Count          int           `mapstructure:"count" envconfig:"WORKER_COUNT" default:"3"`
	ProcessTimeout time.Duration `mapstructure:"process_timeout" envconfig:"WORKER_PROCESS_TIMEOUT" default:"30s"`
}

// RetryConfig содержит конфигурацию повторных попыток
type RetryConfig struct {
	PublisherAttempts int           `mapstructure:"publisher_attempts" envconfig:"RETRY_PUBLISHER_ATTEMPTS" default:"3"`
	PublisherDelay    time.Duration `mapstructure:"publisher_delay" envconfig:"RETRY_PUBLISHER_DELAY" default:"1s"`
	PublisherBackoff  float64       `mapstructure:"publisher_backoff" envconfig:"RETRY_PUBLISHER_BACKOFF" default:"2"`
	ConsumerAttempts  int           `mapstructure:"consumer_attempts" envconfig:"RETRY_CONSUMER_ATTEMPTS" default:"3"`
	ConsumerDelay     time.Duration `mapstructure:"consumer_delay" envconfig:"RETRY_CONSUMER_DELAY" default:"1s"`
	ConsumerBackoff   float64       `mapstructure:"consumer_backoff" envconfig:"RETRY_CONSUMER_BACKOFF" default:"2"`
	MaxRetries        int           `mapstructure:"max_retries" envconfig:"RETRY_MAX_RETRIES" default:"3"`
}

// LoadConfig загружает конфигурацию из файла и переменных окружения
func LoadConfig() (*Config, error) {
	var cfg Config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err == nil {
		log.Info().Str("config_file", viper.ConfigFileUsed()).Msg("Using configuration file")
		if err := viper.Unmarshal(&cfg); err != nil {
			return nil, err
		}
	} else {
		log.Info().Msg("YAML config not found, using defaults and env")
	}

	if err := godotenv.Load(); err != nil {
		log.Info().Msg(".env file not found, using system environment variables")
	}
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

// GetDSN возвращает строку подключения к базе данных
func (d *DBConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.Username, d.Password, d.Database, d.SSLMode)
}

// Validate валидирует конфигурацию
func (c *Config) Validate() error {
	if c.HTTP.Port <= 0 || c.HTTP.Port > 65535 {
		return fmt.Errorf("invalid HTTP port: %d", c.HTTP.Port)
	}
	if c.RabbitMQ.URL == "" {
		return fmt.Errorf("RabbitMQ URL is required")
	}
	if c.Redis.URL == "" {
		return fmt.Errorf("redis URL is required")
	}
	if c.RabbitMQ.MaxRetries < 0 {
		return fmt.Errorf("RabbitMQ MaxRetries must be non-negative")
	}
	if c.Redis.NotificationTTL <= 0 {
		return fmt.Errorf("redis NotificationTTL must be positive")
	}

	return nil
}
