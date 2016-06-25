package controllers

import (
	"io"
	"net/url"
	"os"

	"github.com/astaxie/beego"
)

type AttachmentController struct {
	beego.Controller
}

func (this *AttachmentController) Get() {
	filePath, err := url.QueryUnescape(this.Ctx.Request.RequestURI[1:])
	if err != nil {
		this.Ctx.WriteString(err.Error())
		return
	}

	f, err := os.Open(filePath)
	if err != nil {
		this.Ctx.WriteString(err.Error())
		return
	}

	defer f.Close()

	_, err = io.Copy(this.Ctx.ResponseWriter, f)
	if err != nil {
		this.Ctx.WriteString(err.Error())
	}
}
