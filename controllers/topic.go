package controllers

import (
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

}

func (this *TopicController) Add() {
	this.TplName = "topic_add.html"
}
