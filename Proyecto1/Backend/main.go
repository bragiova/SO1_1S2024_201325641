package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

var conexion = ConexionMysql()

type DataUsage struct {
	Ram float32 `json:"ram_porcentaje"`
	Cpu float32 `json:"cpu_porcentaje"`
}

type Childs struct {
	Pid   int    `json: "pid"`
	Name  string `json: "name"`
	State int    `json: "state"`
}

type Process struct {
	Pid   int      `json: "pid"`
	Name  string   `json: "name"`
	User  int      `json: "user"`
	State int      `json: "state"`
	Ram   int      `json: "ram"`
	Child []Childs `json: "child"`
}

type CpuModel struct {
	Processes []Process `json: "processes"`
	Running   int       `json: "running"`
	Sleeping  int       `json: "sleeping"`
	Zombie    int       `json: "zombie"`
	Stopped   int       `json: "stopped"`
	Total     int       `json: "total"`
	Cpu       float32   `json: "cpu"`
}

type RamModel struct {
	Libre float32 `json: "libre"`
	Uso   float32 `json "uso"`
}

type CpuRamHistorical struct {
	Percentage float32 `json: "percentage"`
	Date       string  `json "date"`
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})
	origins := handlers.AllowedOrigins([]string{"*"})
	// router.Use(mux.CORSMethodMiddleware(router))

	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/realtime", cpuRamRealTime).Methods("GET")
	router.HandleFunc("/historicalram", getRamHistorical).Methods("GET")
	router.HandleFunc("/historicalcpu", getCpuHistorical).Methods("GET")
	fmt.Println("Server running in 3000 port")
	go cpuRamRoutineDB()
	http.ListenAndServe(":3000", handlers.CORS(headers, methods, origins)(router))
}

func home(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("Hello from Go api"))
}

func cpuRamRealTime(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	var objCpu CpuModel = readCpuInfo()
	var ramInfo = readRamInfo()

	var dataReal DataUsage

	dataReal.Cpu = objCpu.Cpu
	dataReal.Ram = ramInfo

	json.NewEncoder(response).Encode(dataReal)
}

func cpuRamRoutineDB() {
	// Lanzar un goroutine que ejecute la función cada n segundos
	interval := 10 // segundos
	ticker := time.NewTicker(time.Second * time.Duration(interval))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			var objCpu CpuModel = readCpuInfo()
			if objCpu.Cpu > 0 {
				saveDBCpuRamInfo(objCpu.Cpu, true)
			}

			var ramInfo = readRamInfo()
			if ramInfo >= 0 {
				saveDBCpuRamInfo(ramInfo, false)
			}
		}
	}
}

func readCpuInfo() CpuModel {
	fmt.Println("DATOS OBTENIDOS DESDE EL MODULO CPU:")
	fmt.Println("")
	var cpu_obj CpuModel

	cmd := exec.Command("sh", "-c", "cat /proc/cpu_so1_1s2024")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error al ejecutar el comando:", err)
		fmt.Println("Mensaje de error:", stderr.String())
		return cpu_obj
	}

	if err := json.Unmarshal(stdout.Bytes(), &cpu_obj); err != nil {
		fmt.Println("Error al decodificar JSON:", err)
		return cpu_obj
	}
	// fmt.Println(cpu_obj.Cpu)
	return cpu_obj
}

func saveDBCpuRamInfo(percentage float32, isCpu bool) {
	var timeNow = time.Now().Format("2006-01-02 15:04:05")

	tableName := (map[bool]string{true: "cpu_info", false: "ram_info"})[isCpu]

	query := "INSERT INTO " + tableName + "(porcentaje, fecha) VALUES (?,?);"
	result, err := conexion.Exec(query, percentage, timeNow)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)
}

func readRamInfo() float32 {
	fmt.Println("DATOS OBTENIDOS DESDE EL MODULO RAM:")
	fmt.Println("")
	var ram_obj RamModel

	cmd := exec.Command("sh", "-c", "cat /proc/ram_so1_1s2024")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error al ejecutar el comando:", err)
		fmt.Println("Mensaje de error:", stderr.String())
		return 0
	}

	if err := json.Unmarshal(stdout.Bytes(), &ram_obj); err != nil {
		fmt.Println("Error al decodificar JSON:", err)
		return 0
	}
	// fmt.Println(ram_obj)

	totalRam := ram_obj.Libre + ram_obj.Uso
	percentageRam := (ram_obj.Uso / totalRam) * 100

	return percentageRam
}

func ConexionMysql() *sql.DB {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	conexion, err := sql.Open("mysql", connString)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Conexion con MySQL Correcta")
	}
	return conexion
}

func getRamHistorical(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var lista []CpuRamHistorical
	query := "select porcentaje, fecha from ram_info;"
	result, err := conexion.Query(query)
	if err != nil {
		fmt.Println(err)
	}

	for result.Next() {
		var logc CpuRamHistorical

		err = result.Scan(&logc.Percentage, &logc.Date)
		if err != nil {
			fmt.Println(err)
		}
		lista = append(lista, logc)
	}
	json.NewEncoder(response).Encode(lista)
}

func getCpuHistorical(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var lista []CpuRamHistorical
	query := "select porcentaje, fecha from cpu_info;"
	result, err := conexion.Query(query)
	if err != nil {
		fmt.Println(err)
	}

	for result.Next() {
		var logc CpuRamHistorical

		err = result.Scan(&logc.Percentage, &logc.Date)
		if err != nil {
			fmt.Println(err)
		}
		lista = append(lista, logc)
	}
	json.NewEncoder(response).Encode(lista)
}
