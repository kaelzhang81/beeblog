package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"

	c.Data["TrueCond"] = true
	c.Data["FalseCond"] = false

	type u struct {
		Name string
		Age  int
		Sex  string
	}

	user := &u{
		Name: "lily",
		Age:  18,
		Sex:  "female",
	}

	c.Data["user"] = user

	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	c.Data["nums"] = nums

	c.Data["tplval"] = "hello world"

	c.Data["html"] = "<div>hello beego</div>"

	c.Data["pipe"] = "<div>hello beego</div>"
}
