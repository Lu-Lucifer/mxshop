package db

import (
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
)

type EsOptions struct {
	Host string
	Port int
}

func NewEsClient(opts *EsOptions) (*elastic.Client, error) {
	//初始化连接es
	host := fmt.Sprintf("http://%s:%d", opts.Host, opts.Port)
	logger := log.New(os.Stdout, "elastic-info", log.LstdFlags)
	var err error
	esClient, err := elastic.NewClient(
		elastic.SetURL(host),
		elastic.SetSniff(false),
		elastic.SetErrorLog(log.New(os.Stderr, "elastic-error", log.LstdFlags)),
		elastic.SetInfoLog(logger),
	)
	if err != nil {
		return nil, err
	}
	return esClient, err
}
