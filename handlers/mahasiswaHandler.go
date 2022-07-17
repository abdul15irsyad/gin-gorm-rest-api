package handlers

import (
	"belajar-gin/structs"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func loadData(mahasiswas *[]structs.Mahasiswa) {
	// load data from json file
	file, err := ioutil.ReadFile("./data/mahasiswa.json")
	if err != nil {
		fmt.Println(err)
	}
	_ = json.Unmarshal([]byte(file), &mahasiswas)
}

func GetAllMahasiswa(c *gin.Context) {
	mahasiswas := []structs.Mahasiswa{}
	loadData(&mahasiswas)

	c.JSON(http.StatusOK, gin.H{
		"msg":  "get all mahasiswa",
		"data": mahasiswas,
	})
}

func findMahasiswa(id int, mahasiswas []structs.Mahasiswa) structs.Mahasiswa {
	result := structs.Mahasiswa{}
	for _, mahasiswa := range mahasiswas {
		if mahasiswa.Id == id {
			result = mahasiswa
		}
	}
	return result
}

func GetMahasiswa(c *gin.Context) {
	// convert string to int
	id, _ := strconv.Atoi(c.Param("id"))

	mahasiswas := []structs.Mahasiswa{}
	loadData(&mahasiswas)

	c.JSON(http.StatusOK, gin.H{
		"msg":  "get mahasiswa",
		"data": findMahasiswa(id, mahasiswas),
	})
}
