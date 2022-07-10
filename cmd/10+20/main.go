package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"

	"github.com/Not-cottage-cheese-but-cottage-cheese/final-vk-calls/server"
)

func init() {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No such config file")
		} else {
			log.Println("Read config error")
		}

		log.Println("Get from OS env")
		viper.Set("GROUP_TOKEN", os.Getenv("GROUP_TOKEN"))
		viper.Set("SECRET", os.Getenv("SECRET"))
	}
}

func main() {
	groupToken := viper.GetString("GROUP_TOKEN")
	secretToken := viper.GetString("SECRET")

	server := server.NewServer(groupToken, secretToken)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		<-quit
		server.Shutdown()
	}()

	if err := server.Run(); err != nil {
		log.Panic(err)
	}
}
