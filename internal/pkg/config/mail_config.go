package config

type Mail struct {
	Email    string `env:"MAIL_EMAIL"    env-required:"true"`
	Password string `env:"MAIL_PASSWORD" env-required:"true"`
	Address  string `env:"MAIL_ADDRESS"  env-required:"true"`
	Port     string `env:"MAIL_PORT"     env-required:"true"`
}
