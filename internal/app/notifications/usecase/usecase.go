package usecase

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/config"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/notifications"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user"
	gomail "gopkg.in/mail.v2"
	"os"
	"time"
)

type UseCase struct {
	notifyRepository notifications.Repository
	userRepository user.Repository
}

func NewNotifyUseCase(notifyRepository notifications.Repository, userRepository user.Repository) *UseCase {
	return &UseCase{notifyRepository: notifyRepository, userRepository: userRepository}
}

func (uc *UseCase) SearchStatusChanges() error {
	changes, err := uc.notifyRepository.GetStatusChanges()
	if err != nil {
		return err
	}

	for _, change := range changes {
		if change.Status == "новый" {
			go func() {
				err := uc.ServeNewStatuses(change)
				if err != nil {
					// TODO: -_-
					panic(err)
				}
			}()
		}
	}

	return nil
}

func (uc *UseCase) ServeNewStatuses(change models.ChangedStatus) error {
	notifyInfo, err := uc.GetNotifyInfo(change)
	if err != nil {
		return err
	}

	err = uc.SendEmail(notifyInfo)
	for err != nil {
		panic(err)
		time.Sleep(10 * time.Second)
		err = uc.SendEmail(notifyInfo)
	}
	fmt.Println("Письмо отослал")
	fmt.Println(change.Email)
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

	if change.Status == "доставлен" {
		address, err := uc.notifyRepository.GetAddressInfo(change.OrderId)
		if err != nil {
			return models.NotifyInfo{}, err
		}

		notifyInfo.Address = notifications.FromModelAddressToString(address)
	}

	return notifyInfo, nil
}

func (uc *UseCase) SendEmail(notifyInfo models.NotifyInfo) error {
	file, err := os.Open("/home/bush/GolangTP/Ozon/2021_2_Good_Vibes/internal/app/notifications/templates/new_status.html")
	if err != nil{
		panic(err)
	}
	defer file.Close()

	wr := bytes.Buffer{}
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		wr.WriteString(sc.Text())
	}

	message := gomail.NewMessage()
	message.SetHeader("From", config.ConfigApp.Email.Address)
	message.SetHeader("To", notifyInfo.UserEmail)
	message.SetHeader("Subject", "Обновление статуса заказа")
	// body := fmt.Sprintf(`Здравствуйте %s! Вы сделали заказ в магазине Azot! Информация об обновлении статуса заказа будет приходить вам на почту.`, notifyInfo.UserName)
	message.SetBody("text/html", wr.String())
	dialer := gomail.NewDialer(config.ConfigApp.Email.Server,
		config.ConfigApp.Email.ServerPort,
		config.ConfigApp.Email.Address,
		config.ConfigApp.Email.Password)

	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := dialer.DialAndSend(message); err != nil {
		return err
	}

	fmt.Println("Письмо доставлено")

	return nil
}
