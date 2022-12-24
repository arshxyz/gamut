package utils

import (
	"log"
	"os"
	"path"

	"github.com/spf13/viper"
)

type Secrets struct {
	ClientID     string `yaml:"clientid"`
	ClientSecret string `yaml:"clientsecret"`
	RefreshToken string `yaml:"refreshtoken"`
}

// Recursively generate config dir in $HOME/.config/gamut
// and the current directory
func createConfigDir() (cfgDirPath string, err error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}
	cfgDirPath = path.Join(homedir, ".config", "gamut")
	err = os.MkdirAll(cfgDirPath, os.ModePerm)
	return cfgDirPath, err
}

// Set Viper config defaults
func InitViper() error {
	cfgDirPath, err := createConfigDir()
	if err != nil {
		return err
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(cfgDirPath)
	viper.AddConfigPath(".")
	return nil
}

// Fetch secrets from config paths
func GetSecrets() (*Secrets, error) {
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	conf := &Secrets{}

	err := viper.Unmarshal(conf)
	return conf, err

}

// Write secrets to disk.
func WriteSecrets(sc Secrets) (err error) {
	viper.Set("ClientID", sc.ClientID)
	viper.Set("ClientSecret", sc.ClientSecret)
	viper.Set("RefreshToken", sc.RefreshToken)
	// Creates a new file if one does not exist
	viper.SafeWriteConfig()
	// Write to a file if it exists, else returns error
	err = viper.WriteConfig()
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return
}
