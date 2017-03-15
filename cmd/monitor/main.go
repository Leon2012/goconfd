package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	_ "path/filepath"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/Leon2012/goconfd/libs/version"
	"github.com/Leon2012/goconfd/monitor"
	"github.com/judwhite/go-svc/svc"
	"github.com/mreiferson/go-options"
)

var (
	flagSet     = flag.NewFlagSet("monitor", flag.ExitOnError)
	config      = flagSet.String("config", "", "config file")
	showVersion = flagSet.Bool("version", false, "show version")
	rpcAddress  = flagSet.String("rpc_address", "0.0.0.0:3002", "rpc address")
	dbUrl       = flagSet.String("db_url", "127.0.0.1:27017", "mongodb url")
	dbTimeout   = flagSet.Int("db_timeout", 5, "mongodb timeout")
	dbName      = flagSet.String("db_name", "goconfd", "mongodb name")
	dbUser      = flagSet.String("db_user", "", "mongodb user ")
	dbPass      = flagSet.String("db_pass", "", "mongodb pass")
)

type program struct {
	monitor *monitor.Monitor
}

func main() {
	prg := &program{}
	if err := svc.Run(prg, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.Fatal(err)
	}
}

func (p *program) Init(env svc.Environment) error {
	return nil
}

func (p *program) Stop() error {
	p.monitor.Exit()
	return nil
}

func (p *program) Start() error {
	flagSet.Parse(os.Args[1:])
	if *showVersion {
		fmt.Println(version.String("monitor"))
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

	opts := monitor.NewOptions()
	options.Resolve(opts, flagSet, cfg)
	log.Println(opts.String())
	err := opts.Valid()
	if err != nil {
		return err
	}
	daemon := monitor.NewMonitor(opts)
	daemon.Main()
	p.monitor = daemon
	return nil
}
