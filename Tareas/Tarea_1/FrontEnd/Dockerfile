FROM node:18.17.1-alpine

WORKDIR /frontend

COPY package.json ./
COPY package-lock.json ./

RUN npm install

COPY . ./

EXPOSE 5173

CMD ["npm", "run", "dev"]
