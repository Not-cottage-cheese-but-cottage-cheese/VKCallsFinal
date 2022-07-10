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

	sendMessage := func(peerID int, message string) (int, error) {
		log.Printf("[INFO] send to %d: %s", peerID, message)
		builder := params.NewMessagesSendBuilder()
		builder.Message(message)
		builder.RandomID(0)
		builder.PeerID(peerID)
		return api.MessagesSend(vk_api.Params(builder.Params))
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
						sendMessage(operatorID, fmt.Sprintf("У вас новый клиент. Напоминаю вашу ссылку:\n%s", cc.GetLink(operatorID)))
					default:
						waitList.Store(mno.Message.PeerID, struct{}{})
						message = "Свободных операторов нет. С вами свяжутся при первой возможности"
						go func() {
							operatorID := <-cc.HasFree
							sendMessage(mno.Message.PeerID, fmt.Sprintf("Ссылка для связи с оператором: %s", cc.GetLink(operatorID)))
							sendMessage(operatorID, fmt.Sprintf("У вас новый клиент. Напоминаю вашу ссылку:\n%s", cc.GetLink(operatorID)))
							cc.SetBusy(operatorID)
						}()
					}
				}
			} else if strings.EqualFold(mno.Message.Text, "хочу быть оператором") {
				if cc.IsOperator(mno.Message.PeerID) {
					message = "Вы уже оператор!"
				} else {
					var res map[string]interface{}
					err := userAPI.RequestUnmarshal("messages.startCall", &res, vk_api.Params{})
					if err != nil {
						log.Println("[ERROR] messages.startCall: ", err)
					} else {
						link, ok := res["join_link"]
						if !ok {
							log.Println("[ERROR] messages.startCall: no link")
						} else {
							cc.AddOperator(mno.Message.PeerID, res["join_link"].(string))
							log.Printf("[INFO] new operator %d with link: %s\n", mno.Message.PeerID, link)
							message = fmt.Sprintf("Поздравляю! Теперь вы оператор! Ваша ссылка:\n%s", cc.GetLink(mno.Message.PeerID))
						}
					}
				}
			} else if strings.EqualFold(mno.Message.Text, "я свободен") {
				if cc.IsOperator(mno.Message.PeerID) {
					if cc.IsBusy(mno.Message.PeerID) {
						message = "Ожидайте следующего клиента"
						cc.SetFree(mno.Message.PeerID)
					} else {
						message = "Вы уже сообщали, что свободны. Ожидайте следующего клиента"
					}
				}
			}

			sendMessage(mno.Message.PeerID, message)
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
	log.Println("[INFO] Start serving")
	return s.lp.Run()
}

func (s *Server) Shutdown() {
	fmt.Println("[INFO] Gracefully shuttinh down...")
	s.lp.Shutdown()
}
