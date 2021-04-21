package config

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
	Backup struct {
		Path struct {
			Local  string `yaml:"local"`
			Remote string `yaml:"remote"`
		}
		Command struct {
			Main string   `yaml:"main"`
			Args []string `yaml:"args"`
		}
		Remove bool `yaml:"remove"`
	}
	Build struct {
		Frontend struct {
			Root        string   `yaml:"root"`
			ClearPath string   `yaml:"clear-path"`
			Parallel    int      `yaml:"parallel"`
			Recursive   []string `yaml:"recursive"`
			IfExistFile string   `yaml:"if-exist-file"`
			Command     struct {
				Main string   `yaml:"main"`
				Args []string `yaml:"args"`
			}
		}
	}
}
