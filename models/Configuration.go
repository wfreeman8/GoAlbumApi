package models

import (
  "fmt"
  "os"
  "io"
  "errors"
  "encoding/json"
  "encoding/hex"
  "crypto/sha256"
  "crypto/rand"
)

type Config struct {
  Salt string `json:"salt"`
  AdminPassword string `json:"adminPassword"`
  AlbumBasePath string `json:"albumBasePath"`
  GalleryTitle string  `json:"gallery_title"`
}

type ConfigFormatted struct {
  GalleryTitle string `json:"gallery_title"`
}

const configFilename = "config.json"


func HasConfiguration() bool {
  configFileMeta, err := os.Stat(configFilename)
  return err == nil && configFileMeta.Size() > 5
}

func GetConfiguration() (*Config, error) {
  _, err := os.Stat(configFilename)
  var config = new(Config)
  if err == nil {
    configFs, err := os.Open(configFilename)
    defer configFs.Close()
    if err == nil {
      configContent, _ := io.ReadAll(configFs)
      err = json.Unmarshal(configContent, &config)
      return config, nil
    }

  }
  return config, err
}

func (config *Config) UpsertBaseFolder() error {
  var err error
  baseFolderStat, err := os.Stat(config.AlbumBasePath)
  if errors.Is(err, os.ErrNotExist) {
    err = os.Mkdir(config.AlbumBasePath, 0664)
  } else if !baseFolderStat.IsDir() {
    err = errors.New("path already in use and is not folder")
  }
  return err
}

func (config *Config) ValidatePassword(password string) bool {
  saltBytes, _ := hex.DecodeString(config.Salt)
  passwordBytes := []byte(password)
  passwordBytes = append(passwordBytes, saltBytes...)
  hash := sha256.New()
  hash.Write(passwordBytes)
  hashBytes := hash.Sum(nil)
  return hex.EncodeToString(hashBytes) == config.AdminPassword
}

func (config *Config) generateSalt() []byte {
  salt := make([]byte, 15)
  _, err := rand.Read(salt[:])
  if err != nil {
    panic(err)
  }
  return salt
}

func (config *Config) SetPassword(newPassword string) {
  salt := config.generateSalt()
  config.Salt = hex.EncodeToString(salt)

  passwordBytes := []byte(newPassword)
  passwordBytes = append(passwordBytes, salt...)

  hash := sha256.New()
  hash.Write(passwordBytes)
  hashBytes := hash.Sum(nil)

  config.AdminPassword = hex.EncodeToString(hashBytes)
}

func (config *Config) Save(createConfig bool) (error) {
  _, err := os.Stat(configFilename)
  configFs := new(os.File)
  var openFileErr error

  if errors.Is(err, os.ErrNotExist) && createConfig {
    configFs, openFileErr = os.Create(configFilename)
  } else if createConfig == false {
    configFs, openFileErr = os.OpenFile(configFilename, os.O_TRUNC | os.O_WRONLY, 0644)
  } else {
    openFileErr = errors.New("Config could not be saved")
  }
  defer configFs.Close()

  if openFileErr != nil {
    os.Remove(configFilename)
    return openFileErr
  }

  configJson, err := json.Marshal(config)
  if err == nil {
    configFs.Write(configJson)
    fmt.Println("config.json saved")
  } else {
    configFs.Close()
    os.Remove(configFilename)
  }
  return err
}

func (config *Config) GetConfigFormatted() ConfigFormatted {
  return ConfigFormatted{
    GalleryTitle: config.GalleryTitle,
  }
}