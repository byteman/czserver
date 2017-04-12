package controllers

import (
	"encoding/json"
	"fmt"
	"cz400/czserver"
	"cz400/models"
	"path"
	con "strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type MainController struct {
	beego.Controller
}
type WeightController struct {
	beego.Controller
}

type OnlineController struct {
	beego.Controller
}
type ParamController struct {
	beego.Controller
}
type GpsController struct {
	beego.Controller
}
type UploadController struct {
	beego.Controller
}
type UploadDevController struct {
	beego.Controller
}
type GpsListController struct {
	beego.Controller
}
type LoginController struct {
	beego.Controller
}
type UserLogin struct {
	Name string
	Pwd string
}
//登陆逻辑处理
func (c *LoginController) Post(){
	fmt.Println("Login")
	res := make(map[string]interface{})
	result := -1
	role:=0
	message := "ok"
	fmt.Println(c.Ctx.Input.RequestBody)
	u:=UserLogin{}
	defer func() {
		res["result"] = result
		res["role"] = role
		res["message"] = message
		c.Data["json"] = res
		c.ServeJSON()
	}()
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &u)
	if err != nil {
		fmt.Println(err)
		result=-2
		return
	}
	if u.Name=="admin" {
		if u.Pwd == "123456"{
			role=1
			result=0
		}else if u.Pwd == "123321"{
			role=2
			result=0
		}
	}
	fmt.Printf("%+v", u)
	return
}
func (c *UploadDevController) Post() {

}

func (c *MainController) Get() {
	c.Redirect("static/html/index.html", 301)
}

type JsonGpsData struct {
	Error    string
	Result   uint32
	Total    uint32
	PageSize uint32
	Gps      []models.Gps
}

func (c *GpsController) Get() {

	fmt.Println("weight reqeust")
	var cond = ""
	page, err := c.GetInt("pages")
	if err != nil {
		fmt.Println(err)
		page = 1
	}
	id, err := c.GetInt("id")
	if err == nil {

		cond = fmt.Sprintf(" and dev_id=%d", id)
	}

	fmt.Println("page=", page, "id=", id)

	all, pagesize, rs := GetPagesInfo("gps", page, 10, cond)

	gps := make([]models.Gps, 0)

	rs.QueryRows(&gps)

	res := JsonGpsData{}

	res.Error = "ok"
	res.Result = 0
	res.Total = uint32(all)
	res.PageSize = uint32(pagesize)
	res.Gps = gps
	fmt.Println(all, pagesize, gps)
	c.Data["json"] = &res

	c.ServeJSON()
}

type JsonData struct {
	Error    string
	Result   uint32
	Total    uint32
	PageSize uint32
	Weights  []models.OneWeight
}

/**
 * 分页函数，适用任何表
 * 返回 总记录条数,总页数,以及当前请求的数据RawSeter,调用中需要"rs.QueryRows(&tblog)"就行了  --tblog是一个Tb_log对象
 * 参数：表名，当前页数，页面大小，条件（查询条件,格式为 " and name='zhifeiya' and age=12 "）
 */
func GetPagesInfo(tableName string, currentpage int, pagesize int, conditions string) (int, int, orm.RawSeter) {
	if currentpage <= 1 {
		currentpage = 1
	}
	if pagesize == 0 {
		pagesize = 20
	}
	var rs orm.RawSeter
	o := orm.NewOrm()
	var totalItem, totalpages int = 0, 0 //总条数,总页数
	sql := "SELECT count(*) FROM " + tableName + "  where 1>0 " + conditions
	o.Raw(sql).QueryRow(&totalItem) //获取总条数
	fmt.Println(sql, totalItem)
	if totalItem <= pagesize {
		totalpages = 1
	} else if totalItem > pagesize {
		temp := totalItem / pagesize
		if (totalItem % pagesize) != 0 {
			temp = temp + 1
		}
		totalpages = temp
	}
	sql = "select *  from  " + tableName + " where 1>0" + conditions + " order by id desc " + " LIMIT " + con.Itoa((currentpage-1)*pagesize) + "," + con.Itoa(pagesize)
	fmt.Println(sql)
	rs = o.Raw(sql)
	return totalItem, totalpages, rs
}

func (c *WeightController) Get() {
	fmt.Println("weight reqeust")
	var cond = ""
	page, err := c.GetInt("pages")
	if err != nil {
		fmt.Println(err)
		page = 1
	}
	id, err := c.GetInt("id")
	if err == nil {
		cond = fmt.Sprintf(" and dev_id=%d", id)
	}
	fmt.Println("page=", page, "id=", id)

	all, pagesize, rs := GetPagesInfo("one_weight", page, 10, cond)

	//o := orm.NewOrm()
	ws := make([]models.OneWeight, 0)

	//sql := "select * from one_weight order by id desc limit 10"
	//fmt.Println(sql)
	//_, er := o.Raw(sql).QueryRows(&ws)
	rs.QueryRows(&ws)
	res := JsonData{}

	res.Error = "ok"
	res.Result = 0
	res.Total = uint32(all)
	res.PageSize = uint32(pagesize)
	res.Weights = ws
	fmt.Println(all, pagesize, ws)
	c.Data["json"] = &res

	c.ServeJSON()
	return
}

func (c *OnlineController) Get() {

	clients := czserver.GetClient()
	//fmt.Println("client=", clients)
	c.Data["json"] = &clients
	c.ServeJSON()

}

func (c *ParamController) Post() {
	fmt.Println("params post")
	u := czserver.DevInfo{}
	res := make(map[string]interface{})
	result := 0
	message := "ok"
	defer func() {
		res["result"] = result
		res["message"] = message
		c.Data["json"] = res
		c.ServeJSON()
	}()
	fmt.Println(c.Ctx.Input.RequestBody)
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &u)
	if err != nil {
		fmt.Println(err)
		result = 1
		message = "json convert failed"
		return
	}
	fmt.Printf("%+v", u)

	czserver.UpdateParam(&u)

}
func (c *ParamController) Get() {
	fmt.Println("params get")
	clients := czserver.GetClient()
	c.Data["json"] = &clients
	c.ServeJSON()
}

func (c *UploadController) Post() {
	fmt.Println("upload 1")
	f, fh, err := c.GetFile("file")
	fmt.Println("upload 3")
	if err != nil {
		fmt.Println("err:", err)
		c.ServeJSON()
		return
	}
	fmt.Println("upload 2", fh.Filename)
	p := path.Join("file", fh.Filename)
	fmt.Println(p)
	f.Close()
	err = c.SaveToFile("file", "static/upload/"+fh.Filename)
	if err != nil {
		fmt.Println(err)
	}
	c.ServeJSON()
}

type RController struct {
	beego.Controller
}

func (c *RController) Login() {
	fmt.Println("params get")
	fmt.Println(c.Ctx.Input.Param(":id"))
	c.Data["json"] = "{hello:1234}"
	c.ServeJSON()
}
