package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/goburrow/modbus"
	"github.com/olebedev/config"
)

type test_struct struct {
	Register int `json:"register"`
	Set      int `json:"set"`
}

type Data struct {
	client    modbus.Client
	start_bit int
	stop_bit  int
}

func (main_data *Data) datas(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	infos := make(map[int]float64)
	for i := main_data.start_bit; i <= main_data.stop_bit; i++ {
		result, err := main_data.client.ReadInputRegisters(uint16(i), 1)
		if err != nil {
			panic(err)
		}
		value := binary.BigEndian.Uint16(result)
		value_fl64 := float64(value) / 10
		infos[i] = value_fl64
	}

	jsonResp, err := json.Marshal(infos)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonResp)
}

func (main_data *Data) action(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	var m test_struct
	json.Unmarshal(body, &m)
	result, err := main_data.client.WriteSingleRegister(uint16(m.Register), uint16(m.Set))
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func main() {
	main_data := new(Data)

	file, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}
	yamlString := string(file)

	cfg, err := config.ParseYaml(yamlString)
	if err != nil {
		panic(err)
	}

	server_url, err := cfg.String("server_url")
	if err != nil {
		panic(err)
	}

	main_data.start_bit, err = cfg.Int("start_bit")
	if err != nil {
		panic(err)
	}

	main_data.stop_bit, err = cfg.Int("stop_bit")
	if err != nil {
		panic(err)
	}

	port, err := cfg.String("port")
	if err != nil {
		panic(err)
	}

	baudrate, err := cfg.Int("baudrate")
	if err != nil {
		panic(err)
	}

	databits, err := cfg.Int("databits")
	if err != nil {
		panic(err)
	}

	parity, err := cfg.String("parity")
	if err != nil {
		panic(err)
	}

	stopbits, err := cfg.Int("stopbits")
	if err != nil {
		panic(err)
	}

	slaveid, err := cfg.Int("slaveid")
	if err != nil {
		panic(err)
	}

	slaveid_b := byte(slaveid)

	handler := modbus.NewRTUClientHandler(port)
	handler.BaudRate = baudrate
	handler.DataBits = databits
	handler.Parity = parity
	handler.StopBits = stopbits
	handler.SlaveId = slaveid_b
	handler.Timeout = 5 * time.Second

	err_conn := handler.Connect()
	if err_conn != nil {
		panic(err)
	}

	main_data.client = modbus.NewClient(handler)
	defer handler.Close()

	fmt.Println("Config read OK")
	fmt.Println("Create client OK")

	html := http.FileServer(http.Dir("./dist"))

	http.HandleFunc("/datas", main_data.datas)
	http.HandleFunc("/action", main_data.action)
	http.Handle("/", html)

	err = http.ListenAndServe(server_url, nil)
	if err != nil {
		log.Fatal(err)
	}
}
