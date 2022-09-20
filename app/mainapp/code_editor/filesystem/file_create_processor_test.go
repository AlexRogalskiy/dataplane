package filesystem

import (
	"dataplane/mainapp/config"
	"dataplane/mainapp/database"
	"dataplane/mainapp/database/models"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
Run Super secret squirrel function test
go test -count=1 -timeout 30s -v -run ^TestFileCreateProcessor$ dataplane/mainapp/filesystem
*/
func TestFileCreateProcessor(t *testing.T) {

	dpconfig.LoadConfig()
	database.DBConnect()

	nodeTypeDesc := "python"

	start := time.Now()
	node := models.PipelineNodes{
		EnvironmentID: "test-environment-id",
		NodeID:        "test-node-id",
		PipelineID:    "test-pipeline-id",
	}

	// FileCreateProcessor(n.NodeTypeDesc, dpconfig.CodeDirectory+rfolder+"/", node)
	output, err := FileCreateProcessor(nodeTypeDesc, "", "testFolderID", node)
	if err != nil {
		t.Error(err)
	}

	log.Println("File location: ", dpconfig.CodeDirectory+output)
	stop := time.Now()
	// Do something with response
	log.Println("🐆 Runtime:", fmt.Sprintf("%f", float32(stop.Sub(start))/float32(time.Millisecond))+"ms")

	dat, err := os.ReadFile(dpconfig.CodeDirectory + output)
	if err != nil {
		t.Error(err)
	}
	fmt.Print(string(dat))

	expected := `print("Pipeline id: test-pipeline-id")
print("Node id: test-node-id")`

	assert.Equalf(t, expected, string(dat), "Processor entry file.")

}
