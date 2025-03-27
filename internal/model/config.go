package model

type Config struct {
	Server   Server   `json:"server"`
	Postgres Postgres `json:"postgres"`
	Sqlite   Sqlite   `json:"sqlite"`
}

type Server struct {
	Port string `json:"port"`
}

type Sqlite struct {
	Driver string `json:"driver"`
	Host   string `json:"host"`
}

type Postgres struct {
	Driver string `json:"driver"`
	Host   string `json:"host"`
}
