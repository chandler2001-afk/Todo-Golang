package main_api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type Task struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}
type Delete_Task struct{
	T_id int `json:"t_id"`
}
var (
	tasks   []Task
	taskID  int
	taskMux sync.Mutex
)

func CreateTask(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Error reading request body!", http.StatusInternalServerError)
		return
	}
	var task Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		http.Error(res, "Error unmarshalling task!", http.StatusInternalServerError)
		return
	}
	taskMux.Lock()
	taskID++
	task.ID = taskID
	tasks = append(tasks, task)
	taskMux.Unlock()

	fmt.Printf("Task added: %+v\n", task)
	fmt.Printf("Total tasks: %d\n", len(tasks))

	res.Header().Set("Content-Type", "application/json")
	taskJSON, err := json.Marshal(task)
	if err != nil {
		http.Error(res, "Error Marshalling the task", http.StatusInternalServerError)
		return
	}
	res.Write(taskJSON)
}

func ShowTasks(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	taskMux.Lock()
	defer taskMux.Unlock()

	fmt.Printf("Number of tasks: %d\n", len(tasks))

	res.Header().Set("Content-Type", "application/json")

	if len(tasks) == 0 {
		fmt.Println("No tasks found, returning empty array")
		res.Write([]byte("[]"))
		return
	}

	taskJSON, err := json.Marshal(tasks)
	if err != nil {
		fmt.Printf("Error marshalling tasks: %v\n", err)
		http.Error(res, "Error loading the tasks!!", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Returning JSON: %s\n", string(taskJSON))

	_, err = res.Write(taskJSON)
	if err != nil {
		fmt.Printf("Error writing response: %v\n", err)
	}
}


//deleting a task

func DeleteTask(res http.ResponseWriter,req *http.Request){
	if req.Method!=http.MethodDelete{
		http.Error(res,"Method Not Allowed",http.StatusMethodNotAllowed)
		return
	}
	body,err:=ioutil.ReadAll(req.Body)
	if err!=nil{
		http.Error(res,"Error reading the request Body",http.StatusBadRequest)
		return
	}
	var dt Delete_Task
    err = json.Unmarshal(body, &dt)
    if err != nil {
        http.Error(res, "Error unmarshalling the request body", http.StatusBadRequest)
        return
    }

    fmt.Printf("Received delete request for task ID: %v\n", dt.T_id)

    taskMux.Lock()
    defer taskMux.Unlock()
	fmt.Printf("Current tasks: %+v\n", tasks)
	for i ,task:=range tasks{
		fmt.Printf("Comparing task ID %v with requested ID %v\n", task.ID, dt.T_id)
		if dt.T_id==task.ID{
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Printf("Task with ID %d deleted\n",dt.T_id)
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(fmt.Sprintf("Task with ID %d deleted", dt.T_id)))
            return
		}
		}
		http.Error(res,"Task Not Found",http.StatusNotFound)
}
