package config

type RedisConfiguration struct {
	Address  string
	Password string
	Ttl      string
}

type DatabaseConfiguration struct {
	Protocol string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type ServerConfiguration struct {
	HttpUrl            string
	HttpPort           string
	HttpAllowedOrigins string
	GrpcUrl            string
	GrpcPort           string
}

type AppConfiguration struct {
	Name string
	Env  string
}

type TokenConfiguration struct {
	Duration string
	Secret   string
	Issuer   string
}

type DiscoveryConfiguration struct {
	OwnerUrl string
}

type RabbitMqConfiguration struct {
	User     string
	Password string
	Host     string
}

type Configuration struct {
	App       AppConfiguration
	Server    ServerConfiguration
	Database  DatabaseConfiguration
	Redis     RedisConfiguration
	Token     TokenConfiguration
	Discovery DiscoveryConfiguration
	Rabbitmq  RabbitMqConfiguration
}
