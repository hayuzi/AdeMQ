package remote

import "encoding/json"

type Formatted struct {
	Cmd    string   `json:"cmd"`
	Params []string `json:"params"`
}

func FormatRequest(cmd string, params []string) ([]byte, error) {
	data := Formatted{
		Cmd:    cmd,
		Params: params,
	}
	return json.Marshal(data)
}
