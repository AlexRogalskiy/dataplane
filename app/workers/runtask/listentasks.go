package runtask

import (
	"context"
	"dataplane/mainapp/database/models"
	"dataplane/workers/config"
	"dataplane/workers/messageq"
	"log"
	"os"
	"syscall"
)

type TaskResponse struct {
	R string
	M string
}

func ListenTasks() {

	// Responding to a task request
	messageq.NATSencoded.Subscribe("task."+os.Getenv("worker_group")+"."+config.WorkerID, func(subj, reply string, msg models.WorkerTaskSend) {
		// log.Println(msg)

		response := "ok"
		message := "ok"
		if os.Getenv("worker_env") != msg.EnvironmentID {
			response = "failed"
			message = "Incorrect environment"
		}

		x := TaskResponse{R: response, M: message}
		messageq.NATSencoded.Publish(reply, x)

		if x.R == "ok" {
			TaskID := msg.TaskID
			ctx, cancel := context.WithCancel(context.Background())
			var task Task

			task.ID = TaskID
			task.Context = ctx
			task.Cancel = cancel

			Tasks[task.ID] = task
			// command := `for((i=1;i<=10000; i+=1)); do echo "Welcome $i times"; sleep 1; done`
			// command := `find . | sed -e "s/[^ ][^\/]*\// |/g" -e "s/|\([^ ]\)/| \1/"`
			go worker(ctx, msg.RunID, TaskID, msg.Commands)
		}
	})
	if os.Getenv("debug") == "true" {
		log.Println("Listening for tasks on subject:", "task."+os.Getenv("worker_group")+"."+config.WorkerID)
	}

	messageq.NATSencoded.Subscribe("taskcancel."+os.Getenv("worker_group")+"."+config.WorkerID, func(subj, reply string, msg models.WorkerTaskSend) {
		// Respond to cancelling a task
		id := msg.TaskID

		if Tasks[id].PID != 0 {
			_ = syscall.Kill(-Tasks[id].PID, syscall.SIGKILL)
		}
		Tasks[id].Cancel()
		TasksStatus[id] = "cancel"

		response := "ok"
		message := "ok"
		x := TaskResponse{R: response, M: message}
		messageq.NATSencoded.Publish(reply, x)

		// TaskUpdate := modelmain.WorkerTasks{
		// 	TaskID: id,
		// 	EndDT:  time.Now().UTC(),
		// 	Status: "Fail",
		// 	Reason: "Cancelled",
		// }
		// var response TaskResponse
		// _, errnats := messageq.MsgReply("taskupdate", TaskUpdate, &response)

		// if errnats != nil {
		// 	logging.PrintSecretsRedact("Cancel task error nats:", errnats)
		// }

		// delete(Tasks, id)
		// delete(TasksStatus, id)
	})

}
