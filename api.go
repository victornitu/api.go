package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type (
	API struct {
		Client        *http.Client
		baseUrl       string
		Authorization string
	}
	Status struct {
		Code int `json:"code"`
		Message string `json:"message"`
	}
)

func (status Status) IsError() bool {
	return status.Code >= 400
}

func New(baseUrl string) API {
	return API{&http.Client{Timeout: 10 * time.Second}, baseUrl, ""}
}

func (api *API) Get(url string, target interface{}) (status Status, err error) {
	return api.request("GET", url, nil, target)
}

func (api *API) Head(url string, target interface{}) (status Status, err error) {
	return api.request("GET", url, nil, target)
}

func (api *API) Post(url string, body interface{}, target interface{}) (status Status, err error) {
	return api.push("POST", url, body, target)
}

func (api *API) Put(url string, body interface{}, target interface{}) (status Status, err error) {
	return api.push("PUT", url, body, target)
}

func (api *API) Patch(url string, body interface{}, target interface{}) (status Status, err error) {
	return api.push("PATCH", url, body, target)
}

func (api *API) push(method string, url string, body interface{}, target interface{}) (status Status, err error) {
	b, err := json.Marshal(body)
	if err != nil {
		return
	}
	return api.request(method, url, bytes.NewBuffer(b), target)
}

func (api *API) request(method string, url string, body io.Reader, target interface{}) (status Status, err error) {
	req, err := http.NewRequest(method, api.baseUrl+url, body)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", api.Authorization)
	req.Header.Set("Content-Type", "application/json")
	res, err := api.Client.Do(req)
	if err != nil {
		return
	}
	status = Status{res.StatusCode, ""}
	switch {
	case status.Code / 100 == 4:
		status.Message = "external service: request failed"
		return
	case status.Code / 100 == 5:
		status.Message = "external service: unavailable"
		return
	}
	if target == nil {
		return
	}
	err = json.NewDecoder(res.Body).Decode(target)
	return
}


