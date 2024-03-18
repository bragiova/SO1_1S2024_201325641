# Imagen oficial node
FROM node:latest as build

# Directorio de trabajo
WORKDIR /app

# Copiar package.json y package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the entire application to the container
COPY . .

# Build the React app for production
RUN npm run build

# NGINX
FROM nginx:1.21.4-alpine

# Copiar configuración de nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf
# COPY nginx.conf //etc/nginx/nginx.conf

# Copiar archivos contruidos antes de nginx
COPY --from=build /app/dist /usr/share/nginx/html

# Expone puerto 80
EXPOSE 80

# Comando para iniciar nginx
CMD ["nginx", "-g", "daemon off;"]