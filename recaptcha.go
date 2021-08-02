package recaptcha

import(
	"encoding/json"
	"net/http"
	"strings"
	"net/url"
	"errors"
	"fmt"
)

type Client struct {
	PrivateKey	string
	HttpClient	*http.Client
}

type Response struct{
	Success		bool		`json:"success"`
	HostName	string		`json:"hostname"`
	ErrorCodes	[]string	`json:"error-codes"`
}

// Preparation client to use package
func CaptchaClient(PrivateKey string) *Client{
	return &Client{
		PrivateKey: PrivateKey,
		HttpClient: http.DefaultClient,
	}
}

// Verify response
func (client *Client) VerifyResponse(response, ip string) (bool, string, error){
	// Preparation form data
	data := url.Values{
		"secret":		{client.PrivateKey},
		"response":		{response},
		"remoteip":		{ip},
	}

	// Preparation request
	req, err := http.NewRequest("POST", "https://www.google.com/recaptcha/api/siteverify", strings.NewReader(data.Encode()))
	if err != nil {
		return false, "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	// Decode the response
	var responseBody Response
	json.NewDecoder(resp.Body).Decode(&responseBody)

	if len(responseBody.ErrorCodes) > 0 {
		return false, "", errors.New(fmt.Sprint("Error code: ", responseBody.ErrorCodes[0]))
	}

	return responseBody.Success, responseBody.HostName, err
}