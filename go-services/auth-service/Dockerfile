FROM golang:1.23.2-bookworm

WORKDIR /app 

COPY . ./

RUN go build -o ./auth-service

EXPOSE 5000

CMD [ "./auth-service" ]
