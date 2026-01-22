package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Host != "localhost" {
		t.Errorf("expected host localhost, got %s", cfg.Host)
	}
	if cfg.Port != 8423 {
		t.Errorf("expected port 8423, got %d", cfg.Port)
	}
}

func TestBaseURL(t *testing.T) {
	cfg := Config{Host: "192.168.1.5", Port: 9000}
	expected := "http://192.168.1.5:9000"
	if cfg.BaseURL() != expected {
		t.Errorf("expected %s, got %s", expected, cfg.BaseURL())
	}
}

func TestBaseURLDefault(t *testing.T) {
	cfg := DefaultConfig()
	expected := "http://localhost:8423"
	if cfg.BaseURL() != expected {
		t.Errorf("expected %s, got %s", expected, cfg.BaseURL())
	}
}

func TestLoadMissingFile(t *testing.T) {
	// Override config path to a non-existent file
	original := os.Getenv("HOME")
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	defer os.Setenv("HOME", original)

	cfg := Load()
	if cfg.Host != "localhost" || cfg.Port != 8423 {
		t.Errorf("expected defaults, got host=%s port=%d", cfg.Host, cfg.Port)
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cfg := Config{Host: "10.0.0.1", Port: 9999}
	if err := Save(cfg); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded := Load()
	if loaded.Host != "10.0.0.1" {
		t.Errorf("expected host 10.0.0.1, got %s", loaded.Host)
	}
	if loaded.Port != 9999 {
		t.Errorf("expected port 9999, got %d", loaded.Port)
	}
}

func TestLoadDefaultsForMissingFields(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	dir := filepath.Join(tmp, ".config", "itsyhome")
	os.MkdirAll(dir, 0755)
	os.WriteFile(filepath.Join(dir, "config.json"), []byte(`{}`), 0644)

	cfg := Load()
	if cfg.Host != "localhost" {
		t.Errorf("expected default host, got %s", cfg.Host)
	}
	if cfg.Port != 8423 {
		t.Errorf("expected default port, got %d", cfg.Port)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	dir := filepath.Join(tmp, ".config", "itsyhome")
	os.MkdirAll(dir, 0755)
	os.WriteFile(filepath.Join(dir, "config.json"), []byte(`not json`), 0644)

	cfg := Load()
	if cfg.Host != "localhost" || cfg.Port != 8423 {
		t.Errorf("expected defaults for invalid JSON, got host=%s port=%d", cfg.Host, cfg.Port)
	}
}

func TestPath(t *testing.T) {
	p := Path()
	if p == "" {
		t.Skip("cannot determine home dir")
	}
	if filepath.Base(p) != "config.json" {
		t.Errorf("expected config.json, got %s", filepath.Base(p))
	}
}
