# Imagen base
FROM golang:latest

# Install mpstat
RUN apt-get update && apt-get install -y sysstat

# Directorio de trabajo dentro del contenedor
WORKDIR /api

# Copiar los archivos go.mod y go.sum
COPY go.mod go.sum ./

# Descargar todas las dependencias
RUN go mod download

# Copiar el código fuente desde el directorio actual al contenedor de Docker
COPY . .

# Compilar la aplicación
RUN go build -o main .

# Exponer el puerto por el que la aplicación será accesible
EXPOSE 3000

# Comando para ejecutar la aplicación
CMD ["./main"]