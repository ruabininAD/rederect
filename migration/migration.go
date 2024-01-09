package main

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"rederect/internal/config"
	"sync"
)

var wg sync.WaitGroup

func main() {
	config.Init()

	m := sql.DB{}
	migrationMariadb(&m)
}

func migrationMariadb(db *sql.DB) {

	err := errors.New("")
	// Подключение к базе данных
	dataSourceName :=
		"root:" + config.Cfg.DB.Password + "@tcp(" + config.Cfg.DB.Ip + ":" + config.Cfg.DB.Port + ")/"
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		config.Log.Panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		config.Log.Panic(err.Error())
	}

	config.Log.Debug("connected with mariadb")

	var domain string

	createTableQuestion := "CREATE TABLE iav2.domains ( domain_name varchar(32)  NOT NULL  PRIMARY KEY, domain_type varchar(8) NOT NULL, domain_add             datetime     DEFAULT CURRENT_TIMESTAMP() NOT NULL, domain_start datetime NULL, domain_reds int(10) DEFAULT 0 NOT NULL, domain_accept_redirect int(1) DEFAULT 0 NOT NULL,domain_parked int(1) DEFAULT 0 NOT NULL,  domain_locale varchar(5)   DEFAULT 'ru' NOT NULL,domain_note varchar(300) DEFAULT '' NOT NULL COMMENT 'Заметка', expire_domain timestamp NULL COMMENT 'Окончание домена',expire_ssl timestamp NULL COMMENT 'Окончание сертификата' ); "
	domainIndexQuestion := "CREATE INDEX expire_domain_index ON iav2.domains (expire_domain);"
	sslIndexQuestion := "CREATE INDEX expire_ssl_index   ON iav2.domains (expire_ssl);"

	Questions := []string{
		"create database iav2;",
		createTableQuestion,
		domainIndexQuestion,
		sslIndexQuestion,
		"INSERT INTO iav2.domains (domain_name, domain_type, domain_add, domain_start, domain_reds, domain_accept_redirect, domain_parked, domain_locale, domain_note, expire_domain, expire_ssl) VALUES ('morelo.space', 'preland', '2023-06-14 07:35:33', null, 0, 0, 0, 'th', '', '2023-07-16 02:59:59', null);",
		" INSERT INTO iav2.domains (domain_name, domain_type, domain_add, domain_start, domain_reds, domain_accept_redirect, domain_parked, domain_locale, domain_note, expire_domain, expire_ssl) VALUES ('hot-news.local', 'news', '2023-06-14 07:27:30', null, 1, 1, 0, 'ru', 'Для smi2', '2023-10-30 11:38:58', '2023-07-18 23:53:22');",
		" INSERT INTO iav2.domains (domain_name, domain_type, domain_add, domain_start, domain_reds, domain_accept_redirect, domain_parked, domain_locale, domain_note, expire_domain, expire_ssl) VALUES ('news.local', 'news', '2023-06-14 07:27:44', null, 1, 1, 0, 'ru', '', null, null); ",
		"INSERT INTO iav2.domains (domain_name, domain_type, domain_add, domain_start, domain_reds, domain_accept_redirect, domain_parked, domain_locale, domain_note, expire_domain, expire_ssl) VALUES ('test.local', 'system', '2023-06-14 13:20:06', null, 0, 0, 0, 'ru', '', null, null);",
		" INSERT INTO iav2.domains (domain_name, domain_type, domain_add, domain_start, domain_reds, domain_accept_redirect, domain_parked, domain_locale, domain_note, expire_domain, expire_ssl) VALUES ('thnews.local', 'news', '2023-06-14 07:27:44', null, 1, 1, 0, 'th', '', null, null);",
	}

	for _, question := range Questions {
		err = db.QueryRow(question).Scan(&domain)
		if err != nil {
			config.Log.Info(err.Error())
		}

	}

}
