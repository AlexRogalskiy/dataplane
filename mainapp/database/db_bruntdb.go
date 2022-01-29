package database

import "github.com/tidwall/buntdb"

var GoDBWorkerGroup *buntdb.DB
var GoDBWorker *buntdb.DB

func GoDBConnect() {
	GoDBWorkerGroup, _ = buntdb.Open(":memory:")

	GoDBWorkerGroup.CreateIndex("environment", "*", buntdb.IndexJSON("Env"))

	GoDBWorker, _ = buntdb.Open(":memory:")

	GoDBWorker.CreateIndex("environment", "*", buntdb.IndexJSON("Env"))
}
