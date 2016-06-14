package controllers

import (
    "beeblog/models"
	"github.com/astaxie/beego"
)

type TopicController struct {
	beego.Controller
}

func (this *TopicController) Get() {
	this.Data["IsTopic"] = true
	this.TplName = "topic.html"
}

func (this *TopicController) Post() {
    if !checkAccount(this.Ctx)
    {
        this.Redirect("/login", 302)
        return
    }
    
    title := this.Input().Get("title")
    content := this.Input().Get("content")
    
    var err error
    err = models.AddTopic(title, content)
    if err != nil{
        beego.Error(err)
    }
    
    this.Redirect("/topic", 302)
}

func (this *TopicController) Add() {
	this.TplName = "topic_add.html"
}
