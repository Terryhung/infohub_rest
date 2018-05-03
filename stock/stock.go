package stock

import (
	"github.com/Terryhung/infohub_rest/utils"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Stock struct {
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	Prediction float64 `json:"prediction"`
	Diff       float64 `json:"diff"`
	Price      float64 `json:"price"`
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

func (s *Stock) GetStockPrice(session *mgo.Session, results interface{}) {
	// Regular Expression

	// Condition
	cond := bson.M{"ts": bson.M{"$gte": utils.NowTSNorm()}}
	col := session.DB("infohub_sandbox").C("infohub_stock_price")
	count, _ := col.Find(cond).Count()

	// results
	if count > 0 {
		col.Find(cond).All(results)
	}
}
