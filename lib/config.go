package lib

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type ConfigDatabase struct {
	Connection            string
	MaxOpenConnections    int `yaml:"maxOpenConnections"`
	MaxIdleConnections    int `yaml:"maxIdleConnections"`
	ConnectionMaxLifetime int `yaml:"connectionMaxLifetime"`
}
type ConfigHttp struct {
	Address string
}
type ConfigPicStorage struct {
	Storage string
	Temp    string
}
type TConfig struct {
	Database ConfigDatabase
	Http     ConfigHttp
	Pictures ConfigPicStorage
}

var (
	Config = TConfig{}
)

func (self *TConfig) Init() error {
	var err error
	err = self.tryReadFromExecutableDir()
	if err == nil {
		return nil
	} else {
		log.Printf("config not parsed from executable dir: %v", err)
	}
	err = self.tryReadFromWorkingDir()
	if err == nil {
		return nil
	} else {
		log.Printf("config not parsed from executable dir: %v", err)
	}
	return fmt.Errorf("can't find any suitable config file")
}

//------------------------- private -------------------------
const configFileName = "api.yml"

func (self *TConfig) postProcess() error {
	var err error

	self.Database.Connection, err = fixFilePath(self.Database.Connection)
	if err != nil {
		return err
	}

	self.Pictures.Storage, err = fixDirectoryPath(self.Pictures.Storage)
	if err != nil {
		return err
	}

	self.Pictures.Temp, err = fixDirectoryPath(self.Pictures.Temp)
	if err != nil {
		return err
	}

	return err
}

func (self *TConfig) tryReadFromExecutableDir() error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}
	var executablePath string
	executablePath, err = filepath.Abs(filepath.Dir(executable))
	if err != nil {
		return err
	}
	return self.tryReadConfigFile(executablePath)
}

func (self *TConfig) tryReadFromWorkingDir() error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}
	return self.tryReadConfigFile(workingDir)
}

func (self *TConfig) tryReadConfigFile(dir string) error {
	path := path.Join(dir, configFileName)
	log.Printf("trying to read config from %v", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	tmpConfig := TConfig{}
	err = yaml.UnmarshalStrict([]byte(data), &tmpConfig)
	if err != nil {
		return err
	}
	err = tmpConfig.postProcess()
	if err != nil {
		return err
	}
	*self = tmpConfig
	log.Printf("config initialized: %v", self)
	return nil
}

func fixFilePath(path string) (string, error) {
	_, err := fixDirectoryPath(filepath.Dir(path))
	if err != nil {
		return "", err
	}
	return filepath.Abs(path)
}

func fixDirectoryPath(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		info, err = os.Stat(path)
	}
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("path %v is not directory", path)
	}
	return filepath.Abs(path)
}
