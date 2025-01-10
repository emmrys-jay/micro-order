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
}

type AdminConfiguration struct {
	Email    string
	Password string
}

type Configuration struct {
	App      AppConfiguration
	Server   ServerConfiguration
	Database DatabaseConfiguration
	Redis    RedisConfiguration
	Token    TokenConfiguration
	Admin    AdminConfiguration
}
