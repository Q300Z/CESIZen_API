FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Copier les fichiers de configuration Go et les dépendances
COPY go.mod go.sum ./
RUN go mod download

# Copier tout le projet dans l'image
COPY . .

# On supprime le Prisma Client déjà existant
RUN rm -R internal/database/prisma/db

# Générer le Prisma Client Go
RUN go run github.com/steebchen/prisma-client-go generate --schema internal/database/prisma/schema.prisma

ENV GIN_MODE=release

# Compiler l'application Go (optimisation avec -ldflags pour la taille de l'image)
RUN go build -ldflags "-s -w" -o app cmd/main.go

FROM golang:1.24.2-alpine

# Installer les certificats SSL nécessaires
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copier l'application compilée depuis l'étape builder
COPY --from=builder /app/app .
COPY --from=builder /app/scripts/entrypoint.sh .
COPY --from=builder /app/internal/database/prisma/ ./internal/database/prisma/
COPY --from=builder /app/go.mod /app/go.sum ./

# On rend executable entrypoint.sh
RUN chmod +x entrypoint.sh

# Copier le fichier .env pour la prod
COPY .env .


ENV DB_HOST=${DB_HOST}
ENV DB_PORT=${DB_PORT}

EXPOSE 8080

CMD ["./entrypoint.sh"]
