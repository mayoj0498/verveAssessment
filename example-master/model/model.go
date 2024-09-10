package model

type Configs struct {
	Host         string `yaml:"HOST"`
	ReadTimeout  int    `yaml:"READ_TIMEOUT "`
	WriteTimeout int    `yaml:"WRITE_TIMEOUT"`
	KafkaServer  string `yaml:"KAFKA_SERVER"`
	KafkaTopic   string `yaml:"KAFKA_TOPIC"`
}
