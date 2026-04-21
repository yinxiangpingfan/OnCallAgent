package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitConfig(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.json")

	configContent := `{
		"server": {
			"host": "testhost",
			"port": 9999
		},
		"embedder": {
			"host": "embedder-host",
			"port": 11434,
			"model": "test-model",
			"dimension": 768
		},
		"qdrant": {
			"host": "qdrant-host",
			"port": 6334,
			"collection": "test-collection"
		}
	}`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	cfg, err := InitConfig(configFile)
	if err != nil {
		t.Fatalf("InitConfig failed: %v", err)
	}

	// 验证 Server 配置
	if cfg.Server.Host != "testhost" {
		t.Errorf("Expected Server.Host=testhost, got %s", cfg.Server.Host)
	}
	if cfg.Server.Port != 9999 {
		t.Errorf("Expected Server.Port=9999, got %d", cfg.Server.Port)
	}

	// 验证 Embedder 配置
	if cfg.Embedder.Host != "embedder-host" {
		t.Errorf("Expected Embedder.Host=embedder-host, got %s", cfg.Embedder.Host)
	}
	if cfg.Embedder.Model != "test-model" {
		t.Errorf("Expected Embedder.Model=test-model, got %s", cfg.Embedder.Model)
	}
	if cfg.Embedder.Dimension != 768 {
		t.Errorf("Expected Embedder.Dimension=768, got %d", cfg.Embedder.Dimension)
	}

	// 验证 Qdrant 配置
	if cfg.Qdrant.Host != "qdrant-host" {
		t.Errorf("Expected Qdrant.Host=qdrant-host, got %s", cfg.Qdrant.Host)
	}
	if cfg.Qdrant.Collection != "test-collection" {
		t.Errorf("Expected Qdrant.Collection=test-collection, got %s", cfg.Qdrant.Collection)
	}
}

func TestInitConfigFromFileNotFound(t *testing.T) {
	_, err := InitConfig("nonexistent.json")
	if err == nil {
		t.Error("Expected error for nonexistent config file")
	}
}

func TestInitConfigFromEnv(t *testing.T) {
	// 设置环境变量
	envVars := map[string]string{
		"SERVER_HOST":        "env-server",
		"SERVER_PORT":        "8080",
		"EMBEDDER_HOST":      "env-embedder",
		"EMBEDDER_PORT":      "11435",
		"EMBEDDER_MODEL":     "env-model",
		"EMBEDDER_DIMENSION": "512",
		"QDRANT_HOST":        "env-qdrant",
		"QDRANT_PORT":        "6335",
		"QDRANT_COLLECTION":  "env-collection",
	}

	// 设置环境变量
	for key, value := range envVars {
		os.Setenv(key, value)
		defer os.Unsetenv(key)
	}

	cfg, err := InitConfigFromEnv()
	if err != nil {
		t.Fatalf("InitConfigFromEnv failed: %v", err)
	}

	// 验证配置从环境变量读取
	if cfg.Server.Host != "env-server" {
		t.Errorf("Expected Server.Host=env-server, got %s", cfg.Server.Host)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("Expected Server.Port=8080, got %d", cfg.Server.Port)
	}
	if cfg.Embedder.Model != "env-model" {
		t.Errorf("Expected Embedder.Model=env-model, got %s", cfg.Embedder.Model)
	}
	if cfg.Qdrant.Collection != "env-collection" {
		t.Errorf("Expected Qdrant.Collection=env-collection, got %s", cfg.Qdrant.Collection)
	}
}

func TestInitConfigFromEnvDefaults(t *testing.T) {
	// 清除可能存在的环境变量
	envKeys := []string{
		"SERVER_HOST", "SERVER_PORT",
		"EMBEDDER_HOST", "EMBEDDER_PORT", "EMBEDDER_MODEL", "EMBEDDER_DIMENSION",
		"QDRANT_HOST", "QDRANT_PORT", "QDRANT_COLLECTION",
	}
	for _, key := range envKeys {
		os.Unsetenv(key)
	}

	cfg, err := InitConfigFromEnv()
	if err != nil {
		t.Fatalf("InitConfigFromEnv failed: %v", err)
	}

	// 验证默认值
	if cfg.Server.Host != "localhost" {
		t.Errorf("Expected default Server.Host=localhost, got %s", cfg.Server.Host)
	}
	if cfg.Server.Port != 8819 {
		t.Errorf("Expected default Server.Port=8819, got %d", cfg.Server.Port)
	}
	if cfg.Embedder.Model != "nomic-embed-text" {
		t.Errorf("Expected default Embedder.Model=nomic-embed-text, got %s", cfg.Embedder.Model)
	}
	if cfg.Embedder.Dimension != 384 {
		t.Errorf("Expected default Embedder.Dimension=384, got %d", cfg.Embedder.Dimension)
	}
}

func TestInitConfigWithFallback(t *testing.T) {
	// 设置环境变量用于 fallback
	os.Setenv("SERVER_HOST", "fallback-server")
	os.Setenv("SERVER_PORT", "9090")
	defer func() {
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("SERVER_PORT")
	}()

	// 测试文件不存在时回退到环境变量
	cfg, err := InitConfigWithFallback("nonexistent.json")
	if err != nil {
		t.Fatalf("InitConfigWithFallback failed: %v", err)
	}

	if cfg.Server.Host != "fallback-server" {
		t.Errorf("Expected Server.Host=fallback-server, got %s", cfg.Server.Host)
	}
	if cfg.Server.Port != 9090 {
		t.Errorf("Expected Server.Port=9090, got %d", cfg.Server.Port)
	}
}

func TestInitConfigWithFallbackUsesFile(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.json")

	configContent := `{
		"server": {
			"host": "file-host",
			"port": 7777
		}
	}`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// 设置环境变量（不应该被使用）
	os.Setenv("SERVER_HOST", "env-host")
	os.Setenv("SERVER_PORT", "8888")
	defer func() {
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("SERVER_PORT")
	}()

	cfg, err := InitConfigWithFallback(configFile)
	if err != nil {
		t.Fatalf("InitConfigWithFallback failed: %v", err)
	}

	// 应该使用配置文件的值，而不是环境变量
	if cfg.Server.Host != "file-host" {
		t.Errorf("Expected Server.Host=file-host, got %s", cfg.Server.Host)
	}
	if cfg.Server.Port != 7777 {
		t.Errorf("Expected Server.Port=7777, got %d", cfg.Server.Port)
	}
}

func TestGetServerAddr(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Host: "localhost", Port: 8080},
	}
	expected := "localhost:8080"
	if addr := cfg.GetServerAddr(); addr != expected {
		t.Errorf("Expected GetServerAddr()=%s, got %s", expected, addr)
	}
}

func TestGetEmbedderAddr(t *testing.T) {
	cfg := &Config{
		Embedder: EmbedderConfig{Host: "embedder", Port: 11434},
	}
	expected := "embedder:11434"
	if addr := cfg.GetEmbedderAddr(); addr != expected {
		t.Errorf("Expected GetEmbedderAddr()=%s, got %s", expected, addr)
	}
}

func TestGetQdrantAddr(t *testing.T) {
	cfg := &Config{
		Qdrant: QdrantConfig{Host: "qdrant", Port: 6334},
	}
	expected := "qdrant:6334"
	if addr := cfg.GetQdrantAddr(); addr != expected {
		t.Errorf("Expected GetQdrantAddr()=%s, got %s", expected, addr)
	}
}
