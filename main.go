package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type endpoint struct {
	UserID   string `json:"UserID"`
	Type     string `json:"Type"`
	Endpoint string `json:"Endpoint"`
}

type requestBody struct {
	UserID string `json:"user_id"`
	Data   string `json:"data"`
}

type responseBody struct {
	Message string `json:"message"`
}

/*
 * I'm considering we'll have a database with user IDs, and one destination endpoint
 * for every kind of message/notification - i.e: message, warning, new-subscription, etc.
 */
var db = []endpoint{
	{UserID: "1", Type: "message", Endpoint: "https://postman-echo.com/post?type=message&user_id=1"},
	{UserID: "1", Type: "warning", Endpoint: "https://postman-echo.com/post?type=warning&user_id=1"},
	{UserID: "2", Type: "message", Endpoint: "https://postman-echo.com/post?type=warning&user_id=2"},
}

/*
 * Simple "I'm alive" message on main endpoint '/'
 */
func mainHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("I'm running!")
	fmt.Fprintln(w, "I'm running!\nTo test me, send a POST request to /webhook/message with body { \"user_id\": \"1\", \"data\": \"<any message here>\"}")
}

/*
 * Handle incoming request to send a webhook of type 'message' to 'user_id' containing 'data'.
 * It'll search the database for the corresponding endpoint, and start a goroutine
 * to POST the data to the endpoint.
 */
func messageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling message webhook!")

	var request requestBody

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendResponseJSON(w,
			http.StatusInternalServerError,
			responseBody{fmt.Sprintf("Ops...")},
		)
		return
	}

	json.Unmarshal(reqBody, &request)

	for _, ep := range db {
		if ep.UserID == request.UserID && ep.Type == "message" {
			go sendHTTPPost(ep.Endpoint, request.Data)

			res := responseBody{fmt.Sprintf("Message sent to %s", ep.Endpoint)}
			sendResponseJSON(w, http.StatusOK, res)
			return
		}
	}

	res := responseBody{fmt.Sprintf("Endpoint for user %s not found!", request.UserID)}
	sendResponseJSON(w, http.StatusNotFound, res)
}

/*
 * Basic JSON response to incoming requests
 */
func sendResponseJSON(w http.ResponseWriter, status int, r responseBody) {
	log.Printf("Send response status=[%d][%v]", status, r)
	res, err := json.Marshal(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(res)
}

/*
 * Function to send an HTTP POST request to the corresponding endpoint.
 * Will be called from a goroutine to execute in parallel with the main thread.
 */
func sendHTTPPost(endpoint string, data string) {
	uuid := uuid.New()
	start := time.Now()

	requestBody, err := json.Marshal(map[string]string{
		"data": data,
	})
	if err != nil {
		log.Printf("Error creating request body: %v", err)
		return
	}

	log.Printf("%s | HTTP request to [%s] containing [%v] ", uuid, endpoint, data)

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("%s | Error sending HTTP request to %s: %v", uuid, endpoint, err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("%s | Error parsing HTTP response: %v", uuid, err)
		return
	}

	log.Printf("%s | HTTP response: [%s], took %s", uuid, string(body), time.Since(start))
}

/*
 * Main function
 */
func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", mainHandler)
	router.HandleFunc("/webhooks/message", messageHandler).Methods("POST")

	log.Println("Starting application on port 8000!")
	log.Fatal(http.ListenAndServe(":8000", router))
}
