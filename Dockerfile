FROM node:19-slim as update-install-git
WORKDIR /usr/src/app
RUN apt-get update
RUN apt-get install git -y

FROM update-install-git as frontend
WORKDIR /usr/src/app

COPY ./web/ ./
RUN rm -rf ./dist
RUN npm ci
RUN rm -rf .env
RUN mv .env.prod .env
RUN npm run build 
RUN ls -la ./dist

FROM golang:1.22 as builder

WORKDIR /server

ARG CGO_ENABLED=0

COPY ./backend/go.mod ./go.mod
COPY ./backend/go.sum ./go.sum
RUN go mod download

COPY ./backend .

RUN mkdir -p /web/dist
COPY --from=frontend /usr/src/app/dist ./web/dist

RUN go build -o binary ./cmd/

CMD ["./binary"]
