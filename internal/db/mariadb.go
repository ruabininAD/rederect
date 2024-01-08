package db

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"rederect/internal/config"
)

type MariaDBS struct {
	db *sql.DB
}

func (m *MariaDBS) GetLast() (string, error) {
	var domain string
	query := "SELECT `domain_name` FROM `domains` WHERE `domain_locale` = 'ru' AND `domain_type` = 'news' ORDER BY `domain_accept_redirect` DESC, `domain_start` DESC LIMIT 1"
	err := m.db.QueryRow(query).Scan(&domain)
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
	// Подключение к базе данных
	dataSourceName :=
		"root:" + config.Cfg.DB.Password +
			"@tcp(" + config.Cfg.DB.Ip + ":" + config.Cfg.DB.Port + ")" +
			"/" + config.Cfg.DB.NameDB
	m.db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		config.Log.Panic(err.Error())
		panic(err.Error())
	}
	// Проверка подключения
	err = m.db.Ping()
	if err != nil {
		config.Log.Panic(err.Error())
		panic(err.Error())
	}

	config.Log.Debug("connected with mariadb")
}
