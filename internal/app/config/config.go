package config

type Config struct {
	Addr    string
	BaseURL string
}

type ConfigBuilder struct {
	conf Config
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{}
}

func (cb *ConfigBuilder) WithAddr(addr string) *ConfigBuilder {
	cb.conf.Addr = addr
	return cb
}

func (cb *ConfigBuilder) WithBaseURL(baseURL string) *ConfigBuilder {
	cb.conf.BaseURL = baseURL
	return cb
}

func (cb *ConfigBuilder) Build() Config {
	return cb.conf
}
