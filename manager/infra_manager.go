package manager

import (
	"database/sql"
	"fmt"
	"polen/config"

	_ "github.com/lib/pq"
)

type InfraManager interface {
	Conn() *sql.DB
}

type infraManager struct {
	db  *sql.DB
	cfg *config.Config
}

func (d *infraManager) initDb() error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", d.cfg.DbConfig.Host, d.cfg.DbConfig.Port, d.cfg.DbConfig.User, d.cfg.DbConfig.Password, d.cfg.DbConfig.Name)

	db, err := sql.Open(d.cfg.DbConfig.Driver, dsn)
	if err != nil {
		return err
	}

	d.db = db
	return nil
}

func (d *infraManager) Conn() *sql.DB {
	return d.db
}

// constructor
func NewInfraManager(cfg *config.Config) (InfraManager, error) {
	conn := &infraManager{
		cfg: cfg,
	}
	err := conn.initDb()
	if err != nil {
		return nil, err
	}
	return conn, nil
}
