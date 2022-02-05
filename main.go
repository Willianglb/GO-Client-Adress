package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Pessoa struct {
	Id 				string 		`json:"id"`
	PrimeiroNome 	string 		`json:"nome"`
	Sobrenome 		string 		`json:"sobrenome"`
	Endereco 		Endereco 	`json:"enderecos"`
}

type Endereco struct {
	Id 			string 	`json:"endereco_id"`
	Logradouro string `json:"logradouro"`
	Cep        string    `json:"cep"`
	Bairro     string `json:"bairro"`
	Cidade 		string 	`json:"cidade"`
	UF 			string 	`json:"uf"`
}

type ClienteCidade struct {
	Cidade 		string 	`json:"cidade"`
	UF 			string 	`json:"uf"`
	ClientePessoa []ClientePessoa `json:"clientes"`
}

type ClientePessoa struct {
	Id 				string 		`json:"id"`
	Nome 	string 		`json:"nome"`
}

var pessoas []Pessoa

func main() {
	getInformationsJSON()

	http.HandleFunc("/pessoas", getAllPessoas)

	http.HandleFunc("/cidade/", buscarPessoaPorCidade)

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "ping")

	if err != nil {
		log.Fatal(err)
	}
}

func readPessoasFromArchive() []byte {
	byteValueJSON, err := ioutil.ReadFile("./db.json")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully readed file in bytes")

	return byteValueJSON
}

func getInformationsJSON() {
	byteValueJSON := readPessoasFromArchive()

	err := json.Unmarshal(byteValueJSON, &pessoas)

	if err != nil {
		log.Fatal(err)
	}
}

func getAllPessoas(w http.ResponseWriter, r *http.Request) {

	err := json.NewEncoder(w).Encode(pessoas)

	if err != nil {
		fmt.Fprint(w, err)
	}

}

func buscarPessoaPorCidade(w http.ResponseWriter, r *http.Request) {

	partsOfURL := strings.Split(r.URL.Path, "/")

	dadoRecebido := partsOfURL[2]

	if len(partsOfURL) > 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	filteredPessoas := ClienteCidade{}

	for _, cidade := range pessoas {
		if dadoRecebido == cidade.Endereco.Cidade {
			filteredPessoas.Cidade = cidade.Endereco.Cidade
			filteredPessoas.UF = cidade.Endereco.UF
			cliente := ClientePessoa{}
			cliente.Id = cidade.Id
			cliente.Nome = cidade.PrimeiroNome + " " + cidade.Sobrenome
			filteredPessoas.ClientePessoa = append(filteredPessoas.ClientePessoa, cliente)
		}
	}

	err := json.NewEncoder(w).Encode(filteredPessoas)

	if err != nil {
		log.Fatal(err)
	}
}
