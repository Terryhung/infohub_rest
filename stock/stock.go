package stock

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Stock struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func (s *Stock) GetStockList(session *mgo.Session, results interface{}) {
	// Regular Expression

	// Condition
	cond := bson.M{}
	col := session.DB("infohub_sandbox").C("infohub_stock")
	count, _ := col.Find(cond).Count()

	// results
	if count > 0 {
		col.Find(cond).All(results)
	}
}
