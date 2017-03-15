package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	_ "path/filepath"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/Leon2012/goconfd/dashboard"
	"github.com/Leon2012/goconfd/libs/version"
	"github.com/judwhite/go-svc/svc"
	"github.com/mreiferson/go-options"
)

var (
	flagSet        = flag.NewFlagSet("dashboard", flag.ExitOnError)
	config         = flagSet.String("config", "", "config file")
	showVersion    = flagSet.Bool("version", false, "show version")
	httpAddress    = flagSet.String("http_address", "0.0.0.0:3003", "http address")
	dbUrl          = flagSet.String("db_url", "127.0.0.1:27017", "mongodb url")
	dbTimeout      = flagSet.Int("db_timeout", 5, "mongodb timeout")
	dbName         = flagSet.String("db_name", "goconfd", "mongodb name")
	dbUser         = flagSet.String("db_user", "", "mongodb user ")
	dbPass         = flagSet.String("db_pass", "", "mongodb pass")
	hosts          = flagSet.String("hosts", "localhost:2379", "etcd hosts")
	dialTimeout    = flagSet.Int("dial_timeout", 5, "dial timeout")
	requestTimeout = flagSet.Int("request_timeout", 5, "request timeout")
	templatePath   = flagSet.String("template_path", "./template", "template path")
)

type program struct {
	dashboard *dashboard.Dashboard
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
	p.dashboard.Exit()
	return nil
}

func (p *program) Start() error {
	flagSet.Parse(os.Args[1:])
	if *showVersion {
		fmt.Println(version.String("dashboard"))
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

	opts := dashboard.NewOptions()
	options.Resolve(opts, flagSet, cfg)
	log.Println(opts.String())
	err := opts.Valid()
	if err != nil {
		return err
	}
	daemon := dashboard.NewDashboard(opts)
	daemon.Main()
	p.dashboard = daemon
	return nil
}
