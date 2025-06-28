FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Copier les fichiers de configuration Go et les dépendances
COPY go.mod go.sum ./
RUN go mod download

# Copier tout le projet dans l'image
COPY . .

# Générer le Prisma Client Go
RUN go run github.com/steebchen/prisma-client-go generate --schema internal/database/prisma/schema.prisma

ARG VERSION
ENV VERSION=${VERSION}

RUN echo $VERSION >> VERSION

# Compiler l'application Go (optimisation avec -ldflags pour la taille de l'image)
RUN go build -ldflags "-s -w" -o app cmd/main.go

FROM golang:1.24.2-alpine

LABEL org.opencontainers.image.source=https://github.com/Q300Z/CESIZen_API


RUN apk --no-cache add curl

WORKDIR /app

# Copier l'application compilée depuis l'étape builder
COPY --from=builder /app/app .
COPY --from=builder /app/scripts/entrypoint.sh .
COPY --from=builder /app/internal/database/prisma/ ./internal/database/prisma/
COPY --from=builder /app/go.mod /app/go.sum ./
COPY --from=builder /app/.air.toml ./

# Installer les dépendances nécessaires pour Prisma
RUN go get github.com/steebchen/prisma-client-go

# On rend executable entrypoint.sh
RUN chmod +x entrypoint.sh



# Configurer l'utilisateur pour éviter les problèmes de permission
RUN addgroup -S -g 1000 cesizen && adduser -S -u 1000 -G cesizen cesizen


# On configure les permissions pour l'utilisateur non-root
RUN chown -R cesizen:cesizen /app

# Passer à l'utilisateur non-root
USER cesizen

# Arguments de build
ARG GIN_MODE
ARG VERSION

# Variables d’environnement
ENV GIN_MODE=${GIN_MODE:-release}
ENV VERSION=${VERSION}

EXPOSE 8080

CMD ["./entrypoint.sh"]
