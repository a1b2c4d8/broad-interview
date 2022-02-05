package client

type StopsResponse struct {
	Stops []*StopResponseData `json:"data"`
}

type StopResponseData struct {
	Id             string                  `json:"id"`
	StopAttributes *StopResponseAttributes `json:"attributes"`
}

type StopResponseAttributes struct {
	Name string `json:"name"`
}
