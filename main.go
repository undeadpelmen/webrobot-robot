package webrobotrobot

func main() {

}

type Robot interface {
	SetSpeed(int) error
	Forward() error
	Stop() error
	Backward() error
	Left() error
	Right() error
}

func RobotControlFunc(cmd chan string, errch chan error, robot Robot) {

}
