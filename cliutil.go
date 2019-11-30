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

func MustGenerateFlags(config interface{}) []cli.Flag {
	flags, err := GenerateFlags(config)
	if err != nil {
		panic(err)	
	}
	return flags
}

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
	case reflect.Float64:
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
	case reflect.Int:
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
	case reflect.Int64:
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
