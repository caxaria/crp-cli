package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"

	pass "code.google.com/p/gopass"
)

const AUTH_ADDRESS = ADDRESS + "/tokens"
const TOKEN_PATH = "/.crowdprocess/token"

func login(username, password string) error {
	var err error

	// username not provided, ask for it
	if username == "" {
		_, err = fmt.Print("Email: ")
		if err != nil {
			panic(err.Error())
		}

		var user string
		_, err = fmt.Scanln(&user)
		if err != nil {
			panic(err.Error())
		}
		username = user
	}

	// password not provided, ask for it
	if password == "" {
		password, err = pass.GetPass("Password: ")
		if err != nil {
			panic(err.Error())
		}
	}

	request, err := http.NewRequest("POST", AUTH_ADDRESS, nil)
	if err != nil {
		panic(err.Error())
	}

	request.SetBasicAuth(username, password)

	response, err := client.Do(request)
	if err != nil {
		panic(err.Error())
	}

	if response.StatusCode == 401 {
		return errors.New("Authentication failed")
	}

	if response.StatusCode != 201 {
		return errors.New("Something went wrong")
	}

	var token struct {
		Id string
	}
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&token)
	if err != nil {
		panic(err.Error())
	}

	usr, err := user.Current()
	if err != nil {
		panic(err.Error())
	}

	err = ioutil.WriteFile(usr.HomeDir+TOKEN_PATH, []byte(token.Id), 0600)
	if err != nil {
		panic(err.Error())
	}

	return nil
}

func logout() error {
	var err error

	usr, err := user.Current()
	if err != nil {
		panic(err.Error())
	}

	tokenPath := usr.HomeDir + TOKEN_PATH

	token, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return errors.New("Already logged out")
	}

	request, err := http.NewRequest("DELETE", AUTH_ADDRESS+"/"+string(token), nil)
	if err != nil {
		panic(err.Error())
	}

	authenticateRequest(request)

	response, err := client.Do(request)
	if err != nil {
		panic(err.Error())
	}

	if response.StatusCode == 401 {
		return errors.New("Authentication failed")
	}

	if response.StatusCode != 204 {
		return errors.New("Something went wrong")
	}

	err = os.Remove(tokenPath)
	if err != nil {
		panic(err.Error())
	}

	return nil
}

func authenticateRequest(request *http.Request) {
	var err error

	// username + password
	if USERNAME != "" {
		if PASSWORD == "" {
			PASSWORD, err = pass.GetPass("Password: ")
			if err != nil {
				panic(err.Error())
			}
		}

		request.SetBasicAuth(USERNAME, PASSWORD)

		return
	}

	// token
	usr, err := user.Current()
	if err != nil {
		panic(err.Error())
	}

	tokenPath := usr.HomeDir + TOKEN_PATH

	token, err := ioutil.ReadFile(tokenPath)
	if err == nil {
		request.Header.Add("Authorization", "Token "+string(token))

		return
	}
}