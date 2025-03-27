package model

type Config struct {
	Server   Server
	Postgres Postgres
}

type Server struct {
	Port string
}

type Postgres struct {
	Driver string
	Host   string
}
