package controllers

import (
	"beeblog/models"

	"github.com/astaxie/beego"
)

type ReplyControllers struct {
	beego.Controller
}

func (this *ReplyControllers) Add() {
	tid := this.Input().Get("tid")
	err := models.AddReply(tid,
		this.Input().Get("nickname"),
		this.Input().Get("content"))

	if err != nil {
		beego.Error(err)
	}

	this.Redirect("/topic/views/" + tid)
}
