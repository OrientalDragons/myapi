package myapi

import (
	"database/sql"
	"fmt"
	"strings"

	//mysql driver
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

//Db 返回本机的mysql数据库对象，参数为数据库名称
func Db(dbname string) *sql.DB {
	fmt.Println("create db conn")
	var (
		db    *sql.DB
		acess string
	)
	acess = fmt.Sprintf("root:223024715@@tcp(localhost:3306)/%s?charset=utf8", dbname)
	db, _ = sql.Open("mysql", acess)
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(30)
	db.Ping()
	return db
}

//DBBox 返回带封装方法的mysql db对象
type DBBox struct {
	DB *sql.DB
}

//DBBOX 是DBBox的接口化
type DBBOX interface {
	Query(string) ([]map[string]string, error)
	SecurityQuery(string, ...interface{}) ([]map[string]string, error)
	RunSQL(string) error
	SecurityRunSQL(string, ...interface{}) error
}

//NewDb 获取一个连接
func NewDb(dbname string, acess string) *DBBox {
	fmt.Println("create db conn")
	var (
		db *sql.DB
	)
	if acess == "" {
		acess = fmt.Sprintf("root:223024715@@tcp(localhost:3306)/%s?charset=utf8", dbname)
		db, _ = sql.Open("mysql", acess)
	} else {
		ace := strings.SplitN(acess, " ", 2)
		fmt.Println(ace)
		db, _ = sql.Open(ace[0], ace[1])
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(30)
	db.Ping()

	return &DBBox{db}
}

//Query 查询
func (df *DBBox) Query(dosql string) ([]map[string]string, error) {
	//定义一个result map数组接收真正的结果值
	var result []map[string]string
	rows, err := df.DB.Query(dosql)
	if err != nil {
		return result, err
	}

	//定义参数接收结果
	colunms, _ := rows.Columns()
	values := make([]sql.RawBytes, len(colunms))
	scanArgs := make([]interface{}, len(colunms))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		res := make(map[string]string)
		rows.Scan(scanArgs...)
		for i, col := range values {
			res[colunms[i]] = string(col)
		}
		result = append(result, res)
	}
	rows.Close()
	fmt.Println("查询结束")
	return result, nil
}

//SecurityQuery 参数化查询
func (df *DBBox) SecurityQuery(s string, args ...interface{}) ([]map[string]string, error) {
	//定义一个result map数组接收真正的结果值
	var result []map[string]string
	stmt, err := df.DB.Prepare(s)
	if err != nil {
		return result, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return result, err
	}

	colunms, _ := rows.Columns()
	values := make([]sql.RawBytes, len(colunms))
	scanArgs := make([]interface{}, len(colunms))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		res := make(map[string]string)
		rows.Scan(scanArgs...)
		for i, col := range values {
			res[colunms[i]] = string(col)
		}
		result = append(result, res)
	}
	rows.Close()
	return result, nil
}

//RunSQL 执行事务操作
func (df *DBBox) RunSQL(dosql string) error {
	// Dosql := strings.ReplaceAll(dosql, "\n", " ")
	conn, err := df.DB.Begin()
	if err != nil {
		return err
	}
	result, err := df.DB.Exec(dosql)
	if err != nil {
		conn.Rollback()
		return err
	}
	A, _ := result.LastInsertId()
	B, _ := result.RowsAffected()

	df.ExecLog(nil, [2]int64{A, B})
	err = conn.Commit()
	if err != nil {
		return err
	}
	return nil
}
func typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

//SecurityRunSQL 参数化执行
func (df *DBBox) SecurityRunSQL(s string, args ...interface{}) error {
	stmt, err := df.DB.Prepare(s)
	if err != nil {
		return err
	}
	fmt.Println("xxxx", args)
	fmt.Println(typeof(args[0]))
	result, err := stmt.Exec(args...)
	if err != nil {
		return err
	}

	A, _ := result.LastInsertId()
	B, _ := result.RowsAffected()
	df.ExecLog(nil, [2]int64{A, B})
	stmt.Close()
	return nil
}

//ExecLog 执行log
func (df *DBBox) ExecLog(fn func(...interface{}), args ...interface{}) {
	if fn == nil {
		fmt.Println(args...)
		return
	}
	fn(args...)
}
