package server

import (
	"net/http"
	"fmt"
	// "io"
	"time"
	"strconv"
	"encoding/json"
	// "log"

	"github.com/gorilla/mux"
)

/*
*SERIALIZACION DE DATOS
*/
type ActivityDocument struct { //Esto es para serializar los datos
	Activity 	Activity	`json:"Activity"`
}
type ResponseIDDocument struct { //Esto tambien es para serializar los datos. Este funciona solo para devolver un dato.
	ID 		uint64 	`json:"id"` //Este unicamente funciona para devolver un mensaje como respuesta del documento que se ha agregado
}
type ResponseError struct {
	Error 		string	`json:"error"`
}

type JSONResponse interface {
	LaunchResponse(writer http.ResponseWriter, statusCode int)
}

/*
	*SERVER RELATED
*/

type ActivitiesHandler struct { //"Estructura" de datos que guarda las actividades realizadas
	Activities *Activities
}

func NewHttpServer(addr string) *http.Server { //Creacion de funcion la cual devuelve el puntero para instanciar un servidor
	activities := &ActivitiesHandler{ //Instancia adicional de un struct ActivitiesHandler para agregar funciones que agregan datos a la estructura Activities
		Activities: &Activities{},
	}

	r := mux.NewRouter() //Este mux contendra los handler que el servidor puede resolver
	r.HandleFunc("/", renderHome).Methods("GET") //POST que funcione en root /
	r.HandleFunc("/addActivity", activities.addActivity).Methods("POST") //POST que funcione en root /
	// r.HandleFunc("/getActivity/{id}", activities.getActivity).Methods("GET") //POST que funcione en root /
	r.HandleFunc("/getActivity", activities.getActivity).Methods("GET") //POST que funcione en root /
	r.HandleFunc("/editActivity", activities.editActivity).Methods("PATCH") //POST que funcione en root /
	r.HandleFunc("/deleteActivity", activities.deleteActivity).Methods("DELETE") //POST que funcione en root /

	return &http.Server{
		Addr: addr,
		Handler: r,
	}
}

/*
	*HANDLERS
*/

func renderHome(w http.ResponseWriter, r *http.Request) {
	actualTime := time.Now()
	fmt.Fprintf(w, "Welcome to the activities API. Current time is: %s\n", actualTime.Format("15:04:05"))
	// w.Write([]byte("This message comes from root\n"))
}

/*
	* addActivity/{id}
*/
func (s *ActivitiesHandler) addActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo no permitido en esta ruta", http.StatusMethodNotAllowed)
		return
	}

	var request ActivityDocument //1. tener un documento del cual se van a serializar los datos
	err := json.NewDecoder(r.Body).Decode(&request) //2. Debido al post, se necesita decodificar la informacion

	if err != nil { //3. Verificar que la conversion de los datos se haya hecho correctamente
		http.Error(w, "Bad request. Verify your data.", http.StatusInternalServerError)
		//4. Recordar que la manera de manejar errores en las solicitudes es con http.Error
		return
	}
	//Toda esta operacion de aqui la dejo por que es descriptiva (Siguientes 5 lineas)
	inserted := s.Activities.InsertActivity(request.Activity) //5. Simulacion del manejo de persistencia. Esto basado en activities
	response := ResponseIDDocument{ID: inserted} //6. Generar una respuesta. Aqui si debemos de usar el formato de los struct json
	w.Header().Set("Content-Type", "application/json") //7. Asignar header ya que estamos devolviendo formato json a la respuesta
	w.WriteHeader(http.StatusCreated)
	encodedResponseOK := json.NewEncoder(w).Encode(response) //8. Codificar los datos y enviarlos en el Writer.

	if encodedResponseOK != nil { //9. Verificar que se haya hecho la conversion al enviar los datos de manera correcta
		http.Error(w, "En caso de que no se haya convertido correctamente", http.StatusBadRequest)
		return
	}
	// fmt.Printf("Inserted activity was: %s with the id of %d\n", request.Activity.Title, inserted)
}

func (s *ActivitiesHandler) getActivity(w http.ResponseWriter, r *http.Request) {	

	var response JSONResponse

	if r.Method != http.MethodGet { //Esto en caso de que el metodo usado no sea el correcto
		response = ResponseError{Error: "Metodo no permitido en esta ruta"}
		response.LaunchResponse(w, http.StatusMethodNotAllowed)
		return
	}
	
	idStr := r.URL.Query().Get("id")
	if idStr == "" { //Obtencion del parametro id del query
		response = ResponseError{Error: "Cannot retrieve id from query"}
		response.LaunchResponse(w, http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64) //Parseo del query param
	if err != nil {
		response = ResponseError{Error: "Given id could not be parsed, try integer value"}
		response.LaunchResponse(w, http.StatusBadRequest)
		return
	}
	activityItem, err := s.Activities.GetActivity(id)
	if err != nil { //Ejecucion de la funcion para agregar la actividad
		response = ResponseError{Error: "Activity could not be found"}
		response.LaunchResponse(w, http.StatusBadRequest)
		return
	}

	response = ActivityDocument{Activity: activityItem}
	response.LaunchResponse(w, http.StatusOK)
}

func (s *ActivitiesHandler) editActivity(w http.ResponseWriter, r *http.Request){

	var response JSONResponse

	if r.Method != http.MethodPatch {
		http.Error(w, "Unable to use this route because verb", http.StatusMethodNotAllowed)
		return
	}
	//Marshaling data
	var newData ActivityDocument
	unmarshalErr := json.NewDecoder(r.Body).Decode(&newData)
	if unmarshalErr != nil {
		response = ResponseError{Error:"Cannot convert data"}
		response.LaunchResponse(w, http.StatusInternalServerError)
		return
	}

	//Search for the given activity id
	for _, activity := range s.Activities.activities {
		if activity.ID == newData.Activity.ID {
			response = ResponseError{Error:"Given id could not be found"}
			response.LaunchResponse(w, http.StatusBadRequest)
			return
		}
	}

	//Manage to give a json response
	actId, ok := s.Activities.EditActivity(newData.Activity)
	if ok != nil {
		response = ResponseError{Error: "Edit activity could not be applied. Maybe data does not exist"}
		response.LaunchResponse(w, http.StatusOK)
		return
	}
	response = ResponseIDDocument{
		ID: actId,
	}
	response.LaunchResponse(w, http.StatusOK)
}

func (s *ActivitiesHandler) deleteActivity(w http.ResponseWriter, r *http.Request) {
	var response JSONResponse

	idStr := r.URL.Query().Get("id")
	if idStr == "" { //Obtencion del parametro id del query
		response = ResponseError{Error: "Cannot retrieve id from query"}
		response.LaunchResponse(w, http.StatusBadRequest)
		return
	}
	id, parseErr := strconv.ParseUint(idStr, 10, 64) //Parseo del query param
	if parseErr != nil {
		response = ResponseError{Error: "Given id could not be parsed, try integer value"}
		response.LaunchResponse(w, http.StatusBadRequest)
		return
	}
	err := s.Activities.DeleteActivity(id)
	if err != nil {
		response = ResponseError{Error: err.Error()}
		response.LaunchResponse(w, http.StatusOK)
		return
	}

	response = ResponseIDDocument{ID:id}
	response.LaunchResponse(w, http.StatusOK)
}


/*
	* Implementations for LaunchResponse
*/
func (r ResponseError) LaunchResponse(writer http.ResponseWriter, statusCode int) {
	badResponse := ResponseError{Error:r.Error}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	json.NewEncoder(writer).Encode(badResponse)
}

func (r ResponseIDDocument) LaunchResponse(writer http.ResponseWriter, statusCode int) {
	response := ResponseIDDocument{ID: r.ID}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	json.NewEncoder(writer).Encode(response)
}

func (r ActivityDocument) LaunchResponse(writer http.ResponseWriter, statusCode int) {
	response := ActivityDocument{Activity: r.Activity}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	json.NewEncoder(writer).Encode(response)
}



