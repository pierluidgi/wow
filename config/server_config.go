package config

type ServerConfig struct {
	LogLevel       string `yaml:"log-level"`
	QuotesFilename string `yaml:"quotes-filename"`
	Server         struct {
		Listen       string `yaml:"listen"`
		ReadTimeout  int    `yaml:"read-timeout"`
		WriteTimeout int    `yaml:"write-timeout"`
		DDoSRate     int    `yaml:"ddos-rate"`
		TargetBits   int    `yaml:"target-bits"`
		ChallengeTtl int    `yaml:"challenge-ttl"`
	} `yaml:"server"`
	CacheTtl     int `yaml:"cache-ttl"`
	RateInterval int `yaml:"rate-interval"`
	RateSize     int `yaml:"rate-size"`
}
