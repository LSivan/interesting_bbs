package models

import (
	"github.com/astaxie/beego/orm"
)

type Section struct {
	Id   int    `orm:"pk;auto"`
	Name string `orm:"unique"`
}

func FindAllSection() []*Section {
	o := orm.NewOrm()
	var section Section
	var sections []*Section
	o.QueryTable(section).All(&sections)
	return sections
}

func FindSectionById(id int) (bool, Section) {
	o := orm.NewOrm()
	var section Section
	err := o.QueryTable(section).RelatedSel().Filter("Id", id).One(&section)
	return err != orm.ErrNoRows, section
}
