// @APIVersion 1.0.0
// @Title mobile API
// @Description mobile has every tool to get any job done, so codename for the new mobile APIs.
// @Contact astaxie@gmail.com
package routers

import (
	"cz400/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/weight", &controllers.WeightController{})
	beego.Router("/online", &controllers.OnlineController{})
	beego.Router("/params", &controllers.ParamController{})
	beego.Router("/gps", &controllers.GpsController{})
	beego.Router("/upload", &controllers.UploadController{})
	beego.Router("/uploadDevice", &controllers.UploadDevController{})
	beego.Router("/login", &controllers.LoginController{})

	beego.AutoRouter(&controllers.RController{})
}
