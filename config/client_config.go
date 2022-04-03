package config

type ClientConfig struct {
	LogLevel string `yaml:"log-level"`
	Client   struct {
		ServerAddr   string `yaml:"server-addr"`
		ConnTimeout  int    `yaml:"conn-timeout"`
		ReadTimeout  int    `yaml:"read-timeout"`
		WriteTimeout int    `yaml:"write-timeout"`
	} `yaml:"client"`
	ParallelRequests int `yaml:"parallel-requests"`
	NextQuoteDelayMs int `yaml:"next-quote-delay-ms"`
}
