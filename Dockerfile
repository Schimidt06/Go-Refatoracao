FROM golang:1.23-alpine

WORKDIR /app

# Copia os arquivos do módulo e resolve as dependências
COPY go.mod go.sum ./
RUN go mod tidy

# Copia todo o código do projeto
COPY . .

# Compila a aplicação
RUN CGO_ENABLED=0 go build -o main .

EXPOSE 8080

CMD ["./main"]
