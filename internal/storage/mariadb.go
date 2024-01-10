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
	var domain string
	var domainReds string
	var domainStart string

	query := "SELECT `domain_name`, `domain_reds`,  `domain_start` FROM `domains` WHERE `domain_locale` = 'ru' AND `domain_type` = 'news' ORDER BY `domain_accept_redirect` DESC, `domain_start` DESC  LIMIT 1"

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

	m.db.SetMaxIdleConns(50)
	m.db.SetMaxOpenConns(250)
	m.db.SetConnMaxLifetime(60 * time.Second)

	// Проверка подключения
	err = m.db.Ping()
	if err != nil {
		config.Log.Panic(err.Error())
		panic(err.Error())
	}

	config.Log.Debug("connected with mariadb")
}

func (m *MariaDBS) PathId(id string) (string, error) {
	var newPath string
	err := m.db.QueryRow("SELECT CONCAT(`cat_url`, '/', `item_url`) FROM `news_posts` LEFT JOIN `news_posts_cat` ON `news_posts_cat`.`cat_id` = `news_posts`.`item_cat` WHERE `item_id` = ?", id).Scan(&newPath)

	return newPath, err
}
