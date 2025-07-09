# CESIZen API

API pour l'application CESIZen.

## Déploiement (Production & Développement distant)

Ce guide explique comment déployer l'application en utilisant Docker et Traefik comme reverse proxy.

### 1. Prérequis

Avant de commencer, assurez-vous d'avoir les éléments suivants :

*   Un serveur distant (VPS, dédié, etc.) avec un accès `sudo` ou `root`.
*   **Docker** et **Docker Compose** installés sur le serveur. Vous pouvez suivre [la documentation officielle](https://docs.docker.com/engine/install/) pour l'installation.
*   **Noms de domaine** configurés pour pointer vers l'adresse IP de votre serveur. Pour ce projet, vous aurez besoin de :
    *   **Pour la production :**
        *   `cesizen.qalpuch.cc` (pour le front-end)
        *   `cesizen-api.qalpuch.cc` (pour l'API)
    *   **Pour le développement distant :**
        *   `cesizen-dev.qalpuch.cc` (pour le front-end)
        *   `cesizen-api-dev.qalpuch.cc` (pour l'API)
*   Un **compte Cloudflare** pour gérer vos domaines et un **token d'API DNS**. Traefik l'utilisera pour générer automatiquement les certificats SSL/TLS avec Let's Encrypt.

### 2. Configuration

#### a. Cloner le dépôt

Connectez-vous à votre serveur et clonez ce dépôt :

```bash
git clone <URL_DU_DEPOT>
cd CESIZen_API
```

#### b. Fichiers d'environnement

Vous devez créer trois fichiers d'environnement distincts à la racine du projet.

**1. Fichier `.env` (pour Traefik)**

Ce fichier contient les informations pour que Traefik puisse communiquer avec l'API de Cloudflare.

```ini
# .env
CF_API_EMAIL=votre-email@cloudflare.com
CF_DNS_API_TOKEN=votre_token_api_cloudflare
```

**2. Fichier `.env.prod` (pour la production)**

Ce fichier contient les variables pour la stack de production.

```ini
# .env.prod
DB_USER=user_prod
DB_PASSWORD=un_mot_de_passe_solide
DB_HOST=cz-db-prod
DB_PORT=5432
DB_NAME=cesizen
DATABASE_URL="postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?schema=public"
APP_HOST=0.0.0.0
APP_PORT=8080
JWT_SECRET=un_secret_jwt_tres_long_et_aleatoire_pour_la_prod
GIN_MODE=release
VERSION=latest
```

**3. Fichier `.env.dev` (pour le développement distant)**

Ce fichier est similaire au fichier de production, mais pour l'environnement de développement.

```ini
# .env.dev
DB_USER=user_dev
DB_PASSWORD=un_autre_mot_de_passe_solide
DB_HOST=cz-db-dev
DB_PORT=5432
DB_NAME=cesizen
DATABASE_URL="postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?schema=public"
APP_HOST=0.0.0.0
APP_PORT=8080
JWT_SECRET=un_secret_jwt_different_pour_le_dev
GIN_MODE=release
VERSION=dev
```

### 3. Lancement des services

Les commandes doivent être lancées depuis la racine du projet sur votre serveur.

#### a. Lancer les services globaux (Traefik & Watchtower)

Cette commande n'a besoin d'être lancée qu'une seule fois. Traefik gérera le routage et Watchtower mettra à jour vos applications automatiquement.

```bash
docker-compose -f docker-compose.global.yml up -d
```

#### b. Lancer la stack de Production

Pour déployer ou mettre à jour l'environnement de production :

```bash
docker-compose -f docker-compose.prod.yml up -d
```

#### c. Lancer la stack de Développement distant

Pour déployer ou mettre à jour l'environnement de développement :

```bash
docker-compose -f docker-compose.dev.yml -d
```

Vos applications sont maintenant déployées et accessibles via les noms de domaine que vous avez configurés.

---

## Développement Local

Pour le développement sur votre machine locale, suivez ces étapes :

1.  **Créez un fichier `.env`** à la racine du projet en vous basant sur `.env.example`.
2.  Assurez-vous que Docker est en cours d'exécution.
3.  Lancez la stack avec la commande suivante :

```bash
docker-compose up --build
```

L'application sera accessible sur :
*   API : `http://localhost:8080`
*   Front-end : `http://localhost:3000`
