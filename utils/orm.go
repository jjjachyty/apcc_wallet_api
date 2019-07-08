package utils

import (
	"os"
	"sync"
	"time"

	// "time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

//DBEngine var 数据库引擎
var dbEngine *xorm.Engine

var waitgroup sync.WaitGroup

//GSession 全局Session
var GSessions sync.Map //= make( [string]*xorm.Session, 0)

// type UsingSync struct {
// 	sync.RWMutex
// 	Count int
// }

type GSession struct {
	SID string
	// GID      string
	RollBack bool //关闭是否回滚 否就commit
	*xorm.Session
}

//OpenSession 初始化Session Service层使用
func OpenSession() GSession {
	var gid = GetGID()
	session, ok := GSessions.Load(gid)
	if !ok {
		SysLog.Debug("打开session-不存在Session,新开Session")
		session = GSession{SID: gid, Session: dbEngine.NewSession()}
		GSessions.Store(gid, session)
	}
	SysLog.Debug("打开session-已存在Session,返回session")
	return session.(GSession)
}

//OpenSession 初始化Session Service层使用 SID 用户SID号 rollback 事物关闭是否需要回滚 false 为 commit
func OpenSessionAdvanced(SID string, rollback bool) GSession {
	var gsession GSession
	SysLog.Debug("打开Session高级方法")
	session, ok := GSessions.Load(SID)
	if !ok {
		SysLog.Debug("打开session-不存在Session,新开Session")
		gsession = GSession{SID: SID, RollBack: rollback, Session: dbEngine.NewSession()}

	} else {
		SysLog.Debug("打开session-已存在Session,返回session")
		gsession = session.(GSession)
		if gsession.RollBack {
			gsession.Rollback()
		}
		gsession.Close()
		gsession.Session = dbEngine.NewSession()
	}
	GSessions.Store(SID, gsession)

	return gsession
}

//GetSession 获取当前协程的Session 适用与Models
func GetSession(SID ...string) (GSession, bool) {
	SysLog.Debug("获取Session")

	var gid string
	if nil != SID {
		gid = SID[0]
	} else {
		gid = GetGID()
	}
	session, ok := GSessions.Load(gid)
	if !ok {
		//session =
		//GSessions.Store(gid, session)
		return GSession{SID: gid, Session: dbEngine.NewSession()}, true //Service 层未开启了Session
	}
	gsession := session.(GSession)
	return gsession, false //Service 层已经开启了Session
}

func GetDBEngine() *xorm.Engine {
	return dbEngine
}

//CloseSession 释放掉Session
func (gs *GSession) CloseSession() {
	if !gs.IsClosed() {
		sessionStore, ok := GSessions.Load(gs.SID)
		if ok {
			gsession := sessionStore.(GSession)

			if gsession.RollBack {
				SysLog.Debug("GSession方法关闭session-事物回滚")
				gsession.Rollback()
			}

			GSessions.Delete(gsession.SID)
		}
		gs.Close()
	}
}

//CloseSession 释放掉Session
func CloseSession() {
	session, _ := GSessions.Load(GetGID())
	session.(GSession).Close()
	GSessions.Delete(GetGID())
	SysLog.Debug("方法关闭session-关闭已存在的Session")

}

//CloseSession 释放掉Session
func CloseSessionGID(gid string) {
	session, ok := GSessions.Load(gid)
	if ok {
		session.(GSession).Close()
		GSessions.Delete(gid)
		SysLog.Debug("根据gid方法关闭session-关闭已存在的Session")
	}
}

func Session(fc func(sessin *xorm.Session) error) error {
	session := dbEngine.NewSession()

	defer session.Close()
	return fc(session)
}

//初始化数据库引擎
func InitDB() {
	var err error
	var dataSrouceName = ""
	dbcfg := GetDataBaseCfg()
	if "mysql" == dbcfg.Type {

		dataSrouceName = dbcfg.UserName + ":" + dbcfg.PassWord + "@(" + dbcfg.Server + ":" + dbcfg.Port + ")/" + dbcfg.SID + "?charset=utf8mb4"
	}
	SysLog.Debug("driverName--", dbcfg.DriverName, "--dataSrouceName--", dataSrouceName)
	dbEngine, err = xorm.NewEngine(dbcfg.DriverName, dataSrouceName)
	if nil != err {
		SysLog.Error("初始化数据库出错", err)
		os.Exit(-1)
	}
	dbEngine.ShowSQL(true)

	// tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, "lrm_")
	// dbEngine.SetTableMapper(tbMapper)
	cqtz, err := time.LoadLocation("Asia/Chongqing")
	if nil != err {
		SysLog.Error("设置时区Asia/Chongqing错误", err)
	}
	dbEngine.SetTZDatabase(cqtz)
	dbEngine.SetTZLocation(cqtz)

	// sl := xorm.NewSimpleLogger(SysLog.Writer())
	// sl.SetLevel(core.LOG_INFO)
	// dbEngine.SetLogger(sl)

	// dbEngine.Ping()
	// fmt.Println(result)
}
