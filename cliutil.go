package cliutil

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli"
)

func GenerateFlags(config interface{}) ([]cli.Flag, error) {
	t := reflect.TypeOf(config)
	if t.Kind() != reflect.Ptr {
		return nil, errors.New("cliutil: Config has to be a pointer of struct")
	}

	v := reflect.ValueOf(config)
	if v.Elem().Kind() != reflect.Struct {
		return nil, errors.New("cliutil: Config has to be a pointer of struct")
	}

	result := []cli.Flag{}

	el := t.Elem()
	for i := 0; i < el.NumField(); i++ {
		addr := v.Elem().Field(i).Addr().Interface()
		result = append(result,
			newFlag(
				el.Field(i).Type.Kind(),
				el.Field(i),
				addr,
			),
		)
	}

	return result, nil
}

func newFlag(kind reflect.Kind, f reflect.StructField, dest interface{}) cli.Flag {
	flagName, flagEnvName := getFlagName(f), getFlagEnvName(f)
	switch kind {
	case reflect.String:
		return &cli.StringFlag{
			Name:        flagName,
			EnvVar:      flagEnvName,
			Destination: dest.(*string),
			Required:    isTagExist(f, "required"),
			Hidden:      isTagExist(f, "hidden"),
			Usage:       getFlagUsage(f),
			Value:       getFlagDefault(f),
		}
	case reflect.Bool:
		return &cli.BoolFlag{
			Name:     flagName,
			EnvVar:   flagEnvName,
			Required: isTagExist(f, "required"),
			Hidden:   isTagExist(f, "hidden"),
			Usage:    getFlagUsage(f),
		}
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(getFlagDefault(f), 64)
		if err != nil {
			val = 0
		}
		return &cli.Float64Flag{
			Name:        flagName,
			EnvVar:      flagEnvName,
			Destination: dest.(*float64),
			Required:    isTagExist(f, "required"),
			Hidden:      isTagExist(f, "hidden"),
			Usage:       getFlagUsage(f),
			Value:       val,
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		val, err := strconv.ParseInt(getFlagDefault(f), 10, 64)
		if err != nil {
			val = 0
		}
		return &cli.IntFlag{
			Name:        flagName,
			EnvVar:      flagEnvName,
			Destination: dest.(*int),
			Required:    isTagExist(f, "required"),
			Hidden:      isTagExist(f, "hidden"),
			Usage:       getFlagUsage(f),
			Value:       int(val),
		}
	case reflect.Int64, reflect.Uint64:
		switch f.Type.String() {
		case "time.Duration":
			val, err := time.ParseDuration(getFlagDefault(f))
			if err != nil {
				val = time.Second
			}
			return &cli.DurationFlag{
				Name:        flagName,
				EnvVar:      flagEnvName,
				Destination: dest.(*time.Duration),
				Required:    isTagExist(f, "required"),
				Hidden:      isTagExist(f, "hidden"),
				Usage:       getFlagUsage(f),
				Value:       val,
			}
		default:
			val, err := strconv.ParseInt(getFlagDefault(f), 10, 64)
			if err != nil {
				val = 0
			}
			return &cli.Int64Flag{
				Name:        flagName,
				EnvVar:      flagEnvName,
				Destination: dest.(*int64),
				Required:    isTagExist(f, "required"),
				Hidden:      isTagExist(f, "hidden"),
				Usage:       getFlagUsage(f),
				Value:       val,
			}
		}
	default:
		panic(f.Type.String() + " is not supported type of flag")
	}
}

func getFlagName(s reflect.StructField) string {
	if val, ok := s.Tag.Lookup("flag"); ok {
		return val
	}

	return toSnakeCase(s.Name)
}

func getFlagUsage(s reflect.StructField) string {
	if val, ok := s.Tag.Lookup("usage"); ok {
		return val
	}

	return "cliutil: tag usage is not set"
}

func getFlagEnvName(s reflect.StructField) string {
	if val, ok := s.Tag.Lookup("env"); ok {
		return val
	}

	return strings.ToUpper(
		strings.ReplaceAll(
			getFlagName(s),
			"-",
			"_",
		),
	)
}

func isTagExist(s reflect.StructField, tagName string) bool {
	if val, ok := s.Tag.Lookup(tagName); ok && val == "true" {
		return true
	}

	return false
}

func getFlagDefault(s reflect.StructField) string {
	if val, ok := s.Tag.Lookup("default"); ok {
		return val
	}

	return ""
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}-${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}-${2}")
	return strings.ToLower(snake)
}

// // AppHelpTemplate is the text template for the Default help topic.
// // cli.go uses text/template to render templates. You can
// // render custom help text by setting this variable.
// const AppHelpTemplate = `Usage: {{if .UsageText}}{{.UsageText}}{{else}}potato {{if .VisibleFlags}}[options]{{end}}{{if .ArgsUsage}}{{.ArgsUsage}}{{else}} <agent|server>{{end}}{{end}}
// Options:
//    {{range $index, $option := .VisibleFlags}}{{if $index}}
//    {{end}}{{$option}}{{end}}
// `

// func NewApp(name string, config interface{}) {
// 	var config struct {
// 		NodeID           string
// 		LogLevel         string
// 		ListenAddress    string
// 		AdvertiseAddress string
// 		MigrationPath    string
// 		DatabaseDSN      string
// 		SkipMigration    bool
// 	}
// 	cli.AppHelpTemplate = AppHelpTemplate

// 	log := log.New(os.Stdout, "", 0)

// 	app := cli.NewApp()
// 	app.Version = printVersion()
// 	app.Flags = []cli.Flag{
// 		cli.StringFlag{
// 			Name:  "env.file, e",
// 			Usage: "Environment variable file path",
// 		},
// 		cli.StringFlag{
// 			Name:        "node-id",
// 			Usage:       "Node identifier",
// 			EnvVar:      "NODE_ID",
// 			Destination: &config.NodeID,
// 		},
// 		cli.StringFlag{
// 			Name:        "log-level",
// 			Usage:       "Log level",
// 			Value:       "info",
// 			EnvVar:      "LOG_LEVEL",
// 			Destination: &config.LogLevel,
// 		},
// 		cli.StringFlag{
// 			Name:        "listen-address",
// 			Usage:       "Port to wait incoming request from client",
// 			Value:       "0.0.0.0:9000",
// 			EnvVar:      "LISTEN_ADDRESS",
// 			Destination: &config.ListenAddress,
// 		},
// 		cli.StringFlag{
// 			Name:        "advertise-address",
// 			Usage:       "Port to advertise metrics over http",
// 			Value:       "0.0.0.0:9001",
// 			EnvVar:      "ADVERTISE_ADDRESS",
// 			Destination: &config.AdvertiseAddress,
// 		},
// 		cli.StringFlag{
// 			Name:        "migration-path",
// 			Usage:       "Database migration path for updating servers",
// 			Value:       "./migrations",
// 			EnvVar:      "MIGRATION_PATH",
// 			Destination: &config.MigrationPath,
// 		},
// 		cli.StringFlag{
// 			Name:        "database-dsn",
// 			Usage:       "Database data source name",
// 			Value:       "postgres://potato:potato@localhost:5432/potato?sslmode=disable",
// 			EnvVar:      "DATABASE_DSN",
// 			Destination: &config.DatabaseDSN,
// 		},
// 		cli.BoolFlag{
// 			Name:        "skip-migration",
// 			Usage:       "Skip database migration",
// 			EnvVar:      "SKIP_MIGRATION",
// 			Destination: &config.SkipMigration,
// 		},
// 	}

// 	app.Before = func(ctx *cli.Context) error {
// 		if envFile := ctx.String("env.file"); envFile != "" {
// 			return godotenv.Load(envFile)
// 		}

// 		if config.NodeID == "" {
// 			config.NodeID, _ = os.Hostname()
// 		}

// 		return nil
// 	}

// 	app.Action = func(ctx *cli.Context) (err error) {
// 		l := logrus.New()
// 		l.SetLevel(logrus.InfoLevel)

// 		dockerClient, err := client.NewEnvClient()
// 		if err != nil {
// 			return err
// 		}

// 		ctxx := context.Background()
// 		switch ctx.Args().First() {
// 		case "server":
// 			if !config.SkipMigration {
// 				err = migrateUp(config.MigrationPath, config.DatabaseDSN)
// 			}
// 			if err == nil {
// 				err = server.NewServer(l, config.ListenAddress, config.DatabaseDSN).Serve(ctxx)
// 			}
// 		case "agent":
// 			err = agent.NewAgent(l, dockerClient, config.NodeID, config.ListenAddress, config.AdvertiseAddress).Start(ctxx)
// 		default:
// 			return errors.New("This command takes one argument: <agent|server>\nFor additional help try 'potato -help'")
// 		}
// 		return
// 	}

// 	if err := app.Run(os.Args); err != nil {
// 		log.Printf("%v", colors.Bold(colors.Red)(err))
// 		os.Exit(1)
// 	}
// }
