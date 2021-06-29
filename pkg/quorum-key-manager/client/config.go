package client

type Config struct {
	URL        string
	MetricsURL string
}

func NewConfig(url string) *Config {
	return &Config{
		URL: url,
	}
}
