# Proyecto 1

## Base de Datos
Se creó la base de datos en MySql con el siguiente script:

![script](./imagenes/bd1.png)

En el cual se crea la base, dos tablas para el monitoreo de ram y cpu

![bd](./imagenes/bd2.png)

## Módulos Kernel
Se crearon dos módulos de kernel para obtener información de procesos y porcentajes de uso tanto de la RAM como del CPU.

### RAM
Importación de cabeceras o librerías necesarias

![ram1](./imagenes/ram1.png)

Se devuelve un json con la información

![ram2](./imagenes/ram2.png)

Finalmente se crea el archivo con la información obtenida

![ram4](./imagenes/ram4.png)

### CPU
Importación de cabeceras o librerías necesarias

![cpu1](./imagenes/cpu1.png)

Se obtienen los procesos y se forma el json a devolver

![cpu2](./imagenes/cpu2.png)

![cpu3](./imagenes/cpu3.png)

Finalmente se crea el archivo con la información obtenida

![cpu4](./imagenes/cpu4.png)

## Golang
Se realiza la conexión a la BD

![gobd1](./imagenes/gobd1.png)

Para insertar la información en la BD, se ejecuta el comando

![gobd2](./imagenes/gobd2.png)

Para obtener el histórico, se realiza un select a las tablas

![gobd3](./imagenes/gobd3.png)

Para leer la información de los módulos, se ejecuta el comando cat de la siguiente manera

![go1](./imagenes/go1.png)

Toda la información la vamos a serializar en JSON para la comunicación con el frontend

Se tiene una rutina ejecutándose cada 10 segundos que se encarga de almacenar la información en la BD para el monitoreo histórico

![go2](./imagenes/go2.png)

Finalmente tenemos nuestras rutas y la llamada a la rutina antes mencionada

![go3](./imagenes/go3.png)

## FrontEnd
Una vez iniciada la aplicación, nos mostrará de inicio el monitoreo en tiempo real. Este monitoreo se muestra con una gráficas de Pie indicando el porcentaje de Uso y Libre, tanto de la RAM como del CPU

![front1](./imagenes/front1.png)

En el caso del monitoreo histórico, nos muestra unas gráficas de líneas con la fecha y porcentaje en cada eje.

![front2](./imagenes/front2.png)

