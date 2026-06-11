package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	MySQL  MySQLConfig  `json:"mysql"`
	Admin  AdminConfig  `json:"admin"`
	Server ServerConfig `json:"server"`
	JWT    JWTConfig    `json:"jwt"`
}

type MySQLConfig struct {
	Address      string `json:"address"`
	Port         int    `json:"port"`
	DBName       string `json:"db-name"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Config       string `json:"config"`
	MaxIdleConns int    `json:"max-idle-conns"`
	MaxOpenConns int    `json:"max-open-conns"`
}

type AdminConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ServerConfig struct {
	Port int `json:"port"`
}

type JWTConfig struct {
	Secret      string `json:"secret"`
	ExpireHours int    `json:"expire-hours"`
}

var C Config

func Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &C)
}

// DSN 返回MySQL连接字符串
func (m MySQLConfig) DSN() string {
	return m.Username + ":" + m.Password +
		"@tcp(" + m.Address + ":" + itoa(m.Port) + ")/" + m.DBName +
		"?" + m.Config
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}
