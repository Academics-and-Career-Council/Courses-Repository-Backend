package main

import (
	"context"
	"log"

	"github.com/Mshivam2409/AnC-Courses/api"
	"github.com/Mshivam2409/AnC-Courses/services"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

func init() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/config/fiber/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panicf("Fatal error config file: %s", err)
	}
	services.ConnectMongo()
	services.InitalizeOryServices()
}

func main() {
	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		err := api.FiberListenAndServe()
		if err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		err := services.SMTPListenAndServe()
		if err != nil {
			return err
		}
		return nil
	})

	g.Wait()
}
