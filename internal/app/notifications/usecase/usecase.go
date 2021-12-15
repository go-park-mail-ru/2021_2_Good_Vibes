package usecase

import (
	"crypto/tls"
	"fmt"
	"github.com/flosch/pongo2/v4"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/config"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/notifications"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/microservice/orders"
	gomail "gopkg.in/mail.v2"
	"log"
	"sync"
	"time"
)

var templateMap map[string]string

func init() {
	templateMap = make(map[string]string, 3)
	templateMap["новый"] = "/home/bush/GolangTP/Ozon/2021_2_Good_Vibes/internal/app/notifications/templates/new_status.html"
	templateMap["в обработке"] = "/home/bush/GolangTP/Ozon/2021_2_Good_Vibes/internal/app/notifications/templates/processing_status.html"
	templateMap["передан курьеру"] = "/home/bush/GolangTP/Ozon/2021_2_Good_Vibes/internal/app/notifications/templates/courier_status.html"
}

type UseCase struct {
	notifyRepository notifications.Repository
	userRepository   user.Repository
	orderRepository  orders.Repository
}

func NewNotifyUseCase(notifyRepository notifications.Repository, userRepository user.Repository, orderRepository orders.Repository) *UseCase {
	return &UseCase{notifyRepository: notifyRepository, userRepository: userRepository, orderRepository: orderRepository}
}

func (uc *UseCase) SearchStatusChanges() error {
	changes, err := uc.notifyRepository.GetStatusChanges()
	if err != nil {
		return err
	}

	if changes == nil {
		return nil
	}

	wg := sync.WaitGroup{}
	wg.Add(len(changes))

	for _, change := range changes {
		change := change

		go func() {
			err := uc.ServeStatus(change)
			if err != nil {
				log.Printf("Cannot send email to %s, orderId: %d, userId: %d, error: %s",
					change.Email, change.OrderId, change.UserId, err.Error())
			}

			wg.Done()
		}()
	}
	wg.Wait()

	err = uc.notifyRepository.StableStatuses(changes)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) ServeStatus(change models.ChangedStatus) error {
	notifyInfo, err := uc.GetNotifyInfo(change)
	if err != nil {
		return err
	}

	err = uc.SendEmail(notifyInfo)

	tries := 0
	for err != nil {
		tries += 1

		time.Sleep(10 * time.Second)
		err = uc.SendEmail(notifyInfo)

		if tries > 2 {
			break
		}
		fmt.Println(tries)
	}
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) GetNotifyInfo(change models.ChangedStatus) (models.NotifyInfo, error) {
	var notifyInfo models.NotifyInfo

	userData, err := uc.userRepository.GetUserDataById(uint64(change.UserId))
	if err != nil {
		return models.NotifyInfo{}, err
	}

	notifyInfo.UserName = userData.Name
	notifyInfo.UserEmail = change.Email
	notifyInfo.OrderStatus = change.Status
	notifyInfo.OrderId = change.OrderId

	if change.Status == "передан курьеру" {
		address, err := uc.notifyRepository.GetAddressInfo(change.OrderId)
		if err != nil {
			return models.NotifyInfo{}, err
		}

		notifyInfo.Address = notifications.FromModelAddressToString(address)
	}

	if change.Status == "новый" {
		order, err := uc.orderRepository.GetOrderById(change.OrderId)
		if err != nil {
			return models.NotifyInfo{}, err
		}

		notifyInfo.OrderData = order
	}

	return notifyInfo, nil
}

func (uc *UseCase) SendEmail(notifyInfo models.NotifyInfo) error {
	tmp := pongo2.Must(pongo2.FromFile(templateMap[notifyInfo.OrderStatus]))

	out, err := uc.putContext(tmp, notifyInfo)
	if err != nil {
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", config.ConfigApp.Email.Address)
	message.SetHeader("To", notifyInfo.UserEmail)
	message.SetHeader("Subject", "Обновление статуса заказа")
	message.SetBody("text/html", out)
	dialer := gomail.NewDialer(config.ConfigApp.Email.Server,
		config.ConfigApp.Email.ServerPort,
		config.ConfigApp.Email.Address,
		config.ConfigApp.Email.Password)

	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := dialer.DialAndSend(message); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) putContext(tmp *pongo2.Template, notifyInfo models.NotifyInfo) (string, error) {
	var ctx pongo2.Context

	if notifyInfo.OrderStatus == "новый" {
		ctx = pongo2.Context{"name": notifyInfo.UserName, "order_id": notifyInfo.OrderId}
	}
	if notifyInfo.OrderStatus == "в обработке" {
		ctx = pongo2.Context{"name": notifyInfo.UserName, "order_id": notifyInfo.OrderId}
	}
	if notifyInfo.OrderStatus == "передан курьеру" {
		ctx = pongo2.Context{"name": notifyInfo.UserName, "order_id": notifyInfo.OrderId, "address": notifyInfo.Address}
	}
	out, err := tmp.Execute(ctx)
	if err != nil {
		return "", err
	}

	return out, nil
}
