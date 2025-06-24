package database

import (
	"database/sql"
	"fmt"
	"framework/util/ini"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// Config 配置信息
type Config struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
	Charset  string
	MaxOpen  int
	MaxIdle  int
}

func NewMysql() {
	port, _ := strconv.Atoi(ini.Iniconfig["mysql.port"])
	cfg := Config{
		User:     ini.Iniconfig["mysql.user"],
		Password: ini.Iniconfig["mysql.password"],
		Host:     ini.Iniconfig["mysql.addr"],
		Port:     port,
		DBName:   ini.Iniconfig["mysql.dbname"],
		Charset:  "utf8mb4",
		MaxOpen:  10,
		MaxIdle:  5,
	}

	if err := Init(cfg); err != nil {
		log.Fatal("数据库初始化失败:", err)
	}

}

// Init 初始化数据库连接
func Init(cfg Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Charset)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// 设置连接池参数
	DB.SetMaxOpenConns(cfg.MaxOpen)
	DB.SetMaxIdleConns(cfg.MaxIdle)

	if err := DB.Ping(); err != nil {
		return err
	}

	log.Println("数据库连接成功")
	return nil
}

func QueryRows(query string, args ...any) ([]map[string]any, error) {
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	results := []map[string]any{}
	for rows.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}

		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}

		row := make(map[string]any)
		for i, col := range cols {
			row[col] = vals[i]
		}
		results = append(results, row)
	}

	return results, nil
}

func QueryRow(query string, args ...any) (map[string]any, error) {
	rows, err := QueryRows(query, args...)
	if err != nil || len(rows) == 0 {
		return nil, err
	}
	return rows[0], nil
}
func Exec(query string, args ...any) (int64, error) {
	result, err := DB.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func Query(queryStr string) ([]map[string]string, error) {
	var records []map[string]string
	rows, err := DB.Query(queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		columns, _ := rows.Columns()
		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))

		for i := range values {
			scanArgs[i] = &values[i]
		}

		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		record := make(map[string]string)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		records = append(records, record)
	}
	return records, nil
}

func execSQL(queryStr string) (int64, error) {
	result, err := DB.Exec(queryStr)
	if err != nil {
		return -1, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}
	return rowsAffected, nil
}

func InsertData(insertStr string) (int64, error) {
	return execSQL(insertStr)
}

func UpdateData(UpdateStr string) (int64, error) {
	return execSQL(UpdateStr)
}

func DeleteData(delStr string) (int64, error) {
	return execSQL(delStr)
}

func Close() {
	DB.Close()
}
