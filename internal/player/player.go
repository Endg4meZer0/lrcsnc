package player

type controller interface {
	Connect() error
	Disconnect() error

	GetPosition() (float64, error)
	SetPosition(pos float64) error
}

var Controller controller
