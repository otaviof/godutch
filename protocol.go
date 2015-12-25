package godutch

import (
	"encoding/json"
	"log"
)

//
// A Request is the basic query unit towards a GoDutch client, any incoming
// communication msut be wrapped on a "Request"
//
type Request struct {
	Command   string   `json:"command"`
	Arguments []string `json:"arguments"`
}

//
// To match Request, we also implement Response type, to hold a reponse and
// it's attributes.
//
type Response struct {
	Name    string           `json:"name"`
	Status  int              `json:"status"`
	Stdout  []string         `json:"stdout"`
	Metrics []map[string]int `json:"metrics,omitempty"`
	Error   string           `json:"error,omitempty"`
}

// Creates a slice of bytes that maches the JSON representation of informed
// arguments, the straight forward input to a socket.
func NewRequest(commandName string, arguments []string) ([]byte, error) {
	var err error
	var q Request
	var payload []byte

	q = Request{Command: commandName, Arguments: arguments}

	if payload, err = json.Marshal(q); err != nil {
		log.Fatalln("Error on JSON Marshal: ", err)
		return nil, err
	}

	return payload, nil
}

// Creates a struct representation of informed slice of bytes, which by default
// validate data structure against Response type.
func NewResponse(payload []byte) (*Response, error) {
	var err error
	var resp *Response = &Response{}

	if err = json.Unmarshal(payload, resp); err != nil {
		log.Fatalln("Payload:", string(payload[:]))
		log.Fatalln("Error on JSON marchal:", err)
		return nil, err
	}

	return resp, nil
}

/* EOF */
