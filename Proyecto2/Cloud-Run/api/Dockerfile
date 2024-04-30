# Define la imagen base de Node.js 18.16.0 en su versión "alpine"
FROM node:18.16.0-alpine

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

EXPOSE 5000

CMD ["node", "index.js"]
