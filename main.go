package main

import (
	interaction "github.com/seed30/TikTok/kitex_gen/douyin/interaction/interactionservice"
	"log"
)

func main() {
	svr := interaction.NewServer(new(InteractionServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
