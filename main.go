package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Tipos tareas
type task struct {
	ID      int    `json:"ID"`
	Name    string `json:"Name"`
	Content string `json:"Content"`
}

type allTasks []task

// Primera tarea
var tasks = allTasks{
	{
		ID:      1,
		Name:    "Task One",
		Content: "Some Content",
	},
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)                     //recibimos las variables vamos a seleccionar "id"
	taskID, err := strconv.Atoi(vars["id"]) // el "id" recibido como string es convertido a entero
	var updatedTask task
	if err != nil {
		fmt.Fprintf(w, "Invalid ID") //para el caso que el id no se pueda covertir a entero "que manden fruta"
		return
	}
	reqBody, err := ioutil.ReadAll(r.Body) //leer la data brindada por el usuario
	if err != nil {
		fmt.Fprintf(w, "Please Enter Valid Data") //para el caso que el id no se pueda covertir a entero "que manden fruta"
	}
	json.Unmarshal(reqBody, &updatedTask)

	for i, t := range tasks { //range sobre lista de tareas
		if t.ID == taskID { //comparacion entre taskID obtenido y el ID dentro de task
			tasks = append(tasks[:i], tasks[i+1:]...) //borramos la data almacenada
			updatedTask.ID = taskID                   //asignamos el ID del [] borrado
			tasks = append(tasks, updatedTask)        //agregamos la nueva tarea
			fmt.Fprintf(w, "The task with ID: %v has been successfully updated.", taskID)
		}
	}

}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task //con esto recibimos la información
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insert a valid task")
	} //aviso en caso de error
	json.Unmarshal(reqBody, &newTask) //asigno esa informacion a variable newTask

	newTask.ID = len(tasks) + 1    //le genero un nuevo ID que es el largo del slice+1
	tasks = append(tasks, newTask) //guardo la tarea en la lista de tasks con su nuevo ID

	w.Header().Set("Content-Type", "application/json") //devolvemos un header que informa el tipo de dato regresado
	w.WriteHeader(http.StatusCreated)                  //aviso de la creacion exitosa de la tarea
	json.NewEncoder(w).Encode(newTask)                 //respondo al cliente con la tarea creada + ID
}

func getTask(w http.ResponseWriter, r *http.Request) { //atento aca que es getTask en singular
	vars := mux.Vars(r)                     //recibimos las variables vamos a seleccionar "id"
	taskID, err := strconv.Atoi(vars["id"]) // el "id" recibido como string es convertido a entero
	if err != nil {
		fmt.Fprintf(w, "Invalid ID") //para el caso que el id no se pueda covertir a entero "que manden fruta"
		return
	}

	for _, task := range tasks { //range sobre lista de tareas
		if task.ID == taskID { //comparacion entre taskID obtenido y el ID dentro de task
			w.Header().Set("Content-Type", "application/json") //devolvemos un header que informa el tipo de dato regresado
			json.NewEncoder(w).Encode(task)                    //desde el writer devolvemos la tarea
		}
	}

}

//Empleando Postman hacemos GET, POST, PUT, DELETE

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) //copio del anterior porque el metodo de busqueda es el mismo
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}
	for i, task := range tasks { //misma busqueda que en el caso anterior
		if task.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...) //en realidad no borra sino que arma un slice con los elementos previos y posteriores al indice dado
			fmt.Fprintf(w, "The task with ID: %v has been successfully removed.", taskID)
		}
	}
}

func getTasks(w http.ResponseWriter, r *http.Request) { //envia las tareas en formato json
	w.Header().Set("Content-Type", "application/json") //devolvemos un header que informa el tipo de dato regresado
	json.NewEncoder(w).Encode(tasks)
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Wecome to my GO REST API!")
}

func main() {
	router := mux.NewRouter().StrictSlash(true) //obligado a escribir correctamente la url

	router.HandleFunc("/", indexRoute)                             //index
	router.HandleFunc("/tasks", getTasks).Methods("GET")           //consultar tasks
	router.HandleFunc("/tasks", createTask).Methods("POST")        //post tasks
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")       //consulta por ID
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE") //borra por ID
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")    //actualiza por ID
	log.Fatal(http.ListenAndServe(":3000", router))                //asignacion de puerto
}
