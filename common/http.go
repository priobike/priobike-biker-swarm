package common

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func Get(url string, serviceName string, headers map[string]string) []byte {
	client := http.Client{
		Timeout: Timeout(),
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(serviceName + ": " + err.Error())
	}
	for key, value := range headers {
		request.Header.Add(key, value)
	}
	response, requestErr := client.Do(request)
	if requestErr != nil {
		panic(serviceName + ": " + requestErr.Error())
	}
	defer response.Body.Close()
	fmt.Println(serviceName+" status:", response.Status)
	if response.StatusCode != 200 {
		io.Copy(os.Stdout, response.Body)
		panic(serviceName + ": request failed")
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(serviceName + ": " + err.Error())
	}
	return body
}

func PostJson(url string, serviceName string, requestBody io.Reader) []byte {
	client := http.Client{
		Timeout: Timeout(),
	}
	response, requestErr := client.Post(url, "application/json", requestBody)
	if requestErr != nil {
		panic(serviceName + ": " + requestErr.Error())
	}
	defer response.Body.Close()
	fmt.Println(serviceName+" status:", response.Status)
	if response.StatusCode != 200 {
		io.Copy(os.Stdout, response.Body)
		panic(serviceName + ": request failed")
	}
	responseBody, responseErr := io.ReadAll(response.Body)
	if responseErr != nil {
		panic(serviceName + ": " + requestErr.Error())
	}
	return responseBody
}
