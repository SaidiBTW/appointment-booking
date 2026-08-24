package config

import "testing"

func TestLoadUsesDefaultSSLModeWhenNotProvided(t *testing.T) {
	t.Setenv("DATABASE_HOST", "db.internal")
	t.Setenv("DATABASE_PORT", "5432")
	t.Setenv("DATABASE_USER", "app")
	t.Setenv("DATABASE_PASSWORD", "secret")
	t.Setenv("DATABASE_NAME", "appointments")
	t.Setenv("DATABASE_SSLMODE", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.SSLMode != DefaultPostgresConfig().SSLMode {
		t.Fatalf("SSLMode = %q, want %q", cfg.SSLMode, DefaultPostgresConfig().SSLMode)
	}

	if cfg.Host != "db.internal" || cfg.Port != 5432 || cfg.User != "app" || cfg.Password != "secret" || cfg.DBName != "appointments" {
		t.Fatalf("Load() returned unexpected config: %+v", cfg)
	}
}

func TestLoadFallsBackToDefaultConfigWhenDBNameMissing(t *testing.T) {
	t.Setenv("DATABASE_HOST", "db.internal")
	t.Setenv("DATABASE_PORT", "5432")
	t.Setenv("DATABASE_USER", "app")
	t.Setenv("DATABASE_PASSWORD", "secret")
	t.Setenv("DATABASE_NAME", "")
	t.Setenv("DATABASE_SSLMODE", "require")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	defaultCfg := DefaultPostgresConfig()
	if *cfg != *defaultCfg {
		t.Fatalf("Load() = %+v, want %+v", cfg, defaultCfg)
	}
}
