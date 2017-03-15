package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	_ "path/filepath"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/Leon2012/goconfd/agent"
	"github.com/Leon2012/goconfd/libs/version"
	"github.com/judwhite/go-svc/svc"
	"github.com/mreiferson/go-options"
)

var (
	flagSet           = flag.NewFlagSet("agent", flag.ExitOnError)
	config            = flagSet.String("config", "", "config file")
	showVersion       = flagSet.Bool("version", false, "show version")
	httpAddress       = flagSet.String("http_address", "0.0.0.0:3001", "http address")
	dialTimeout       = flagSet.Int("dial_timeout", 5, "dial timeout")
	requestTimeout    = flagSet.Int("request_timeout", 5, "request timeout")
	keyPrefix         = flagSet.String("key_prefix", "", "key prefix")
	savePath          = flagSet.String("save_path", "/tmp/", "save path")
	saveType          = flagSet.Int("save_type", 1, "1=file, 2=shm")
	hosts             = flagSet.String("hosts", "localhost:2379", "etcd hosts")
	fileExt           = flagSet.String("file_ext", "php", "file ext")
	autoLoad          = flagSet.Bool("auto_load", false, "auto load local keys")
	heartbeatInterval = flagSet.Int("heartbeat_interval", 5, "heartbeat interval")
)

type program struct {
	agent *agent.Agent
}

func main() {
	prg := &program{}
	if err := svc.Run(prg, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.Fatal(err)
	}
}

func (p *program) Init(env svc.Environment) error {
	//log.Printf("is win service? %v\n", env.IsWindowsService())
	return nil
}

func (p *program) Stop() error {
	p.agent.Exit()
	return nil
}

func (p *program) Start() error {
	flagSet.Parse(os.Args[1:])
	if *showVersion {
		fmt.Println(version.String("agent"))
		os.Exit(0)
	}

	var cfg map[string]interface{}
	if *config != "" {
		_, err := toml.DecodeFile(*config, &cfg)
		if err != nil {
			log.Fatal("ERROR: failed to load config file %s - %s", *config, err.Error())
			return err
		}
	}

	opts := agent.NewOptions()
	options.Resolve(opts, flagSet, cfg)
	log.Println(opts.String())
	err := opts.Valid()
	if err != nil {
		return err
	}
	daemon := agent.NewAgent(opts)
	daemon.Main()
	p.agent = daemon
	return nil
}
