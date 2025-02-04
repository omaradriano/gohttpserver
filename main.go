package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const portNum string := ":3333"

func getRoot(w http.ResponseWriter, r *http.Request) {
	//Context viene con cada una de las request
	// context := r.Context()
	var str string
	firstArgExist := r.URL.Query().Has("name") //Verify if the param exist
	firstArg := r.URL.Query().Get("name") //Get it in case it exist
	w.WriteHeader(200)
	if firstArgExist {
		str = fmt.Sprintf("This is my website, welcome %s\n", firstArg)
	}else{
		str = "This is my website\n"
	}
	io.WriteString(w,str)
}
func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

func deployHome(w http.ResponseWriter, r *http.Request){
	fmt.Println("Se ha accedido a la ruta /home")
	io.WriteString(w, "Usted esta en /home\n")
}

func main(){

	multiplexer := http.NewServeMux()

	multiplexer.HandleFunc("/", getRoot)
	//Tener en cuenta que root maneja cualquier otra ruta que no este registrada
	multiplexer.HandleFunc("/hello", getHello)
	multiplexer.HandleFunc("/home", deployHome)

	err := http.ListenAndServe(portNum, multiplexer)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
	// jsondata := map[string]interface{}{
	// 	"name":	"Omar Acosta",
	// 	"age":	24,
	// 	"bornDate": time.Date(2000, 1, 2, 9, 10, 0, 0, time.UTC), //printed format is RFC 3339
	// 	"hobbies": map[string]interface{}{
	// 		"favoriteGames": []string{"Skyrim", "Dark Souls"},
	// 	},
	// 	"working": nil, //This will be interpreted as null as JSON format
	// }

	//creating json from a STRUCT

	type MyValues struct {
		MyName string `json:"variable extra xd"` //This tag converts name 
		MyAge int
		CurrentlyWorking bool `json:"currentlyWorking,omitempty"`

	}

	// structData := &MyValues{
	// 	MyName: "Omar Acosta", 
	// 	MyAge: 24,
	// 	CurrentlyWorking: true,
	// }

	// convertedData, ok := json.Marshal(jsondata)
	// convertedStructData, structOk := json.Marshal(structData)

	// if ok != nil {
	// 	fmt.Println("Se ha generado un nuevo error al serializar el json")
	// }else{
	// 	fmt.Printf("Generated data is %s\n",convertedData)
	// }
		
	// if structOk != nil {
	// 	fmt.Println("Se ha generado un nuevo error al serializar el json")	
	// }else{
	// 	fmt.Printf("Generated data from struct is %s\n",string(convertedStructData))
	// }

}