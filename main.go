package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Router struct {
	rules map[string]http.HandlerFunc
}
type Server struct {
	port   string
	router *Router
}

type Info struct {
	Results []struct {
		Login struct {
			UUID     string `json:"uuid"`
			Username string `json:"username"`
			Password string `json:"password"`
			Salt     string `json:"salt"`
			Md5      string `json:"md5"`
			Sha1     string `json:"sha1"`
			Sha256   string `json:"sha256"`
		} `json:"login"`
	} `json:"results"`
	Info struct {
		Seed    string `json:"seed"`
		Results int    `json:"results"`
		Page    int    `json:"page"`
		Version string `json:"version"`
	} `json:"info"`
}

type Personas struct {
	Results []struct {
		Gender string
		Name   struct {
			Title string
			First string
			Last  string
		}
		Email string
		Login struct {
			UUID     string
			Username string
			Password string
			Salt     string
			Md5      string
			Sha1     string
			Sha256   string
		}
	}
	Info struct {
		Seed    string
		Results int
		Page    int
		Version string
	}
}

func main() {
	server := NewServer(":3000")
	server.Listen()
}

func (route *Router) Handle(w http.ResponseWriter, r *http.Request) {
	Handler, ok := route.rules[r.URL.Path]

	if ok {
		Handler(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func GetPersonas(w http.ResponseWriter, r *http.Request) {
	var personas Personas
	var info Info

	resp, err := http.Get("https://randomuser.me/api/?inc=name,email,gender,login&results=5000")
	if err != nil {
		log.Fatalln(err)
	}
	time.Sleep(5 * time.Second)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(body, &personas)
	if err != nil {
		fmt.Println(err)
		return
	}
	seen := make(map[string]bool)

	for i := 0; i < len(personas.Results); i++ {
		uuid := personas.Results[i].Login.UUID
		if seen[uuid] {
			resp, err := http.Get("https://randomuser.me/api/?inc=login")
			if err != nil {
				log.Fatalln(err)
			}
			time.Sleep(5 * time.Second)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			err = json.Unmarshal(body, &info)
			if err != nil {
				fmt.Println(err)
				return
			}
			personas.Results[i].Login.UUID = info.Results[0].Login.UUID
		}
	}

	for i := 0; i < len(personas.Results); i++ {
		fmt.Fprintf(w, fmt.Sprintf("Nombre: %s, Correo: %s, UUID: %s \n", personas.Results[i].Name, personas.Results[i].Email, personas.Results[i].Login.UUID))
	}

}

func NewServer(port string) *Server {
	return &Server{
		port:   port,
		router: NewRouter(),
	}
}

func (server *Server) Listen() error {
	http.HandleFunc("/", server.router.Handle)

	err := http.ListenAndServe(server.port, nil)
	if err != nil {
		return err
	}
	return nil
}

func NewRouter() *Router {
	router := &Router{
		rules: make(map[string]http.HandlerFunc),
	}
	router.rules["/getPersonas"] = GetPersonas

	return router
}

func Existe(info []int, dest int) int {
	c := 0
	i := -1

	for i, num := range info {
		if num == dest {
			c++
			if c == 2 {
				i = i
				break
			}
		}
	}
	return i
}
