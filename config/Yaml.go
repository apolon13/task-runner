package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"task-runner/cmd"
	"task-runner/utils/builder/frontend"
	"task-runner/utils/grpc"
)

type Branch struct {
	Name    string      `yaml:"name"`
	Amend   bool        `yaml:"amend"`
	Command cmd.Command `yaml:"command"`
}

type Yaml struct {
	Connections struct {
		Ssh struct {
			Username   string `yaml:"username"`
			Host       string `yaml:"host"`
			Port       int    `yaml:"port"`
			Password   string `yaml:"password"`
			PrivateKey string `yaml:"private_key"`
		}
	}
	Restore struct {
		Db struct {
			Path struct {
				Local  string `yaml:"local"`
				Remote string `yaml:"remote"`
			}
			Command cmd.Command
			Remove  bool `yaml:"remove"`
		}
	}
	Build struct {
		Frontend frontend.ProcessParams
	}
	GRPC struct {
		Client grpc.ProtocParams
		Server grpc.ProtocParams
	}
	Git struct {
		Release struct {
			Intermediate []Branch `yaml:"intermediate"`
		}
	}
}

func New(path string) Yaml {
	yamlFile := Yaml{}
	confContent, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("config file not found: %s", err))
	}
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	confContent = []byte(os.ExpandEnv(string(confContent)))
	err = yaml.Unmarshal(confContent, &yamlFile)
	if err != nil {
		log.Fatal(fmt.Sprintf("unmarshal config file error: %s", err))
	}
	return yamlFile
}
