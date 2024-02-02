# Lenguaje o imagen
FROM golang:1.21.4

# Directorio virtual
WORKDIR /tarea1

# Descarga de mod Go
COPY go.mod go.sum ./
RUN go mod download

# Copia del código
COPY *.go ./

#Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /api_tarea1

EXPOSE 3000

#RUN
CMD ["/api_tarea1"]