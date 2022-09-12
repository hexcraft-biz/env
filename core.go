package env

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultDBType = "mysql"
)

type Prototype struct {
	// Application Setup
	AppHost    string
	AppPath    string
	AppPort    string
	GinMode    string
	Location   *time.Location
	TrustProxy string

	// Database
	AutoCreateDBSchema bool
	DBType             string
	DBInitUser         string
	DBInitPassword     string
	DBInitParams       string
	DBUser             string
	DBPassword         string
	DBHost             string
	DBPort             string
	DBName             string
	DBParams           string
	DBMaxOpen          int
	DBMaxIdle          int
	DBLifeTime         int
	DBIdleTime         int
}

func Fetch() (*Prototype, error) {
	var (
		err                error
		appPath            string
		loc                *time.Location
		autoCreateDBSchema bool
		dbType             string
	)

	if appPath = os.Getenv("APP_PATH"); !strings.HasPrefix(appPath, "/") {
		appPath = "/" + appPath
	}

	if loc, err = time.LoadLocation(os.Getenv("TIMEZONE")); err != nil {
		return nil, err
	}

	if autoCreateDBSchema, err = strconv.ParseBool(os.Getenv("AUTO_CREATE_DB_SCHEMA")); err != nil {
		return nil, err
	}

	dbType = os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = DefaultDBType
	}

	e := &Prototype{
		AppHost:            os.Getenv("APP_HOST"),
		AppPath:            appPath,
		AppPort:            os.Getenv("APP_PORT"),
		GinMode:            os.Getenv("GIN_MODE"),
		Location:           loc,
		TrustProxy:         os.Getenv("TRUST_PROXY"),
		AutoCreateDBSchema: autoCreateDBSchema,
		DBType:             dbType,
	}

	switch e.DBType {
	case "mysql":
		err = e.fetchEnvMysql()
	default:
	}

	return e, err
}

func (e *Prototype) fetchEnvMysql() error {
	var (
		err        error
		dbMaxOpen  int
		dbMaxIdle  int
		dbLifeTime int
		dbIdleTime int
	)

	if dbMaxOpen, err = strconv.Atoi(os.Getenv("DB_MAX_OPEN")); err != nil {
		return err
	}

	if dbMaxIdle, err = strconv.Atoi(os.Getenv("DB_MAX_IDLE")); err != nil {
		return err
	}

	if dbLifeTime, err = strconv.Atoi(os.Getenv("DB_LIFE_TIME")); err != nil {
		return err
	}

	if dbIdleTime, err = strconv.Atoi(os.Getenv("DB_IDLE_TIME")); err != nil {
		return err
	}

	e.DBInitUser = os.Getenv("DB_INIT_USER")
	e.DBInitPassword = os.Getenv("DB_INIT_PASSWORD")
	e.DBInitParams = os.Getenv("DB_INIT_PARAMS")
	e.DBUser = os.Getenv("DB_USER")
	e.DBPassword = os.Getenv("DB_PASSWORD")
	e.DBHost = os.Getenv("DB_HOST")
	e.DBPort = os.Getenv("DB_PORT")
	e.DBName = os.Getenv("DB_NAME")
	e.DBParams = os.Getenv("DB_PARAMS")
	e.DBMaxOpen = dbMaxOpen
	e.DBMaxIdle = dbMaxIdle
	e.DBLifeTime = dbLifeTime
	e.DBIdleTime = dbIdleTime

	return nil
}

func (e *Prototype) GetAppRootURL() string {
	return "https://" + e.AppHost + e.AppPath
}

func (e *Prototype) MysqlConnectWithMode(init bool) (*sqlx.DB, error) {
	if init {
		return e.MysqlConnect(
			e.DBType,
			e.DBInitUser,
			e.DBInitPassword,
			e.DBHost,
			e.DBPort,
			"",
			e.DBInitParams,
			1,
			1,
			30,
			30,
		)
	} else {
		return e.MysqlConnect(
			e.DBType,
			e.DBUser,
			e.DBPassword,
			e.DBHost,
			e.DBPort,
			e.DBName,
			e.DBParams,
			e.DBMaxOpen,
			e.DBMaxIdle,
			e.DBLifeTime,
			e.DBIdleTime,
		)
	}
}

func (e *Prototype) MysqlConnect(typ, user, password, host, port, name, params string, maxOpen, maxIdle, lifeTime, idleTime int) (*sqlx.DB, error) {
	if db, err := sqlx.Open(typ, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", user, password, host, port, name, params)); err != nil {
		return nil, err
	} else {
		db.SetMaxOpenConns(maxOpen)
		db.SetMaxIdleConns(maxIdle)
		db.SetConnMaxLifetime(time.Duration(lifeTime) * time.Second)
		db.SetConnMaxIdleTime(time.Duration(idleTime) * time.Second)
		return db, nil
	}
}

func (e *Prototype) MysqlDBInit(sqlDir string, sortedFiles []string) error {
	db, err := e.MysqlConnectWithMode(true)
	defer db.Close()
	if err != nil {
		return err
	}

	if _, err := db.Exec(`CREATE DATABASE IF NOT EXISTS ` + e.DBName + ` COLLATE 'utf8mb4_unicode_ci' CHARACTER SET 'utf8mb4';`); err != nil {
		return err
	} else {
		db.Exec(`USE ` + e.DBName + `;`)
	}

	if len(sortedFiles) > 0 {
		for i := range sortedFiles {
			if err = execSqlFromFile(db, sqlDir+sortedFiles[i]); err != nil {
				break
			}
		}
	} else {
		err = filepath.Walk(sqlDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			} else if info.IsDir() || filepath.Ext(path) != ".sql" {
				return nil
			}

			return execSqlFromFile(db, path)
		})
	}

	return err
}

func execSqlFromFile(db *sqlx.DB, path string) error {
	if data, err := os.ReadFile(path); err != nil {
		return err
	} else if _, err := db.Exec(string(data)); err != nil {
		return err
	}

	return nil
}
