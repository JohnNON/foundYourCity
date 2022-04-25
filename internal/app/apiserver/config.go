package apiserver

// Config - содержит конфигурацию для запуска сервера
type Config struct {
	BindAddr string
	EndPoint string
	Source   string
}

// NewConfig - инициализация конфига по умолчанию
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		EndPoint: "/api/search",
		Source:   "cities500.zip",
	}
}
