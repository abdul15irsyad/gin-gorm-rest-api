package structs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Mahasiswas struct {
	Mahasiswas []Mahasiswa
}

// load data from json file
func (ms Mahasiswas) LoadData() Mahasiswas {
	file, err := ioutil.ReadFile("./data/mahasiswa.json")
	if err != nil {
		fmt.Println(err)
	}
	_ = json.Unmarshal([]byte(file), &ms.Mahasiswas)
	return ms
}

// find mahasiswa
func (ms Mahasiswas) FindMahasiswa(callback func(m Mahasiswa) bool) Mahasiswa {
	result := Mahasiswa{}
	for _, mahasiswa := range ms.Mahasiswas {
		if callback(mahasiswa) {
			result = mahasiswa
			break
		}
	}
	return result
}

type Mahasiswa struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Year int    `json:"year"`
}
