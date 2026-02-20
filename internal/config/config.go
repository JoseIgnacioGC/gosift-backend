package config

import (
	"log"
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type GinConfig struct {
	TrustedProxies []string `env:"TRUSTED_PROXIES" envSeparator:"," envDefault:"*"`
	GinMode        string   `env:"GIN_MODE" envDefault:"debug"`
	Port           string   `env:"PORT" envDefault:"8080"`
}

var (
	once      sync.Once
	ginConfig *GinConfig
)

func Get() *GinConfig {
	once.Do(func() {
		if err := godotenv.Load(".env"); err != nil {
			log.Printf("[WARNING] No .env file found: %v\n", err)
		}

		ginConfig = &GinConfig{}

		if err := env.Parse(ginConfig); err != nil {
			log.Fatalf("[ERROR] Can't parse environment variables: %v", err)
		}
	})

	return ginConfig
}

func ConfigureProxies(router *gin.Engine, ginConfig *GinConfig) {
	if ginConfig.GinMode == "debug" || (len(ginConfig.TrustedProxies) > 0 && ginConfig.TrustedProxies[0] == "*") {
		log.Println("[DEBUG WARNING]: All proxies are trusted.")
		return
	}

	if len(ginConfig.TrustedProxies) == 0 {
		log.Fatalln("[ERROR]: No trusted proxies configured in release. Set TRUSTED_PROXIES env var.")
	}

	log.Printf("[INFO]: Configuring trusted proxies: %v\n", ginConfig.TrustedProxies)
	router.SetTrustedProxies(ginConfig.TrustedProxies)

}
