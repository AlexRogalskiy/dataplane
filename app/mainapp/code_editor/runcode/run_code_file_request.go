package runcode

import (
	"errors"

	"github.com/dataplane-app/dataplane/app/mainapp/code_editor/filesystem"
	"github.com/dataplane-app/dataplane/app/mainapp/database"
	"github.com/dataplane-app/dataplane/app/mainapp/database/models"

	"github.com/google/uuid"
)

type Command struct {
	Command string `json:command`
}

/*
Task status: Queue, Allocated, Started, Failed, Success
*/
func RunCodeFile(workerGroup string, fileID string, envID string, pipelineID string, nodeID string, nodeTypeDesc string, runid string) (models.CodeRun, error) {

	// Important not to update status to avoid timing issue where it can overwrite a success a status
	if runid == "" {
		runid = uuid.NewString()
	}
	// ------ Obtain folder structure from file id

	filesdata := models.CodeFiles{}
	database.DBConn.Where("file_id = ? and environment_id =? and level = ?", fileID, envID, "node_file").Find(&filesdata)

	parentfolderdata := ""
	var err error
	if filesdata.FolderID != "" {
		parentfolderdata, err = filesystem.FolderConstructByID(database.DBConn, filesdata.FolderID, envID, "pipelines")
		if err != nil {
			return models.CodeRun{}, errors.New("File record not found")
		}
	} else {
		return models.CodeRun{}, errors.New("File record not found")
	}

	// if dpconfig.Debug == "true" {
	// 	if _, err := os.Stat(dpconfig.CodeDirectory + folderMap); os.IsExist(err) {
	// 		log.Println("Dir exists:", dpconfig.CodeDirectory+folderMap)

	// 	}
	// }

	// ------ Construct run command
	var commands []string
	var runSend models.CodeRun
	switch nodeTypeDesc {
	case "python":
		commands = append(commands, "python3 -u ${{nodedirectory}}"+filesdata.FileName)
		runSend, err = RunCodeServerWorker(envID, nodeID, workerGroup, runid, commands, filesdata, parentfolderdata, filesdata.FolderID)
		if err != nil {
			return runSend, err
		}
	case "rpa-python":
		commands = append(commands, "python3 -u ${{nodedirectory}}"+filesdata.FileName)
		runSend, err = RunCodeRPAWorker(envID, nodeID, workerGroup, runid, commands, filesdata, parentfolderdata, filesdata.FolderID)
		if err != nil {
			return runSend, err
		}
	default:
		return models.CodeRun{}, errors.New("Code run type not found.")
	}

	//
	return runSend, nil

}
