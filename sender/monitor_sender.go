package sender

type MonitorSender interface {
	Send(message string)
	String() string
}
