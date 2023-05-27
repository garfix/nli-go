package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/server"
	"testing"
	"time"
)

// Mimics Terry Winograd's SHRDLU dialog, but in the NLI-GO way

func TestBlocksWorld(t *testing.T) {

	srv := server.NewServer("3334")
	srv.RunInBackground()
	defer srv.Close()

	time.Sleep(500 * time.Millisecond)

	client := server.CreateTestClient("blocks")
	defer client.Close()

	client.RunFile(common.Dir() + "/../../resources/blocks/test/test1.yml")
}
