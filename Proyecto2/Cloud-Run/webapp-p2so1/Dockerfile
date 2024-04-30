# Utilizar la imagen de Node.js como base
FROM node:latest

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar los archivos package.json y package-lock.json
COPY package*.json ./

# Instalar las dependencias
RUN npm install

# Copiar el resto de los archivos de la aplicación
COPY . .

# Construir la aplicación Vue.js para producción
RUN npm run build

# Utilizar la imagen de Nginx como base para servir la aplicación
FROM nginx:stable-alpine

# Copiar los archivos de construcción de Vue.js al directorio de Nginx
COPY --from=0 /app/dist /usr/share/nginx/html

# Exponer el puerto en el que se ejecutará la aplicación
EXPOSE 80

# Comando para iniciar Nginx
CMD ["nginx", "-g", "daemon off;"]