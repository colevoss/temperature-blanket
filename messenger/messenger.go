package messenger

import "log"

type Messenger interface {
	SendMessage(to string, message string) error
}

type MockMessenger struct {
}

func NewMockMessenger() *MockMessenger {
	return &MockMessenger{}
}

func (m *MockMessenger) SendMessage(to string, message string) error {
	log.Println("\n" + message)

	return nil
}
