package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/igorpollo/go-custom-log"
	"golang.org/x/crypto/ssh"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

type Server struct {
	Commands []string `yaml:"commands"`
	Config   string   `yaml:"config"`
	Region   string   `yaml:"region"`
	Image    string   `yaml:"image"`
	Provider string   `yaml:"provider"`
}

type CreatedServer struct {
	ID        int    `yaml:"id"`
	CreatedAt string `yaml:"created_at"`
	IP        string `yaml:"ip"`
}
type Configs struct {
	Main struct {
		Port string `yaml:"port"`
	} `yaml:"main"`
	Servers     map[string]Server
	CreatedByIP map[string]CreatedServer
}

type CreateServer struct {
	ServerName string `json:"server_name"`
}

type DeleteIP struct {
	IP string `json:"ip"`
}


 var pat = os.Getenv("DO_TOKEN")


type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func readConfig() Configs {
	configData, _ := ioutil.ReadFile("./config.yml")
	config := Configs{}
	err := yaml.Unmarshal([]byte(configData), &config)
	if err != nil {
		fmt.Println(err)
	}
	return config
}

func createServer(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var s CreateServer
	err := decoder.Decode(&s)
	if err != nil {
		panic(err)
	}
	config := readConfig()
	if _, ok := config.Servers[s.ServerName]; !ok {
		w.Write([]byte("This server doesn't exist."))
		return
	}
	server := config.Servers[s.ServerName]
	switch server.Provider {
	case "DigitalOcean":
		err := createDigitalOcean(server, s.ServerName)
		if err != nil {
			w.Write([]byte("Error creating the server"))
			return
		}
		w.Write([]byte("Server created successfully"))
	case "Google":
		err := createGoogleCloud(server, s.ServerName)
		if err != nil {
			w.Write([]byte("Error creating the server"))
			return
		}
		w.Write([]byte("Server created successfully"))
	}
}

func runCommand(cmd string, conn *ssh.Client) {
	sess, err := conn.NewSession()
	if err != nil {
		panic(err)
	}

	defer sess.Close()

	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		panic(err)
	}

	go io.Copy(os.Stdout, sessStdOut)

	sessStderr, err := sess.StderrPipe()
	if err != nil {
		panic(err)
	}

	go io.Copy(os.Stderr, sessStderr)

	err = sess.Run(cmd) // eg., /usr/bin/whoami
	if err != nil {
		panic(err)
	}
}

func hai(path string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	signer, err := ssh.ParsePrivateKeyWithPassphrase(key, []byte(os.Getenv("SSH_PASSWORD")))
	if err != nil {
		panic(err)
	}
	return ssh.PublicKeys(signer)
}

func deleteIP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var s DeleteIP
	err := decoder.Decode(&s)
	if err != nil {
		
	}
	if s.IP == "" {
		ip := strings.Split(req.RemoteAddr, ":")
		s.IP = ip[0]
	}
	err = DeleteByIP(s.IP)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("Error deleting droplet"))
		return
	}
	w.Write([]byte("Droplet deleted succesfully"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/create", createServer).Methods("POST")
	r.HandleFunc("/deletebyip", deleteIP).Methods("POST")

	http.Handle("/", r)
	log.Info("🚀 Server Started")
	http.ListenAndServe(":5000", nil)
}
