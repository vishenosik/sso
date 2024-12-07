package config

type Cache struct {
	Redis Redis
}
type Redis struct {
	Options Options
}

type Options struct {
	User     string
	Password string
	DB       int
	Host     string
	Port     int
}
