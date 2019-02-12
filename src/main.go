package main

import (
	// Importamos las librerías necesarias, aunque con solo guardar se importan automáticamente :D
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) //clients Conectados
var broadcast = make(chan Message)           // Broadcast canal de transmision

// upgrader Este es solo un objeto con métodos para tomar una conexión HTTP normal y actualizarla a un WebSocket
var upgrader = websocket.Upgrader{}

//Message Definiremos un objeto para guardar nuestros mensajes, para interactuar con el servicio ***Gravatar*** que nos proporcionará un avatar único.
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

<<<<<<< HEAD
// esta funcion manejara nuestras conexiones  WebSocket entrantes
=======
func main() {
	// Create a simple file server
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	// Configure webSocket route
	http.HandleFunc("/ws", handleConnections)

	// Start listening for incomming chat messages
	go handleMessages()

	// Start the server on localhost port 8080 and log any errors
	log.Println("http server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

>>>>>>> c87cb8edc90bad1b8df6c51ba901d3b7d800e37c
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// El método Upgrade() permite cambiar nuesra solicitud GET inicial a una completa en WebSocket, si hay un error lo mostramos en consola pero no salimos.
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// para cerrar la conexion una vez termina la funcion
	defer ws.Close()

	// Registramos nuestro nuevo cliente al agregarlo al mapa global de "clients" que fue creado anteriormente.
	clients[ws] = true
	// Bucle infinito que espera continuamente que se escriba  un nuevo mensaje en el WebSocket, lo desserializa de JSON a un objeto Message y luego lo arroja al canal de difusión.
	for {
		var msg Message

		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}

		// Send the newly received message to the broadcast channel
		//Enviar el mensaje recién recibido al canal de difusión.
		broadcast <- msg
	}
}

//goroutine llamada "handleMessages"
func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast

		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
