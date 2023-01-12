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
	log.Println("Sending Message to:" + to)
	log.Println("Message:\n" + message)

	return nil
}
