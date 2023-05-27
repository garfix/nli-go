package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/server"
	"testing"
	"time"
)

func TestAll(t *testing.T) {

	srv := server.NewServer("3334")
	srv.RunInBackground()
	defer srv.Close()

	time.Sleep(500 * time.Millisecond)

	client := server.CreateTestClient()
	defer client.Close()

	client.RunFile("blocks", common.Dir()+"/../../resources/blocks/test/test1.yml")
	client.RunFile("dbpedia", common.Dir()+"/../../resources/dbpedia/test/test1.yml")
	client.RunFile("dualworld", common.Dir()+"/../../resources/dualworld/test/test1.yml")
	client.RunFile("expressions", common.Dir()+"/../../resources/expressions/test/test1.yml")
	client.RunFile("helloworld", common.Dir()+"/../../resources/helloworld/test/test1.yml")
	client.RunFile("shell", common.Dir()+"/../../resources/shell/test/test1.yml")
	client.RunFile("relationships", common.Dir()+"/../../resources/relationships/test/test1.yml")
}
