package server

type Config struct {
	Port       int      `yaml:"port"`
	Path       string   `yaml:"path"`
	ServerName string   `yaml:"server_name"`
	Allowed    []string `yaml:"allowed"`
}
