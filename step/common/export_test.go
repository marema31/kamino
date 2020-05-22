package common

import "github.com/Sirupsen/logrus"

func NewSkipQuery(query string, inverted bool, compareValue int) SkipQuery {
	return SkipQuery{query: query, inverted: inverted, compareValue: compareValue}
}

func (skq *SkipQuery) CompareSkipQuery(query string, inverted bool, compareValue int) bool {
	if skq.query != query {
		logrus.Errorf("Query '%s' != '%s'", skq.query, query)
		return false
	}
	if skq.inverted != inverted {
		logrus.Errorf("inverted '%v' != '%v'", skq.inverted, inverted)
		return false
	}
	if skq.compareValue != compareValue {
		logrus.Errorf("compareValue '%d' != '%d'", skq.compareValue, compareValue)
		return false
	}
	return true
}
