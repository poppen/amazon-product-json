package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/labstack/echo"
	"github.com/poppen/amazing"
)

const (
	version = "0.1.0"
)

// Config ...
type Config struct {
	Port      string          `toml:"port"`
	AmazonAPI AmazonAPIConfig `toml:"amazon"`
}

// AmazonAPIConfig  ...
type AmazonAPIConfig struct {
	AssociateTag  string `toml:"associate_tag"`
	AccessKey     string `toml:"access_key"`
	SecretKey     string `toml:"secret_key"`
	ServiceDomain string `toml:"service_domain"`
	ResponseGroup string `toml:"response_group"`
}

var conf Config

func main() {
	checkVersion()
	loadConfig()
	e := echo.New()
	e.GET("/items/:asin", getItem)
	e.Logger.Fatal(e.Start(":" + conf.Port))
}

func checkVersion() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("amazon-product-json version %s\n", version)
		os.Exit(0)
	}
}

func loadConfig() {
	var configPath string
	flag.StringVar(&configPath, "c", "config.tml", "configuration file path")
	flag.Parse()

	if _, err := toml.DecodeFile(configPath, &conf); err != nil {
		panic(err)
	}
}

func getItem(ctx echo.Context) error {
	asin := ctx.Param("asin")
	if len(asin) == 0 {
		return ctx.String(http.StatusBadRequest, "Asin is empty")
	}

	client, err := amazing.NewAmazing(conf.AmazonAPI.ServiceDomain, conf.AmazonAPI.AssociateTag, conf.AmazonAPI.AccessKey, conf.AmazonAPI.SecretKey)
	params := url.Values{
		"ResponseGroup": []string{conf.AmazonAPI.ResponseGroup},
	}
	res, err := client.ItemLookupAsin(asin, params)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err.Error()))
	}

	return ctx.JSON(http.StatusOK, res)
}
