package bootstrap

import "flag"

type CommandFlags struct {
	Conf   string // 文件配置
	Env    string // 运行环境 dev 开发 ；pro 生产
	Consul string // consul 远程配置 地址
}

func NewCommandFlags() *CommandFlags {
	return &CommandFlags{
		Conf:   "",
		Env:    "",
		Consul: "",
	}
}

func (f *CommandFlags) Init() {
	flag.StringVar(&f.Conf, "conf", "../../configs", "config path, eg: -conf bootstrap.yaml")
	flag.StringVar(&f.Env, "env", "dev", "runtime environment, eg: -env dev")
	flag.StringVar(&f.Consul, "consul", "", "config server host, eg: -consul http://127.0.0.1:8500")
}
