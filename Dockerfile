# Etapa de construcción
FROM golang:1.18 AS builder

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos de Go necesarios para construir el proyecto
COPY go.mod .
COPY go.sum .

# Descarga las dependencias de Go
RUN go mod download

# Copia el código fuente del proyecto y los assets
COPY . .

# Compila la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -a -installsuffix cgo -o zapping_stream ./cmd

# Etapa de ejecución
FROM alpine:latest  

# Instala las herramientas necesarias, si hay alguna
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copia los assets y el binario compilado desde la etapa de construcción
COPY --from=builder /app/zapping_stream .
COPY --from=builder /app/assets ./assets  

# Crea un archivo .env con las variables de entorno en la ubicación correcta
RUN echo "DB_HOST=localhost" >> .env && \
    echo "DB_PORT=5432" >> .env && \
    echo "DB_USER=postgres" >> .env && \
    echo "DB_PASSWORD=penti" >> .env && \
    echo "DB_NAME=zapping" >> .env && \
    echo "DB_SSLMODE=disable" >> .env && \
    echo "SALT=PP12312SSADl.^^&212321(((**(!@3!@LKJ*(ASD77}}" >> .env && \
    echo "JWT_SECRET_KEY=PP1123l.435341(((**(zxc!@LKJa21D77}}" >> .env

# Expone el puerto que tu aplicación utiliza
EXPOSE 8080

# Ejecuta la aplicación de Go
CMD ["./zapping_stream"]
