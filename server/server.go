package server

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	vk_api "github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

type Server struct {
	api     *vk_api.VK
	userAPI *vk_api.VK
	lp      *longpoll.LongPoll
	cc      *CallCenter

	waitList *sync.Map
}

func NewServer(groupToken string, secretToken string) *Server {
	api := vk_api.NewVK(groupToken)
	userAPI := vk_api.NewVK(secretToken)

	cc := NewCallCenter()
	waitList := &sync.Map{}

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
		log.Printf("[INFO] got message from %d: %s", mno.Message.PeerID, mno.Message.Text)
		go func() {
			message := "Уточните запрос"

			if strings.EqualFold(mno.Message.Text, "звонок") {
				var res map[string]interface{}
				err := userAPI.RequestUnmarshal("messages.startCall", &res, vk_api.Params{})
				if err != nil {
					log.Println("[ERROR] messages.startCall: ", err)
				} else {
					link, ok := res["join_link"]
					if !ok {
						log.Println("[ERROR] messages.startCall: no link")
					} else {
						log.Printf("[INFO] peer: %d, link: %s\n", mno.Message.PeerID, link)
						message = fmt.Sprintf("Ссылка на звонок: %s", link)
					}
				}
			} else if strings.EqualFold(mno.Message.Text, "звонок оператору") {
				_, ok := waitList.Load(mno.Message.PeerID)
				if ok {
					message = "Вы уже в очереди. Ожидайте"
				} else {
					select {
					case operatorID := <-cc.HasFree:
						message = fmt.Sprintf("Ссылка для связи с оператором: %s", cc.GetLink(operatorID))

						builder := params.NewMessagesSendBuilder()
						builder.Message(fmt.Sprintf("У вас новый клиент. Напоминаю вашу ссылку:\n%s", cc.GetLink(operatorID)))
						builder.RandomID(0)
						builder.PeerID(operatorID)
						api.MessagesSend(vk_api.Params(builder.Params))
					default:
						waitList.Store(mno.Message.PeerID, struct{}{})
						message = "Свободных операторов нет. С вами свяжутся при первой возможности"
						go func() {
							operatorID := <-cc.HasFree
							builder := params.NewMessagesSendBuilder()
							builder.Message(fmt.Sprintf("Ссылка для связи с оператором: %s", cc.GetLink(operatorID)))
							builder.RandomID(0)
							builder.PeerID(mno.Message.PeerID)
							api.MessagesSend(vk_api.Params(builder.Params))

							builder.Message(fmt.Sprintf("У вас новый клиент. Напоминаю вашу ссылку:\n%s", cc.GetLink(operatorID)))
							builder.RandomID(0)
							builder.PeerID(operatorID)
							api.MessagesSend(vk_api.Params(builder.Params))

							cc.SetBusy(operatorID)
						}()
					}
				}
			} else if strings.EqualFold(mno.Message.Text, "хочу быть оператором") {
				var res map[string]interface{}
				userAPI.RequestUnmarshal("messages.startCall", &res, vk_api.Params{})
				cc.AddOperator(mno.Message.PeerID, res["join_link"].(string))

				message = fmt.Sprintf("Поздравляю! Теперь вы оператор! Ваша ссылка:\n%s", cc.GetLink(mno.Message.PeerID))
			} else if strings.EqualFold(mno.Message.Text, "я свободен") {
				cc.SetFree(mno.Message.PeerID)
				return
			}

			builder := params.NewMessagesSendBuilder()
			builder.Message(message)
			builder.RandomID(0)
			builder.PeerID(mno.Message.PeerID)
			api.MessagesSend(vk_api.Params(builder.Params))
		}()
	})

	rand.Seed(time.Now().UnixNano())
	server := &Server{
		api:      api,
		userAPI:  userAPI,
		lp:       lp,
		cc:       cc,
		waitList: waitList,
	}

	return server
}

func (s *Server) Run() error {
	return s.lp.Run()
}
