package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Embedder EmbedderConfig `mapstructure:"embedder"`
	Qdrant   QdrantConfig   `mapstructure:"qdrant"`
	OpenAI   OpenAIConfig   `mapstructure:"openai"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// EmbedderConfig 嵌入模型配置
type EmbedderConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Model     string `mapstructure:"model"`
	Dimension int    `mapstructure:"dimension"`
}

// QdrantConfig Qdrant 向量数据库配置
type QdrantConfig struct {
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	Collection string `mapstructure:"collection"`
}

// OpenAIConfig OpenAI API 配置
type OpenAIConfig struct {
	APIKey  string `mapstructure:"api_key"`
	Model   string `mapstructure:"model"`
	APIBase string `mapstructure:"api_base"`
}

// InitConfig 从配置文件初始化配置
// configFile: 配置文件路径，如 "config/config.json"
func InitConfig(configFile string) (*Config, error) {
	v := viper.New()

	// 设置配置文件
	v.SetConfigFile(configFile)

	// 设置配置文件类型
	if strings.HasSuffix(configFile, ".json") {
		v.SetConfigType("json")
	} else if strings.HasSuffix(configFile, ".yaml") || strings.HasSuffix(configFile, ".yml") {
		v.SetConfigType("yaml")
	}

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 设置默认值
	setDefaults(v)

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// InitConfigFromEnv 从环境变量初始化配置
// 环境变量命名规则：SERVER_HOST, SERVER_PORT, EMBEDDER_HOST 等
func InitConfigFromEnv() (*Config, error) {
	v := viper.New()

	// 设置默认值
	setDefaults(v)

	// 绑定环境变量
	bindEnvVars(v)

	// 自动绑定环境变量
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from env: %w", err)
	}

	return &cfg, nil
}

// InitConfigWithFallback 优先使用配置文件，失败则回退到环境变量
// configFile: 配置文件路径
// 如果配置文件不存在或读取失败，将尝试从环境变量读取
func InitConfigWithFallback(configFile string) (*Config, error) {
	// 先尝试从配置文件读取
	cfg, err := InitConfig(configFile)
	if err == nil {
		return cfg, nil
	}

	// 配置文件读取失败，尝试从环境变量读取
	fmt.Printf("Config file read failed: %v, falling back to environment variables\n", err)
	return InitConfigFromEnv()
}

// setDefaults 设置默认值
func setDefaults(v *viper.Viper) {
	// Server 默认值
	v.SetDefault("server.host", "localhost")
	v.SetDefault("server.port", 8819)

	// Embedder 默认值
	v.SetDefault("embedder.host", "localhost")
	v.SetDefault("embedder.port", 11434)
	v.SetDefault("embedder.model", "nomic-embed-text")
	v.SetDefault("embedder.dimension", 384)

	// Qdrant 默认值
	v.SetDefault("qdrant.host", "localhost")
	v.SetDefault("qdrant.port", 6334)
	v.SetDefault("qdrant.collection", "oncallagent")

	// OpenAI 默认值
	v.SetDefault("openai.api_key", "")
	v.SetDefault("openai.model", "minimax/minimax-m2.1")
	v.SetDefault("openai.api_base", "https://api.qnaigc.com/v1")
}

// bindEnvVars 绑定环境变量
// 环境变量格式：SERVER_HOST, SERVER_PORT, EMBEDDER_HOST 等
func bindEnvVars(v *viper.Viper) {
	// Server 环境变量
	_ = v.BindEnv("server.host", "SERVER_HOST")
	_ = v.BindEnv("server.port", "SERVER_PORT")

	// Embedder 环境变量
	_ = v.BindEnv("embedder.host", "EMBEDDER_HOST")
	_ = v.BindEnv("embedder.port", "EMBEDDER_PORT")
	_ = v.BindEnv("embedder.model", "EMBEDDER_MODEL")
	_ = v.BindEnv("embedder.dimension", "EMBEDDER_DIMENSION")

	// Qdrant 环境变量
	_ = v.BindEnv("qdrant.host", "QDRANT_HOST")
	_ = v.BindEnv("qdrant.port", "QDRANT_PORT")
	_ = v.BindEnv("qdrant.collection", "QDRANT_COLLECTION")

	// OpenAI 环境变量
	_ = v.BindEnv("openai.api_key", "OPENAI_API_KEY")
	_ = v.BindEnv("openai.model", "OPENAI_MODEL")
	_ = v.BindEnv("openai.api_base", "OPENAI_API_BASE")
}

// GetServerAddr 获取服务器完整地址
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// GetEmbedderAddr 获取嵌入模型服务地址
func (c *Config) GetEmbedderAddr() string {
	return fmt.Sprintf("%s:%d", c.Embedder.Host, c.Embedder.Port)
}

// GetQdrantAddr 获取 Qdrant 服务地址
func (c *Config) GetQdrantAddr() string {
	return fmt.Sprintf("%s:%d", c.Qdrant.Host, c.Qdrant.Port)
}
