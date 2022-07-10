package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	vk_api "github.com/SevereCloud/vksdk/v2/api"
)

func main() {
	if len(os.Args) < 3 {
		log.Panic("invalid argument count")
	}

	token, groupID := os.Args[1], os.Args[2]

	api := vk_api.NewVK(token)

	ticker := time.NewTicker(time.Minute)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	lastTick := time.Now().Unix()

	for {
		select {
		case <-ticker.C:
			resp, err := api.VideoGet(vk_api.Params{
				"owner_id": groupID,
			})
			if err != nil {
				lastTick = time.Now().Unix()
				log.Println(err)
			}
			if len(resp.Items) == 0 {
				log.Println("Ничего нового")
			}

			for _, video := range resp.Items {
				fmt.Println()
				if video.Live && video.Date > int(lastTick) {
					log.Printf("Новая транляция!\nСсылка на трансляцию:%s\n", video.Player)
				}
			}

			lastTick = time.Now().Unix()
		case <-shutdown:
			log.Println("shutdown")
			return
		}
	}

}
