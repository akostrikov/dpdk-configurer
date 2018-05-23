package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func getHostKey(host string) (ssh.PublicKey, error) {
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				return nil, errors.New(fmt.Sprintf("error parsing %q: %v", fields[2], err))
			}
			break
		}
	}

	if hostKey == nil {
		return nil, errors.New(fmt.Sprintf("no hostkey for %s", host))
	}
	return hostKey, nil
}

func renderConnectFormHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("connect_form.html")
	t.Execute(w, nil)
}

func connectHandler(w http.ResponseWriter, r *http.Request) {
	var hostKey ssh.PublicKey
	address := r.FormValue("address")

	hostKey, err := getHostKey(address)
	if err != nil {
		log.Fatal(err)
	}

	key, err := ioutil.ReadFile("/root/.ssh/id_rsa")
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	client, err := ssh.Dial("tcp", address+":22", config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("/usr/bin/whoami"); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())

	t, _ := template.ParseFiles("connect.html")
	t.Execute(w, nil)
}

func main() {
	http.HandleFunc("/connect_form/", renderConnectFormHandler)
	http.HandleFunc("/connect_handler/", connectHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
