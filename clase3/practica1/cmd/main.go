package main

import (
	"database/sql"
	"fmt"
	repository "godatabases/clase3/practica1/repository"

	"github.com/go-sql-driver/mysql"
)

// ConfigServerChi is the configuration for the server
type ConfigServerChi struct {
	// Addr is the address to listen on
	Addr string
	// MySQLDSN is the DSN for the MySQL database
	MySQLDSN string
}

// NewServerChi creates a new instance of the server
func NewServerChi(cfg ConfigServerChi) *ServerChi {
	// default config
	defaultCfg := ConfigServerChi{
		Addr:     ":8080",
		MySQLDSN: "",
	}
	if cfg.Addr != "" {
		defaultCfg.Addr = cfg.Addr
	}
	if cfg.MySQLDSN != "" {
		defaultCfg.MySQLDSN = cfg.MySQLDSN
	}

	return &ServerChi{
		addr:     defaultCfg.Addr,
		mysqlDSN: defaultCfg.MySQLDSN,
	}
}

// ServerChi is the default implementation of the server
type ServerChi struct {
	// addr is the address to listen on
	addr string
	// mysqlDSN is the DSN for the MySQL database
	mysqlDSN string
	imp      *repository.ImplStorageProductMySQL
}

// Run runs the server
func (s *ServerChi) Run() (err error) {

	// - database: connection
	db, err := sql.Open("mysql", s.mysqlDSN)
	if err != nil {
		return
	}
	defer db.Close()
	// - database: ping
	err = db.Ping()
	if err != nil {
		return
	}
	s.imp = repository.NewImplStorageProductMySQL(db)

	result, err := s.imp.GetOne(1)
	fmt.Println(result)

	return
}

func main() {
	addrCfg := ":8080"
	mysqlCfg := mysql.Config{
		User:      "user1",
		Passwd:    "secret_password",
		Net:       "tcp",
		Addr:      "localhost:3306",
		DBName:    "my_db",
		ParseTime: true,
	}
	cfg := ConfigServerChi{Addr: addrCfg, MySQLDSN: mysqlCfg.FormatDSN()}
	// - server
	server := NewServerChi(cfg)
	// - run
	if err := server.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
