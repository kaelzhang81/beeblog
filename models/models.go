package models

import (
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Unknwon/com"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

const (
	// 设置数据库路径
	_DB_NAME = "data/beeblog.db"
	// 设置数据库名称
	_SQLITE3_DRIVER = "sqlite3"
)

// 分类
type Category struct {
	Id              int64
	Title           string
	Created         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	TopicTime       time.Time `orm:"index"`
	TopicCount      int64
	TopicLastUserId int64
}

// 文章
type Topic struct {
	Id              int64
	Uid             int64
	Title           string
	Category        string
	Labels          string
	Content         string `orm:"size(5000)"`
	Attachment      string
	Created         time.Time `orm:"index"`
	Updated         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	Author          string
	ReplyTime       time.Time `orm:"index"`
	ReplyCount      int64
	ReplyLastUserId int64
}

type Comment struct {
	Id      int64
	Tid     int64
	Name    string
	Content string    `orm:"size(1000)"`
	Created time.Time `orm:"index"`
}

func RegisterDB() {
	// 检查数据库文件
	if !com.IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}

	// 注册模型
	orm.RegisterModel(new(Category), new(Topic), new(Comment))
	// 注册驱动（“sqlite3” 属于默认注册，此处代码可省略）
	orm.RegisterDriver(_SQLITE3_DRIVER, orm.DRSqlite)
	// 注册默认数据库
	orm.RegisterDataBase("default", _SQLITE3_DRIVER, _DB_NAME, 10)
}

func AddCategory(name string) error {
	o := orm.NewOrm()

	cate := &Category{
		Title:     name,
		Created:   time.Now(),
		TopicTime: time.Now(),
	}

	// 查询数据
	qs := o.QueryTable("category")
	err := qs.Filter("title", name).One(cate)
	if err == nil {
		return err
	}

	_, err = o.Insert(cate)

	return err
}

func DeleteCategory(id string) error {
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()

	cate := &Category{Id: cid}
	_, err = o.Delete(cate)
	return err
}

func GetAllCategories() ([]*Category, error) {
	o := orm.NewOrm()
	cates := make([]*Category, 0)

	qs := o.QueryTable("category")
	_, err := qs.All(&cates)
	return cates, err
}

func formatLabels(labels string) string {
	return "$" + strings.Join(strings.Split(labels, " "), "#$") + "#"
}

func AddTopic(title, category, labels, content, attachment string) error {
	o := orm.NewOrm()

	topic := &Topic{
		Title:      title,
		Content:    content,
		Category:   category,
		Attachment: attachment,
		Labels:     formatLabels(labels),
		Created:    time.Now(),
		Updated:    time.Now(),
		ReplyTime:  time.Now(),
	}
	_, err := o.Insert(topic)
	if err != nil {
		return err
	}

	return updateCategoryCount(category)
}

func GetTopic(tid string) (*Topic, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}

	o := orm.NewOrm()
	topic := new(Topic)

	qs := o.QueryTable("topic")
	err = qs.Filter("id", tidNum).One(topic)
	if err != nil {
		return nil, err
	}

	topic.Views++
	_, err = o.Update(topic)

	topic.Labels = strings.Replace(strings.Replace(topic.Labels, "#", " ", -1), "$", "", -1)
	return topic, nil
}

func ModifyTopic(tid, title, category, labels, content, attachment string) error {
	var oldCate, oldAttch string
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		oldCate = topic.Category
		oldAttch = topic.Attachment
		topic.Title = title
		topic.Category = category
		topic.Attachment = attachment
		topic.Labels = formatLabels(labels)
		topic.Content = content
		topic.Updated = time.Now()
		_, err = o.Update(topic)
		if err != nil {
			return nil
		}
	}

	if len(oldCate) > 0 {
		err = updateCategoryCount(oldCate)
		if err != nil {
			return err
		}
	}

	if len(oldAttch) > 0 {
		os.Remove(path.Join("attachment", oldAttch))
	}

	return updateCategoryCount(category)
}

func updateCategoryCount(category string) error {
	o := orm.NewOrm()
	cate := new(Category)
	qs := o.QueryTable("category")
	err := qs.Filter("title", category).One(cate)
	if err == nil {
		topics, err := GetAllTopics(category, "", false)
		if err == nil {
			cate.TopicCount = int64(len(topics))
			log.Println(category, cate.TopicCount)
			_, err = o.Update(cate)
		}
	}

	return err
}

func DeleteTopic(tid string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	topic := new(Topic)
	o := orm.NewOrm()
	qs := o.QueryTable("topic")
	err = qs.Filter("id", tidNum).One(topic)
	if err != nil {
		topic.Id = tidNum
	}

	cate := topic.Category
	_, err = o.Delete(topic)
	if err == nil && len(topic.Category) > 0 {
		return updateCategoryCount(cate)
	}

	return err
}

func GetAllTopics(cate, label string, isDesc bool) (topics []*Topic, err error) {
	o := orm.NewOrm()

	topics = make([]*Topic, 0)

	qs := o.QueryTable("topic")
	if len(cate) > 0 {
		qs = qs.Filter("category", cate)
	}
	if len(label) > 0 {
		qs = qs.Filter("labels__contains", "$"+label+"#")
	}
	if isDesc {
		_, err = qs.OrderBy("-created").All(&topics)
	} else {
		_, err = qs.All(&topics)
	}
	return topics, err
}

func AddReply(tid, nickname, content string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	comment := &Comment{
		Tid:     tidNum,
		Name:    nickname,
		Content: content,
		Created: time.Now(),
	}
	o := orm.NewOrm()
	_, err = o.Insert(comment)

	if err == nil {
		topic := &Topic{Id: tidNum}
		if o.Read(topic) == nil {
			topic.ReplyCount++
			topic.ReplyTime = comment.Created
			_, err = o.Update(topic)
		}
	}

	return err
}

func GetAllReplies(tid string) ([]*Comment, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}

	replies := make([]*Comment, 0)
	o := orm.NewOrm()
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tidNum).All(&replies)
	return replies, nil
}

func DeleteReply(rid string) error {
	ridNum, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()

	var tidNum int64
	reply := &Comment{Id: ridNum}
	if o.Read(reply) == nil {
		tidNum = reply.Tid
		_, err = o.Delete(reply)
		if err != nil {
			return err
		}
	}

	replies := make([]*Comment, 0)
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tidNum).OrderBy("-created").All(&replies)
	if err != nil {
		return err
	}

	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		topic.ReplyCount = int64(len(replies))
		if topic.ReplyCount > 0 {
			topic.ReplyTime = replies[0].Created
		}

		_, err = o.Update(topic)
	}

	return err
}
