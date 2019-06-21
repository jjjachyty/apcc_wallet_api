package models

import (
	"apcc_wallet/utils"
	"bytes"
	"net/url"
	"strings"
)

//QueryByPage 分页查询
func QueryByPage(pageData *utils.PageData, bean interface{}, results interface{}, orderBy string, sort string, whereSQL string, args ...interface{}) error {

	err := getResultCount(bean, pageData, whereSQL, args...)
	if nil != err {
		return err
	}
	err = getPageResult(results, pageData, orderBy, sort, whereSQL, args...)

	return err
}

//Query 实体查询
func Query(bean interface{}, results interface{}, orderBy string, sort string, whereSQL string, args ...interface{}) error {

	err := getResult(results, orderBy, sort, whereSQL, args...)
	return err
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
func Create(beanP interface{}) error {
	session, flag := utils.GetSession()
	_, err := session.Insert(beanP)
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

func getPageResult(results interface{}, pageData *utils.PageData, orderBy string, sort string, whereSQL string, args ...interface{}) error {
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
	err := session.Limit(pageData.Page.PageSize, pageData.Page.PageSize*(pageData.Page.CurrentPage-1)).Find(results)
	if flag {
		session.Close()
	}
	if nil != err {
		utils.SysLog.Error("数据库查询出错", err)
		return err
	}
	return nil
}

func getResultCount(bean interface{}, pageData *utils.PageData, whereSQL string, args ...interface{}) error {
	session, flag := utils.GetSession()
	if nil != args {
		session.Where(whereSQL, args...)
	}
	total, err := session.Count(bean)
	if flag {
		session.Close()
	}
	if nil != err {
		utils.SysLog.Error("数据库查询出错", err)
		return err
	}
	pageData.Page.TotalRows = total
	pageSize64 := int64(pageData.Page.PageSize)

	pageData.Page.PageCount = int(pageData.Page.TotalRows % pageSize64)

	return nil
}

func GetSQL(pars url.Values, filedSQL map[string]string) (string, []interface{}, string, string) {
	var whereSQL string
	var args []interface{}
	var AndFlag bool = false
	var orderBy string
	var sort string

	delete(pars, "Sid")
	delete(pars, "UserId")

	if nil != pars["sort"] {
		sort = pars["sort"][0]
	}
	if nil != pars["orderBy"] {
		orderBy = Field2Cols(pars["orderBy"][0])
	}

	delete(pars, "orderBy")
	delete(pars, "sort")

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
	return whereSQL, args, orderBy, sort
}

func Field2Cols(filed string) string {
	byts := []byte(filed)
	col := bytes.NewBuffer(nil)
	for i, v := range byts {
		if (int(v) > 64 && int(v) < 91) && i > 0 { //大写字母
			col.WriteByte('_')
		}
		col.WriteByte(v)
	}
	return col.String()
}
