package main

import (
	"errors"
	"io/ioutil"
	"os/user"

	"golang.org/x/crypto/ssh"
)

func getKeyFromFile(file string) (key ssh.Signer, err error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil,
			errors.New("Couldn't read from the file: '" +
				file + "': " + err.Error())
	}
	return ssh.ParsePrivateKey(buf)
}

func getKey(usr *user.User) (key ssh.Signer, err error) {
	return getKeyFromFile(usr.HomeDir + "/.ssh/id_rsa")
}

func getClientConfiguration() (*ssh.ClientConfig, error) {
	usr, err := user.Current()
	if err != nil {
		err = errors.New("Couldn't retrieve the current user: " +
			err.Error())
		return nil, err
	}

	key, err := getKey(usr)
	if err != nil {
		err = errors.New("Couldn't get the ssh key for user: '" +
			usr.Username + "': " +
			err.Error())
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: usr.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}

	return config, nil
}

func makeClientConn(host string) (*ssh.Client, error) {
	config, err := getClientConfiguration()
	if err != nil {
		err = errors.New("Couldn't get the client config: " +
			err.Error())
		return nil, err
	}

	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		err = errors.New("Couldn't connect to host: '" +
			host + "' as user: '" +
			config.User + "':" +
			err.Error())
		return nil, err
	}
	return client, nil
}
