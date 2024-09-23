package storage

import (
	"github.com/mitchellh/go-homedir"
	"github.com/nutsdb/nutsdb"
	"github.com/rs/zerolog"
	"os"
	"time"
)

// CreateConfigDirectory creates the config directory
func CreateConfigDirectory() (configDir string, err error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	homePath, err := homedir.Expand(home)
	if err != nil {
		return "", err
	}
	configDir = homePath + "/.config/tinyRconConsole/storage"
	_, err = os.Stat(configDir)
	if err != nil {
		err = os.MkdirAll(configDir, 0744)
	}
	if err != nil {
		return "", err
	}
	return configDir, nil
}
func initDatabase(configDir string) (err error) {
	db, err := nutsdb.Open(
		nutsdb.DefaultOptions,
		nutsdb.WithDir(configDir),
	)
	defer db.Close()
	if err != nil {
		return err
	}
	return nil
}
func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	configDir, err := CreateConfigDirectory()
	if err != nil {
		log.Fatal().AnErr("error: ", err).Msgf("fail to create dir")
	}
	err = initDatabase(configDir)
	if err != nil {
		log.Fatal().AnErr("error: ", err).Msgf("fail to init database")
	}
}
