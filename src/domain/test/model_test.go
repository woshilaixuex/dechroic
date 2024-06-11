package test

import (
	"context"
	"encoding/json"
	"flag"
	"testing"

	"github.com/delyr1c/dechoric/src/domain/domainModel/award"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var configFile = flag.String("f", "etc/test.yaml", "the config file")

type Config struct {
	DB struct {
		MySqlDataSource string
	}
}

func TestDomainModel(t *testing.T) {
	flag.Parse()
	var c Config
	conf.MustLoad(*configFile, &c)

	sqlConn := sqlx.NewMysql(c.DB.MySqlDataSource)
	AwardModel := award.NewAwardModel(sqlConn)
	awards, err := AwardModel.FindAll(context.Background())
	if err != nil {
		t.Fatalf("failed to find all awards: %v", err)
	}

	awardsJSON, err := json.Marshal(awards)
	if err != nil {
		t.Fatalf("failed to marshal awards to JSON: %v", err)
	}

	t.Log(string(awardsJSON))
}
