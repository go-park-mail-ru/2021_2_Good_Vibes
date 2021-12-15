package notifications

import (
	"time"
)

type Notifier struct {
	useCase UseCase
}

func NewNotifier(useCase UseCase) *Notifier {
	return &Notifier{useCase: useCase}
}

func (n *Notifier) Run() error {
	go func() {
		for {
			err := n.useCase.SearchStatusChanges()
			if err != nil {
				// TODO: -_-
				panic(err)
			}
			time.Sleep(10 * time.Second)
		}
	}()

	return nil
}
