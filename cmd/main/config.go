package main

import "github.com/go-hao/configurator/ctype"

type Config struct {
	App struct {
		Mode      ctype.String `yaml:"mode" default:"DEBUG"`
		Name      ctype.String `yaml:"name" default:"app"`
		Owner     ctype.String `yaml:"owner" default:"www.example.com"`
		SecretKey ctype.String `yaml:"secret-key" default:"changeme!"`
	} `yaml:"app"`

	Server struct {
		Host ctype.String `yaml:"host" default:"EMPTY"`
		Port ctype.Int    `yaml:"port" default:"8000"`
		Path ctype.String `yaml:"path" default:"EMPTY"`
		CORS struct {
			AllowOrigins     ctype.Slice        `yaml:"allow-origins" default:"[*]"`
			AllowMethods     ctype.Slice        `yaml:"allow-methods" default:"[GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS]"`
			AllowHeaders     ctype.Slice        `yaml:"allow-headers" default:"[]"`
			AllowCredentials ctype.Bool         `yaml:"credentials" default:"true"`
			MaxAge           ctype.TimeDuration `yaml:"max-age" default:"12"`
			MaxAgeUnit       ctype.String       `yaml:"max-age-unit" default:"hr"`
		} `yaml:"cors"`
	} `yaml:"server"`

	Auth struct {
		TokenSigningAlg      ctype.String       `yaml:"token-signing-alg" default:"HS256"`
		TokenLifetimeUnit    ctype.String       `yaml:"token-lifetime-unit" default:"min"`
		AccessTokenLifetime  ctype.TimeDuration `yaml:"access-token-lifetime" default:"10"`
		RefreshTokenLifetime ctype.TimeDuration `yaml:"refresh-token-lifetime" default:"30"`
		IDTokenLifetime      ctype.TimeDuration `yaml:"id-token-lifetime" default:"10"`
		PrivateKeyPath       ctype.String       `yaml:"private-key-path" default:"certs/key.pem"`
		PublicKeyPath        ctype.String       `yaml:"public-key-path" default:"certs/key.pem.pub"`
	} `yaml:"auth"`

	DB struct {
		Host        ctype.String `yaml:"host" default:"127.0.0.1"`
		Port        ctype.Int    `yaml:"port" default:"3306"`
		Name        ctype.String `yaml:"name" default:"sql_db"`
		TablePrefix ctype.String `yaml:"table-prefix" default:"t"`
		Username    ctype.String `yaml:"username" default:"sql"`
		Password    ctype.String `yaml:"password" default:"letmein"`
	} `yaml:"db"`
}
