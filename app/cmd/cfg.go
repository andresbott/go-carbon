package cmd

import "strconv"

var DefaultCfg = appCfg{
	Main: serverCfg{
		BindIp: "",
		Port:   8085,
	},
	Obs: serverCfg{
		BindIp: "",
		Port:   9090,
	},
	Log: logConfig{
		Level: "info",
	},
}

type appCfg struct {
	Main serverCfg
	Obs  serverCfg `config:"Observability"`
	Log  logConfig
}

type serverCfg struct {
	BindIp string
	Port   int
}
type logConfig struct {
	Level string
}

func (c serverCfg) Addr() string {
	if c.BindIp == "" {
		return ":" + strconv.Itoa(c.Port)
	}

	return c.BindIp + ":" + strconv.Itoa(c.Port)
}
