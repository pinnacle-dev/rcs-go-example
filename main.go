package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type PinnacleClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewPinnacleClient(apiKey string) *PinnacleClient {
	return &PinnacleClient{
		apiKey:  apiKey,
		baseURL: "https://www.trypinnacle.dev/api",
		client:  &http.Client{},
	}
}

var client *PinnacleClient

func send_rcs_media_card() {
	url := client.baseURL + "/send/rcs"

	payloadData := map[string]interface{}{
		"from": "test",
		"to":   "+16287261512",
		"cards": []map[string]interface{}{
			{
				"title":    "I do!",
				"subtitle": "This is one of my favorite songs",
				"mediaUrl": "https://www.dropbox.com/scl/fi/99zd9kp92yxdb74k642x4/Hello-World-1.mp4?rlkey=otev4efja2rpefqplw1vn649p&st=t2kys8tu&raw=1", // Raw .mp4 is supported--you can get a raw link usually by adding raw=1 to things like Dropbox links. Sound is also supported.
				"buttons": []map[string]string{
					{
						"title":   "How can I build this?",
						"payload": "https://docs.trypinnacle.app/api-reference/api-reference/send-a-message/rcs",
						"type":    "openUrl",
					},
				},
			},
		},
		"quickReplies": []map[string]string{
			{
				"title":   "Reset",
				"payload": "RESET",
				"type":    "trigger",
			},
		},
	}

	payload, err := json.Marshal(payloadData)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	req, _ := http.NewRequest("POST", url, bytes.NewReader(payload))
	req.Header.Add("PINNACLE-API-KEY", client.apiKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println(res)
	fmt.Println(string(body))
}

func send_rcs_with_quick_replies() {
	url := client.baseURL + "/send/rcs"

	payloadData := map[string]interface{}{
		"from": "test",
		"to":   "+16287261512",
		"text": "Hello, World!",
		"quickReplies": []map[string]string{
			{
				"title":   "Hi, World!",
				"payload": "HI_WORLD",
				"type":    "trigger",
			},
			{
				"title":   "Got any jams?",
				"payload": "ANY_JAMS?",
				"type":    "trigger",
			},
		},
	}

	payload, err := json.Marshal(payloadData)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	req, _ := http.NewRequest("POST", url, bytes.NewReader(payload))
	req.Header.Add("PINNACLE-API-KEY", client.apiKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println(res)
	fmt.Println(string(body))
}

func send_basic_rcs() {
	url := client.baseURL + "/send/rcs"

	payloadData := map[string]interface{}{
		"from": "test",
		"to":   "+16287261512",
		"text": "...",
	}

	payload, err := json.Marshal(payloadData)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	req, _ := http.NewRequest("POST", url, bytes.NewReader(payload))
	req.Header.Add("PINNACLE-API-KEY", client.apiKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println(res)
	fmt.Println(string(body))
}

func say_hi_back() {
	// Send a simple "Hi" response
	url := client.baseURL + "/send/rcs"
	payloadData := map[string]interface{}{
		"from": "test",
		"to":   "+16287261512",
		"text": "Hi! 👋",
		"quickReplies": []map[string]string{
			{
				"title":   "Got any jams?",
				"payload": "ANY_JAMS?",
				"type":    "trigger",
			},
		},
	}

	payload, err := json.Marshal(payloadData)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	req, _ := http.NewRequest("POST", url, bytes.NewReader(payload))
	req.Header.Add("PINNACLE-API-KEY", client.apiKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer res.Body.Close()
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Print the received webhook data
	fmt.Printf("Received webhook: %s\n", string(body))

	// Parse JSON body if needed
	var webhookData map[string]interface{}
	if err := json.Unmarshal(body, &webhookData); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	// Handle action messages
	if messageType, ok := webhookData["messageType"].(string); ok && messageType == "action" {
		// Extract action details
		actionTitle := webhookData["actionTitle"].(string)
		payload, hasPayload := webhookData["payload"].(string)
		actionMetadata, hasMetadata := webhookData["actionMetadata"].(string)

		fmt.Printf("Received action: %s\n", actionTitle)
		if hasPayload {
			fmt.Printf("Payload: %s\n", payload)
		}
		if hasMetadata {
			fmt.Printf("Metadata: %s\n", actionMetadata)
		}

		switch payload {
		case "ANY_JAMS?":
			send_rcs_media_card()
		case "HI_WORLD":
			say_hi_back()
		case "RESET":
			send_basic_rcs()
			send_rcs_with_quick_replies()
		}
	}

	// Send a success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received successfully"))
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	// Initialize the client
	apiKey := os.Getenv("PINNACLE_API_KEY")
	if apiKey == "" {
		log.Fatal("Error: PINNACLE_API_KEY environment variable not set")
	}
	client = NewPinnacleClient(apiKey)

	// Set up webhook endpoint
	http.HandleFunc("/", webhookHandler)

	// Send initial messages
	send_basic_rcs()
	send_rcs_with_quick_replies()

	fmt.Println("Starting webhook server on :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
