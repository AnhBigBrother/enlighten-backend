package cfg

import (
	"os"
)

type EnvCfg struct {
	DbUri           string
	JwtSecret       string
	Port            string
	AccessTokenAge  int64
	RefreshTokenAge int64
}

var Envs EnvCfg

func init() {
	Envs = EnvCfg{
		DbUri:           os.Getenv("DB_URI"),
		JwtSecret:       os.Getenv("JWT_SECRET"),
		Port:            os.Getenv("PORT"),
		AccessTokenAge:  30 * 60 * 1000,
		RefreshTokenAge: 7 * 24 * 60 * 60 * 1000,
	}
}
