package main

import (
	"fmt"
	"html/template"
	"log"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultServicePort  = 55555
	defaultServiceAddr  = "127.0.0.1"
	defaultPageSize     = 20
	initBasketCapacity  = 200
	maxBasketCapacity   = 2000
	defaultDatabaseType = DbTypeMemory
	serviceOldAPIPath   = "baskets"
	serviceAPIPath      = "api"
	serviceUIPath       = "web"
	serviceName         = "request-baskets"
	basketNamePattern   = `^[\w\d\-_\.]{1,250}$`
	sourceCodeURL       = "https://github.com/darklynx/request-baskets"
)

// ServerConfig describes server configuration.
type ServerConfig struct {
	ServerPort   int
	ServerAddr   string
	InitCapacity int
	MaxCapacity  int
	PageSize     int
	MasterToken  string
	DbType       string
	DbFile       string
	DbConnection string
	Baskets      []string
	PathPrefix   string
	Mode         string
	Theme        string
	ThemeCSS     template.HTML
}

type arrayFlags []string

func (v *arrayFlags) String() string {
	return strings.Join(*v, ",")
}

func (v *arrayFlags) Set(value string) error {
	*v = append(*v, value)
	return nil
}

// CreateConfig creates server configuration base on application command line arguments
func CreateConfig() *ServerConfig {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	// check if there environment variables matching current keys
	viper.AutomaticEnv()

	pflag.Int("port", defaultServicePort, "HTTP service port")
	pflag.String("listener", defaultServiceAddr, "HTTP listen address")
	pflag.Int("size", initBasketCapacity, "Initial basket size (capacity)")
	pflag.Int("maxsize", maxBasketCapacity, "Maximum allowed basket size (max capacity)")
	pflag.Int("page", defaultPageSize, "Default page size")
	pflag.String("token", "", "Master token, random token is generated if not provided")
	pflag.String("db-type", defaultDatabaseType, fmt.Sprintf(
		"Baskets storage type: \"%s\" - in-memory, \"%s\" - Bolt DB, \"%s\" - SQL database",
		DbTypeMemory, DbTypeBolt, DbTypeSQL))
	pflag.String("file", "./baskets.db", "Database location, only applicable for file or SQL databases")
	pflag.String("conn", "", "Database connection string for SQL databases, if undefined \"file\" argument is considered")
	pflag.String("prefix", "", "Service URL path prefix")
	pflag.String("mode", ModePublic, fmt.Sprintf(
		"Service mode: \"%s\" - any visitor can create a new basket, \"%s\" - baskets creation requires master token",
		ModePublic, ModeRestricted))
	pflag.String("theme", ThemeStandard, fmt.Sprintf(
		"CSS theme for web UI, supported values: %s, %s, %s",
		ThemeStandard, ThemeAdaptive, ThemeFlatly))
	pflag.StringArray("basket", []string{}, "Name of a basket to auto-create during service startup (can be specified multiple times)")
	// var port = flag.Int("p", defaultServicePort, "HTTP service port")
	// var address = flag.String("l", defaultServiceAddr, "HTTP listen address")
	// var initCapacity = flag.Int("size", initBasketCapacity, "Initial basket size (capacity)")
	// var maxCapacity = flag.Int("maxsize", maxBasketCapacity, "Maximum allowed basket size (max capacity)")
	// var pageSize = flag.Int("page", defaultPageSize, "Default page size")
	// var masterToken = flag.String("token", "", "Master token, random token is generated if not provided")
	// var dbType = flag.String("db", defaultDatabaseType, fmt.Sprintf(
	// 	"Baskets storage type: \"%s\" - in-memory, \"%s\" - Bolt DB, \"%s\" - SQL database",
	// 	DbTypeMemory, DbTypeBolt, DbTypeSQL))
	// var dbFile = flag.String("file", "./baskets.db", "Database location, only applicable for file or SQL databases")
	// var dbConnection = flag.String("conn", "", "Database connection string for SQL databases, if undefined \"file\" argument is considered")
	// var prefix = flag.String("prefix", "", "Service URL path prefix")
	// var mode = flag.String("mode", ModePublic, fmt.Sprintf(
	// 	"Service mode: \"%s\" - any visitor can create a new basket, \"%s\" - baskets creation requires master token",
	// 	ModePublic, ModeRestricted))
	// var theme = flag.String("theme", ThemeStandard, fmt.Sprintf(
	// 	"CSS theme for web UI, supported values: %s, %s, %s",
	// 	ThemeStandard, ThemeAdaptive, ThemeFlatly))

	// var baskets arrayFlags
	// flag.Var(&baskets, "basket", "Name of a basket to auto-create during service startup (can be specified multiple times)")
	// flag.Parse()

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	var token = viper.GetString("token")
	if len(token) == 0 {
		token, _ = GenerateToken()
		log.Printf("[info] generated master token: %s", token)
	}

	return &ServerConfig{
		ServerPort:   viper.GetInt("port"),
		ServerAddr:   viper.GetString("listener"),
		InitCapacity: viper.GetInt("size"),
		MaxCapacity:  viper.GetInt("maxsize"),
		PageSize:     viper.GetInt("page"),
		MasterToken:  token,
		DbType:       viper.GetString("db-type"),
		DbFile:       viper.GetString("file"),
		DbConnection: viper.GetString("conn"),
		Baskets:      viper.GetStringSlice("basket"),
		PathPrefix:   normalizePrefix(viper.GetString("prefix")),
		Mode:         viper.GetString("mode"),
		Theme:        viper.GetString("theme"),
		ThemeCSS:     toThemeCSS(viper.GetString("theme"))}
}

func normalizePrefix(prefix string) string {
	if (len(prefix) > 0) && (prefix[0] != '/') {
		return "/" + prefix
	} else {
		return prefix
	}
}
