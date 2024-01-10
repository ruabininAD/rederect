package storage

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"redirect/internal/config"
	"time"
)

type MariaDBS struct {
	db *sql.DB
}

func (m *MariaDBS) GetLast() (string, error) {

	query := "SELECT `domain_name`, `domain_reds`,  `domain_start` FROM `domains` WHERE `domain_locale` = 'ru' AND `domain_type` = 'news' ORDER BY `domain_accept_redirect` DESC, `domain_start` DESC  LIMIT 1"

	var domain, domainReds, domainStart string
	err := m.db.QueryRow(query).Scan(&domain, &domainReds, &domainStart)
	if err != nil {
		return "", err
	}

	return domain, nil

}
func (m *MariaDBS) Update(domain string) error {
	query := `INSERT INTO domains (domain_name, domain_add, domain_start, domain_reds, domain_type) 
			  VALUES (?, NOW(), NOW(), 1, 'news')
			  ON DUPLICATE KEY UPDATE 
			  domain_start = IFNULL(domain_start, NOW()), domain_reds = domain_reds + 1`
	_, err := m.db.Exec(query, domain)
	if err != nil {
		return err
	}
	return nil
}
func (m *MariaDBS) Connect() {
	err := errors.New("")

	m.db, err = sql.Open("mysql", config.Cfg.DspToDatabase)
	if err != nil {
		config.Log.Panic(err.Error())
		panic(err.Error())
	}

	err = m.db.Ping()
	if err != nil {
		config.Log.Panic(err.Error())
		panic(err.Error())
	}

	m.db.SetMaxIdleConns(50)
	m.db.SetMaxOpenConns(250)
	m.db.SetConnMaxLifetime(60 * time.Second)

	config.Log.Debug("connected with mariadb")
}
