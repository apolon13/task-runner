package config

type Command struct {
	Main string   `yaml:"main"`
	Args []string `yaml:"args"`
}

type Branch struct {
	Name    string  `yaml:"name"`
	Amend   bool    `yaml:"amend"`
	Command Command `yaml:"command"`
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
	Backup struct {
		Path struct {
			Local  string `yaml:"local"`
			Remote string `yaml:"remote"`
		}
		Command Command
		Remove  bool `yaml:"remove"`
	}
	Build struct {
		Frontend struct {
			Root        string   `yaml:"root"`
			CutExecPath string   `yaml:"cut-exec-path"`
			Parallel    int      `yaml:"parallel"`
			Recursive   []string `yaml:"recursive"`
			CheckFile   string   `yaml:"check-file"`
			Command     Command
		}
	}
	Git struct {
		Intermediate []Branch `yaml:"intermediate"`
	}
}
