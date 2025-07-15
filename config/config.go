package config

import (
	"github.com/joho/godotenv"
)

// LoadConfig carrega as variáveis de ambiente do .env
func LoadConfig() error {
	return godotenv.Load()
}
