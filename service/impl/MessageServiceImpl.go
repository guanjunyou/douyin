package impl

import (
	"encoding/json"
	"fmt"
	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/utils"
	"io"
	"net"
	"sort"
	"sync"
)

var chatConnMap = sync.Map{}

type MessageServiceImpl struct {
}

func (messageService MessageServiceImpl) SendMessage(userId int64, toUserId int64, content string) error {
	err := models.SaveMessage(&models.Message{
		CommonEntity: utils.NewCommonEntity(),
		Content:      content,
	})
	if err != nil {
		return err
	}
	err = models.SaveMessageSendEvent(&models.MessageSendEvent{
		CommonEntity: utils.NewCommonEntity(),
		UserId:       userId,
		ToUserId:     toUserId,
		MsgContent:   content,
	})
	if err != nil {
		return err
	}
	err = models.SaveMessagePushEvent(&models.MessagePushEvent{
		CommonEntity: utils.NewCommonEntity(),
		FromUserId:   userId,
		MsgContent:   content,
	})
	if err != nil {
		return err
	}
	return nil
}

func (messageService MessageServiceImpl) GetHistoryOfChat(userId int64, toUserId int64) ([]models.MessageDVO, error) {
	//find from meesageSendEvent table
	messageSendEvents, err := models.FindMessageSendEventByUserIdAndToUserId(userId, toUserId)
	if err != nil {
		return nil, err
	}
	messageSendEventsOpposite, err := models.FindMessageSendEventByUserIdAndToUserId(toUserId, userId)
	if err != nil {
		return nil, err
	}
	messageSendEvents = append(messageSendEvents, messageSendEventsOpposite...)
	sort.Sort(models.ByCreateTime(messageSendEvents))

	var messages []models.MessageDVO
	var wg sync.WaitGroup
	for _, messageSendEvent := range messageSendEvents {
		wg.Add(1)
		go func(messageSendEvent models.MessageSendEvent) {
			defer wg.Done()
			messages = append(messages, models.MessageDVO{
				Id:         messageSendEvent.Id,
				UserId:     messageSendEvent.UserId,
				ToUserId:   messageSendEvent.ToUserId,
				Content:    messageSendEvent.MsgContent,
				CreateTime: messageSendEvent.CreateDate.Unix(),
			})
		}(messageSendEvent)
	}
	wg.Wait()
	return messages, nil
}

func RunMessageServer() {
	listen, err := net.Listen("tcp", "127.0.0.1:9090")
	if err != nil {
		fmt.Printf("Run message sever failed: %v\n", err)
		return
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("Accept conn failed: %v\n", err)
			continue
		}

		go process(conn)
	}
}

func process(conn net.Conn) {
	defer conn.Close()

	var buf [256]byte
	for {
		n, err := conn.Read(buf[:])
		if n == 0 {
			if err == io.EOF {
				break
			}
			fmt.Printf("Read message failed: %v\n", err)
			continue
		}

		var event = models.MessageSendEvent{}
		_ = json.Unmarshal(buf[:n], &event)
		fmt.Printf("Receive Message：%+v\n", event)

		fromChatKey := fmt.Sprintf("%d_%d", event.UserId, event.ToUserId)
		if len(event.MsgContent) == 0 {
			chatConnMap.Store(fromChatKey, conn)
			continue
		}

		toChatKey := fmt.Sprintf("%d_%d", event.ToUserId, event.UserId)
		writeConn, exist := chatConnMap.Load(toChatKey)
		if !exist {
			fmt.Printf("User %d offline\n", event.ToUserId)
			continue
		}

		pushEvent := models.MessagePushEvent{
			FromUserId: event.UserId,
			MsgContent: event.MsgContent,
		}
		pushData, _ := json.Marshal(pushEvent)
		_, err = writeConn.(net.Conn).Write(pushData)
		if err != nil {
			fmt.Printf("Push message failed: %v\n", err)
		}
	}
}
