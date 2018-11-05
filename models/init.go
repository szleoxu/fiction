package models

import(
	_ "github.com/go-sql-driver/mysql"
	"github.com/Unknwon/goconfig"
	"log"
	"database/sql"
)

func InitBook()(*DBBook,error){
	cfg, errConfig := goconfig.LoadConfigFile("config.ini")
	if errConfig != nil {
		log.Fatalf("无法加载配置文件: %s", errConfig)
	}
	username,_:=cfg.GetValue("mysql","username")
	password,_:=cfg.GetValue("mysql","password")
	host,_:=cfg.GetValue("mysql","host")
	port,_:=cfg.GetValue("mysql","port")
	database,_:=cfg.GetValue("mysql","database")
	connectionStr:=username+":"+password+"@tcp("+host+":"+port+")/"+database

	var err error
	dbw := DBBook{
		Dsn: connectionStr,
	}
	dbw.Db, err = sql.Open("mysql",dbw.Dsn)
	if err != nil {
		return nil,err
	}
	if err := dbw.Db.Ping(); err != nil{
		return nil,err
	}else{
		return &dbw,nil
	}
}


func InitChapter()(*DBChapter,error){
	cfg, errConfig := goconfig.LoadConfigFile("config.ini")
	if errConfig != nil {
		log.Fatalf("无法加载配置文件: %s", errConfig)
	}
	username,_:=cfg.GetValue("mysql","username")
	password,_:=cfg.GetValue("mysql","password")
	host,_:=cfg.GetValue("mysql","host")
	port,_:=cfg.GetValue("mysql","port")
	database,_:=cfg.GetValue("mysql","database")
	connectionStr:=username+":"+password+"@tcp("+host+":"+port+")/"+database

	var err error
	dbw := DBChapter{
		Dsn: connectionStr,
	}
	dbw.Db, err = sql.Open("mysql",dbw.Dsn)
	if err != nil {
		return nil,err
	}
	if err := dbw.Db.Ping(); err != nil{
		return nil,err
	}else{
		return &dbw,nil
	}
}