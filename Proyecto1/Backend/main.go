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

func main() {
	router := mux.NewRouter().StrictSlash(true)
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})
	origins := handlers.AllowedOrigins([]string{"*"})
	// router.Use(mux.CORSMethodMiddleware(router))

	router.HandleFunc("/", home).Methods("GET")
	fmt.Println("Server running in 3000 port")
	readCpuInfo()
	http.ListenAndServe(":3000", handlers.CORS(headers, methods, origins)(router))
}

func home(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("Hello from Go api"))
}

func cpuRamRoutine() {
	// Lanzar un goroutine que ejecute la función cada n segundos
	interval := 10 // segundos
	ticker := time.NewTicker(time.Second * time.Duration(interval))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			readCpuInfo()
		}
	}
}

func readCpuInfo() {
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
		return
	}

	if err := json.Unmarshal(stdout.Bytes(), &cpu_obj); err != nil {
		fmt.Println("Error al decodificar JSON:", err)
		return
	}
	fmt.Println(cpu_obj.Cpu)
}

func readRamInfo() {
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
		return
	}

	if err := json.Unmarshal(stdout.Bytes(), &cpu_obj); err != nil {
		fmt.Println("Error al decodificar JSON:", err)
		return
	}
	fmt.Println(cpu_obj.Cpu)
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
