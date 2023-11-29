package configs

type Config struct {
	Mysql         Mysql         `json:"mysql"`
	Logger        Logger        `json:"logger"`
	System        System        `json:"system"`
	Concurrency   Concurrency   `json:"concurrency"`
	WebhookSecret WebhookSecret `json:"webhook_secret"`
}

type Mysql struct {
	Conf            string
	CreateBatchSize int
}

type Logger struct {
	Level string
}

type Concurrency struct {
	Number int64
}

type System struct {
	Mode string
}

type WebhookSecret struct {
	Github string
	Gitlab string
}

var Conf = new(Config)
