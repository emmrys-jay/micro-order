package config

type RedisConfiguration struct {
	Address  string
	Password string
}

type DatabaseConfiguration struct {
	Protocol string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Url      string // override the other params
}

type ServerConfiguration struct {
	HttpUrl            string
	HttpPort           string
	HttpAllowedOrigins string
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
	OwnerUrl   string
	ProductUrl string
}

type Configuration struct {
	App       AppConfiguration
	Server    ServerConfiguration
	Database  DatabaseConfiguration
	Redis     RedisConfiguration
	Token     TokenConfiguration
	Discovery DiscoveryConfiguration
}
