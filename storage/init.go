package storage

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/binary"
	"errors"
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
	if os.IsExist(err) {
		return configDir, errors.New("config directory already exists")
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(configDir, 0744)
	}
	if err != nil {
		return "", err
	}
	return configDir, nil
}
func initDatabase(DBPath string) (err error) {
	db, err := nutsdb.Open(
		nutsdb.DefaultOptions,
		nutsdb.WithDir(DBPath),
	)
	defer func(db *nutsdb.DB) {
		_ = db.Close()
	}(db)
	if err != nil {
		return err
	}
	buckets := make([]func(tx *nutsdb.Tx) error, 0)
	buckets = append(buckets, func(tx *nutsdb.Tx) error {
		return tx.NewBucket(nutsdb.DataStructureSet, "settings")
	})
	buckets = append(buckets, func(tx *nutsdb.Tx) error {
		return tx.NewBucket(nutsdb.DataStructureList, "servers")
	})
	buckets = append(buckets, func(tx *nutsdb.Tx) error {
		return tx.NewBucket(nutsdb.DataStructureBTree, "savedCommand")
	})
	for _, v := range buckets {
		if err = db.Update(v); err != nil {
			return errors.Join(err, errors.New("fail to create bucket"))
		}
	}
	tx, err := db.Begin(true)
	if err != nil {
		return errors.Join(err, errors.New("fail to begin tx"))
	}
	var binData []byte
	_, _ = binary.Encode(binData, binary.NativeEndian, "en-US")
	err = tx.Put("settings", []byte("language"), binData, 0)
	if err != nil {
		_ = tx.Rollback()
		return errors.Join(err, errors.New("fail to put"))
	}
	_, _ = binary.Encode(binData, binary.NativeEndian, "buttIn")
	err = tx.Put("settings", []byte("storageBackend"), binData, 0)
	if err != nil {
		_ = tx.Rollback()
		return errors.Join(err, errors.New("fail to put"))
	}
	_, _ = binary.Encode(binData, binary.NativeEndian, "SM2")
	err = tx.Put("settings", []byte("cryptoBackend"), binData, 0)
	if err != nil {
		_ = tx.Rollback()
		return errors.Join(err, errors.New("fail to put"))
	}
	_, _ = binary.Encode(binData, binary.NativeEndian, true)
	err = tx.Put("settings", []byte("tcpKeepLive"), binData, 0)
	if err != nil {
		_ = tx.Rollback()
		return errors.Join(err, errors.New("fail to put"))
	}
	_, _ = binary.Encode(binData, binary.NativeEndian, time.Now().In(time.Now().Location()))
	err = tx.Put("settings", []byte("timezone"), binData, 0)
	if err != nil {
		_ = tx.Rollback()
		return errors.Join(err, errors.New("fail to put"))
	}
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return errors.Join(err, errors.New("fail to generate key"))
	}
	_, _ = binary.Encode(binData, binary.NativeEndian, privateKey)
	err = tx.Put("settings", []byte("privateKey"), binData, 0)
	if err != nil {
		_ = tx.Rollback()
		return errors.Join(err, errors.New("fail to put"))
	}
	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return errors.Join(err, errors.New("fail to commit setup data"))
	}
	return nil
}
func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	configDir, err := CreateConfigDirectory()
	if !errors.Is(err, errors.New("config directory already exists")) || err != nil {
		log.Fatal().AnErr("error: ", err).Msgf("fail to create dir")
	}
	err = initDatabase(configDir)
	if err != nil {
		log.Fatal().AnErr("error: ", err).Msgf("fail to init database")
	}
}
