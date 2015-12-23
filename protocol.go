package godutch

import (
	"encoding/json"
	"log"
)

type Request struct {
	Command   string   `json:"command"`
	Arguments []string `json:"arguments"`
}

type Response struct {
	Name   string   `json:"name"`
	Status int      `json:"status"`
	Stdout []string `json:"stdout"`
	Error  string   `json:"error"`
}

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

func NewResponse(payload []byte) (*Response, error) {
	var err error
	var resp *Response = &Response{}

	if err = json.Unmarshal(payload, resp); err != nil {
		log.Fatalln("Error on JSON marchal:", err)
		return nil, err
	}

	return resp, nil
}

/* EOF */
