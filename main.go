package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Server struct {
	ServerName    string
	ServerUrl     string
	Status        int
	TempoExecucao float64
	DataFalha     string
}

func criarListaServidores(serverList *os.File) []Server {
	csvReader := csv.NewReader(serverList)
	data, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var servidores []Server
	for i, line := range data {
		if i > 0 {
			servidor := Server{
				ServerName: line[0],
				ServerUrl:  line[1],
			}
			servidores = append(servidores, servidor)
		}
	}
	return servidores
}

func checkServer(servidores []Server) []Server {
	var downServers []Server

	for _, servidor := range servidores {
		agora := time.Now()
		resp, err := http.Get(servidor.ServerUrl)
		if err != nil {
			fmt.Printf("Server: [%s] %v\n", servidor.ServerName, err)
			servidor.Status = 0
			servidor.DataFalha = agora.Format("02/01/2006 15:04:05")
			servidor.TempoExecucao = time.Since(agora).Seconds()
			downServers = append(downServers, servidor)
			continue
		}

		resp.Body.Close()

		servidor.Status = resp.StatusCode
		if servidor.Status != 200 && servidor.Status != 201 {
			servidor.DataFalha = agora.Format("02/01/2006 15:04:05")
			downServers = append(downServers, servidor)
		}
		servidor.TempoExecucao = time.Since(agora).Seconds()
		fmt.Printf("Status: [%d] Tempo de Carga: [%f] URL: [%s]\n", servidor.Status, servidor.TempoExecucao, servidor.ServerUrl)
	}

	return downServers
}

func openfiles(serverListFile string, downTimeFile string) (*os.File, *os.File) {
	serverList, err := os.Open(serverListFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	downTimeList, err := os.OpenFile(downTimeFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return serverList, downTimeList
}

func generateDownTime(downTimeList *os.File, downServers []Server) {
	csvWriter := csv.NewWriter(downTimeList)
	for _, servidor := range downServers {
		line := []string{servidor.ServerName, servidor.DataFalha, fmt.Sprintf("%f", servidor.TempoExecucao), fmt.Sprintf("%d", servidor.Status)}
		csvWriter.Write(line)
	}
	csvWriter.Flush()
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Uso: programa <server_list.csv> <downtime.csv>")
		os.Exit(1)
	}
	serverList, downTimeList := openfiles(os.Args[1], os.Args[2])
	defer serverList.Close()
	defer downTimeList.Close()
	servidores := criarListaServidores(serverList)
	downServers := checkServer(servidores)
	generateDownTime(downTimeList, downServers)
}
