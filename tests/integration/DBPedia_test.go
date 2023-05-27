package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/server"
	"testing"
	"time"
)

func TestDBPedia(t *testing.T) {

	srv := server.NewServer("3334")
	srv.RunInBackground()
	defer srv.Close()

	time.Sleep(500 * time.Millisecond)

	client := server.CreateTestClient("dbpedia")
	defer client.Close()

	client.RunFile(common.Dir() + "/../../resources/dbpedia/test/test1.yml")
}
