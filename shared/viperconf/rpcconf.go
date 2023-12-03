package viperconf

import (
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

var (
	// ConfigLoaded signals if the config has been loaded or not
	ConfigLoaded bool
	ConfigMutex  sync.Mutex
	RPCViper     = viper.New()
)

func InitConfig() {
	ConfigMutex.Lock()
	defer ConfigMutex.Unlock()
	if ConfigLoaded {
		return
	}
	ConfigLoaded = true

	RPCViper.SetConfigName("rpc")
	RPCViper.AddConfigPath("/opt")
	RPCViper.AddConfigPath(".")
	confErr := RPCViper.ReadInConfig()
	if confErr != nil {
		if _, castok := confErr.(viper.ConfigFileNotFoundError); !castok {
			log.Printf("ERROR Reading config: %v", confErr)
			if os.Getenv("CONFIG_OPTIONAL") == "" {
				os.Exit(1)
			}
		}
	}
}
