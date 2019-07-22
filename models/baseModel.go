package models

import (
	"apcc_wallet_api/utils"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

//Query 实体查询
func Query(bean interface{}, results interface{}, orderBy string, sort string, whereSQL string, args ...interface{}) error {

	err := getResult(results, orderBy, sort, whereSQL, args...)
	return err
}

//Query 实体查询
func QueryOne(results interface{}, whereSQL string, args ...interface{}) (bool, error) {
	return getOneResult(results, whereSQL, args...)
}

//GetBean 根据实体非空字段查询单跳数据
func GetBean(beanP interface{}) error {
	session, flag := utils.GetSession()
	if flag {
		session.Close()
	}
	_, err := session.Get(beanP)
	return err
}

//GetBean 根据实体非空字段查询单跳数据
func GetBeans(beans interface{}, condBean interface{}) error {
	session, flag := utils.GetSession()
	if flag {
		session.Close()
	}

	return session.Find(beans, condBean)
}

//GetBean 根据实体非空字段查询单跳数据
func SQLBean(beanP interface{}, sql string, args ...interface{}) error {
	session, flag := utils.GetSession()
	if flag {
		session.Close()
	}
	_, err := session.SQL(sql, args...).Get(beanP)
	return err
}

//GetBean 根据实体非空字段查询单跳数据
func SQLBeans(beans interface{}, sql string, args ...interface{}) error {
	session, flag := utils.GetSession()
	if flag {
		session.Close()
	}
	return session.SQL(sql, args...).Find(beans)
}

//Update 更新表数据
func Update(UUID string, beanP interface{}) error {
	session, flag := utils.GetSession()
	_, err := session.Id(UUID).Update(beanP)
	if flag {
		session.Close()
	}
	return err
}

//UpdateBean 更新表数据beanP 更新的数据 whereBean更新的条件 &User{Name:name} UPDATE user SET ... Where name = ?
func UpdateBean(beanP interface{}, whereBean interface{}) error {
	session, flag := utils.GetSession()
	_, err := session.Update(beanP, whereBean)
	if flag {
		session.Close()
	}
	return err
}
func Create(bean interface{}) error {
	session, flag := utils.GetSession()
	_, err := session.Insert(bean)
	if flag {
		session.Close()
	}
	return err
}

//Update 更新表数据
func Delete(UUID string, beanP interface{}) error {
	session, flag := utils.GetSession()
	_, err := session.Id(UUID).Delete(beanP)
	if flag {
		session.Close()
	}
	return err
}

func getResult(results interface{}, orderBy string, sort string, whereSQL string, args ...interface{}) error {
	session, flag := utils.GetSession()

	if nil != args {
		session.Where(whereSQL, args...)
	}
	if "" != orderBy && "" != sort {
		switch sort {
		case "desc":
			session.Desc(orderBy)
		case "asc":
			session.Asc(orderBy)
		}
	}
	err := session.Find(results)
	if flag {
		session.Close()
	}
	if nil != err {
		utils.SysLog.Error("数据库查询出错", err)
		return err
	}
	return nil
}

func GetBeansPage(pageData *utils.PageData, condBean interface{}) error {
	session, flag := utils.GetSession()
	beanType := reflect.TypeOf(condBean)
	slice := reflect.MakeSlice(reflect.SliceOf(beanType), 0, 15)
	result := reflect.New(slice.Type())
	result.Elem().Set(slice)

	if "" != pageData.Page.OrderBy && "" != pageData.Page.Sort {
		switch pageData.Page.Sort {
		case "desc":
			session.Desc(pageData.Page.OrderBy)
		case "asc":
			session.Asc(pageData.Page.OrderBy)
		}
	}
	total, err := session.Limit(pageData.Page.PageSize, pageData.Page.PageSize*(pageData.Page.CurrentPage-1)).FindAndCount(result.Interface(), condBean)
	if flag {
		session.Close()
	}
	pageData.Page.TotalRows = total
	pageSize64 := int64(pageData.Page.PageSize)

	pageData.Page.PageCount = int(pageData.Page.TotalRows % pageSize64)
	pageData.Rows = result.Interface()
	if nil != err {
		utils.SysLog.Error("数据库查询出错", err)
		return err
	}
	return nil
}

func GetSQLPage(pageData *utils.PageData, where string, args ...interface{}) error {
	session, flag := utils.GetSession()

	if "" != pageData.Page.OrderBy && "" != pageData.Page.Sort {
		switch pageData.Page.Sort {
		case "desc":
			session.Desc(pageData.Page.OrderBy)
		case "asc":
			session.Asc(pageData.Page.OrderBy)
		}
	}
	fmt.Println("----------------------", reflect.TypeOf(pageData.Rows))
	total, err := session.Limit(pageData.Page.PageSize, pageData.Page.PageSize*(pageData.Page.CurrentPage-1)).Where(where, args...).FindAndCount(pageData.Rows)
	if flag {
		session.Close()
	}
	pageData.Page.TotalRows = total
	pageSize64 := int64(pageData.Page.PageSize)

	pageData.Page.PageCount = int(pageData.Page.TotalRows % pageSize64)

	if nil != err {
		utils.SysLog.Error("数据库查询出错", err)
		return err
	}
	return nil
}

func getOneResult(results interface{}, whereSQL string, args ...interface{}) (bool, error) {
	session, flag := utils.GetSession()

	if nil != args {
		session.Where(whereSQL, args...)
	}

	has, err := session.Get(results)
	if flag {
		session.Close()
	}
	if nil != err {
		utils.SysLog.Error("数据库查询出错", err)
		return has, err
	}
	return has, err
}

func GetSQL(pars url.Values, filedSQL map[string]string) (string, []interface{}) {
	var whereSQL string
	var args []interface{}
	var AndFlag = false
	// var orderBy string
	// var sort string

	delete(pars, "order")
	delete(pars, "sort")
	delete(pars, "page")
	for k, v := range pars {
		if AndFlag {
			whereSQL += " AND "
		}
		//通用查询Search 处理
		if utils.WEB_COMMON_QUERY_KEY == k {
			if 0 < len(v) && "" != v[0] {
				//根据占位符？的个数区分需要多少个参数
				argsCount := strings.Count(filedSQL[k], "?")
				for index := 0; index < argsCount; index++ {
					args = append(args, v[0])
				}
				whereSQL += filedSQL[k]
			}
			continue
		}

		if len(v) < 1 {
			args = append(args, v)
			AndFlag = true
			whereSQL += filedSQL[k]
		} else { //一个字段多个值查询
			var orsql = "("
			var hasor = false
			for _, val := range v {
				if hasor {
					orsql += " OR "
				}
				orsql = orsql + filedSQL[k]
				args = append(args, val)
				hasor = true
			}
			orsql += " )"
			AndFlag = true
			whereSQL += orsql
		}

	}
	return whereSQL, args
}
