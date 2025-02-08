package config

const (
	File = "config.json"
)

type Config struct {
	DataDirectory string `json:"dataDirectory"`
}
