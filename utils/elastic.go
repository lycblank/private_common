package utils

import elastic "gopkg.in/olivere/elastic.v5"

var elasticClient *elastic.Client

func InitElastic(addr string) {
	var err error
	elasticClient, err = elastic.NewClient(elastic.SetTraceLog(new(ElasticDebugLog)), elastic.SetInfoLog(new(ElasticInfoLog)),
		elastic.SetErrorLog(new(ElasticErrorLog)), elastic.SetURL(addr))
	if err != nil {
		panic(err)
	}
}

func GetElasticClient() *elastic.Client {
	return elasticClient
}
