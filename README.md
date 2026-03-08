# HobbyManager

Application pour gérer la liste des achats culturels à faire : livres, films, séries, jeux de société, jeux vidéo.

- **Frontend** : Vue 3 (Vite)
- **Backend** : Go (API REST, SQLite)
- **API** : contrat OpenAPI dans [api/openapi.yaml](api/openapi.yaml)

## Structure

- `api/` – Spécification OpenAPI (source de vérité)
- `backend/` – Serveur Go (CRUD œuvres)
- `frontend/` – Application Vue 3
- `tests/` – Tests E2E et intégration / contrat

Voir [AGENTS.md](AGENTS.md) pour le périmètre des agents (frontend, backend, tests).

## Démarrage

1. **Backend** (terminal 1)  
   ```bash
   cd backend && go mod tidy && go run ./cmd/server
   ```  
   API : http://localhost:8080

2. **Frontend** (terminal 2)  
   ```bash
   cd frontend && npm install && npm run dev
   ```  
   App : http://localhost:5173 (proxy `/api` vers le backend)

3. Base de données : par défaut `backend/hobby.db` (créée automatiquement). Variable d'environnement `DB_PATH` pour changer le chemin.

## Tests

- Backend : `cd backend && go test ./...`
- Frontend : `cd frontend && npm run test:run`

## Licence

Apache-2.0
