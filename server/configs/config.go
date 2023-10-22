package configs

type Config struct {
	Mysql  Mysql  `json:"mysql"`
	Logger Logger `json:"logger"`
	System System `json:"system"`
}

type Mysql struct {
	Conf string
}

type Logger struct {
	Level string
}

type System struct {
	Mode string
}

var Conf = new(Config)
