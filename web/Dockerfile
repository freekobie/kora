FROM node:22-alpine

WORKDIR /frontend

COPY package.json .
COPY package-lock.json .

RUN npm run install

COPY . .

EXPOSE 5173

CMD ["npm", "start"]

