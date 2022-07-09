package server

import (
	"context"
	"log"
	"math/rand"
	"strings"
	"time"

	vk_api "github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

type Server struct {
	api     *vk_api.VK
	userAPI *vk_api.VK
	lp      *longpoll.LongPoll
}

func NewServer(groupToken string, secretToken string) *Server {
	api := vk_api.NewVK(groupToken)
	userAPI := vk_api.NewVK(secretToken)

	groups, err := api.GroupsGetByID(nil)
	if err != nil {
		log.Panic(err)
	}

	if len(groups) != 1 {
		log.Panic(err)
	}

	group := groups[0]

	lp, err := longpoll.NewLongPoll(api, group.ID)
	if err != nil {
		log.Panic(err)
	}

	lp.MessageNew(func(ctx context.Context, mno events.MessageNewObject) {
		go func() {
			if strings.EqualFold(mno.Message.Text, "звонок") {
				var res interface{}
				err := userAPI.RequestUnmarshal("messages.startCall", &res, vk_api.Params{
					"group_id": group.ID,
				})
				log.Println(res, err)
			}
		}()
	})

	rand.Seed(time.Now().UnixNano())
	server := &Server{
		api:     api,
		userAPI: userAPI,
		lp:      lp,
	}

	return server
}

func (s *Server) Run() error {
	return s.lp.Run()
}
