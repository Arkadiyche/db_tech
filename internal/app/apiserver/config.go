package apiserver

import "githun.com/Arkadiyche/bd_techpark/internal/pkg/store"

type Config struct {
	BindAddr string `toml:"bind_addr"`
	Store    *store.Config
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8000",
		Store:    store.NewConfig(),
	}
}