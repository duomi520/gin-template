package database

import (
	"../../../config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var SqlDB *sql.DB

func init() {
	var err error
	SqlDB, err = sql.Open("mysql", config.GetEnv().DATABASE_USERNAME+
		":"+config.GetEnv().DATABASE_PASSWORD+"@tcp("+config.GetEnv().DATABASE_IP+
		":"+config.GetEnv().DATABASE_PORT+")/"+config.GetEnv().DATABASE_NAME)
	if err != nil {
		log.Fatal(err.Error())
		panic(err.Error())
	}
	err = SqlDB.Ping()
	if err != nil {
		log.Fatal(err.Error())
		panic(err.Error())
	}

	// 设置数据库最大连接 减少timewait
	SqlDB.SetMaxIdleConns(2000)
	SqlDB.SetMaxOpenConns(2000)

	if err := SqlDB.Ping(); err != nil {
		log.Fatalln(err)
		panic(err)
	}
}

func Query(query string, args ...interface{}) ([]map[string]interface{}, *sql.Rows) {
	rs, err := SqlDB.Query(query, args...)

	if err != nil {
		log.Fatalln(err)
		panic(err)
	}

	col, colErr := rs.Columns()

	if colErr != nil {
		log.Fatalln(colErr)
		panic(colErr)
	}

	typeVal, err := rs.ColumnTypes()
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}

	results := make([]map[string]interface{}, 0)

	for rs.Next() {
		var colVar = make([]interface{}, len(col))
		for i := 0; i < len(col); i++ {
			//log.Println("type: " + typeVal[i].ScanType().Name())
			// TODO: string类型如果为interface返回乱字符串
			if typeVal[i].ScanType().Name() == "RawBytes" {
				var s string
				colVar[i] = &s
			} else {
				var s interface{}
				colVar[i] = &s
			}
		}
		result := make(map[string]interface{})
		if scanErr := rs.Scan(colVar...); scanErr != nil {
			log.Fatalln(scanErr)
			panic(scanErr)
		}
		for j := 0; j < len(col); j++ {
			result[col[j]] = colVar[j]
		}
		results = append(results, result)
	}
	if err := rs.Err(); err != nil {
		log.Fatalln(err)
		panic(err)
	}
	return results, rs
}

func Exec(query string, args ...interface{}) sql.Result {
	rs, err := SqlDB.Exec(query, args...)
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
	return rs
}
