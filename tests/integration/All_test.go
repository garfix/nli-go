package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/server"
	"testing"
	"time"
)

func TestAll(t *testing.T) {

	appDir := common.Dir() + "/../../resources"
	workDir := common.Dir() + "/../../var"

	srv := server.NewServer("3334", appDir, workDir)
	srv.RunInBackground()
	defer srv.Close()

	time.Sleep(500 * time.Millisecond)

	client := server.CreateTestClient(t)
	defer client.Close()

	const all = 1

	if all == 0 {
		client.RunTests("blocks", common.Dir()+"/../../resources/blocks/test/test1.yml")
	} else {
		client.RunTests("blocks", common.Dir()+"/../../resources/blocks/test/test1.yml")
		client.RunTests("blocks", common.Dir()+"/../../resources/blocks/test/test2.yml")
		client.RunTests("dbpedia", common.Dir()+"/../../resources/dbpedia/test/test1.yml")
		client.RunTests("dualworld", common.Dir()+"/../../resources/dualworld/test/test1.yml")
		client.RunTests("expressions", common.Dir()+"/../../resources/expressions/test/test1.yml")
		client.RunTests("helloworld", common.Dir()+"/../../resources/helloworld/test/test1.yml")
		client.RunTests("shell", common.Dir()+"/../../resources/shell/test/test1.yml")
		client.RunTests("relationships", common.Dir()+"/../../resources/relationships/test/test1.yml")
	}
}
