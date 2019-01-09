package mongo_lib

import (
	"strings"

	"gopkg.in/mgo.v2/bson"
)

func GenRedisKey(keys []string) string {
	return strings.Join(keys, "-")
}

func GenCondition(cat string, lang string, cty string, cty_ary bson.M, source_date int32) bson.M {
	cond := bson.M{"source_date_int": bson.M{"$gte": source_date}, "category": cat, "language": lang}

	// Check using country or country_array
	if cty != "" {
		cond["country"] = cty
	} else if len(cty_ary) > 0 {
		cond["country_array"] = cty_ary
	}

	return cond
}
