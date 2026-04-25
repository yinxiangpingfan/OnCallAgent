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
