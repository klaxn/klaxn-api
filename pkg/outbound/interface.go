package outbound

type Sender interface {
	SendMessage(to, message string) error
}
