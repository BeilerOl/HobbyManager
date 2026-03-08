# HobbyManager – Agents

Ce projet est développé de manière incrémentale avec des agents dédiés. Le contrat d’API est la source de vérité partagée.

## Contrat API

- **Fichier** : [api/openapi.yaml](api/openapi.yaml)
- Toute modification des endpoints ou des schémas doit être reflétée dans ce fichier et dans le backend puis le frontend.

## Périmètres

| Agent   | Dossier     | Rôle |
|---------|-------------|------|
| Backend | `backend/`  | Implémentation Go : handlers, repository (SQLite/PostgreSQL), modèles. S’aligner sur `api/openapi.yaml`. |
| Frontend| `frontend/` | Application Vue 3 : pages, composants, appels API. Consommer l’API selon `api/openapi.yaml`. |
| Tests   | `tests/` + tests dans `backend/` et `frontend/` | Tests unitaires et d’intégration dans chaque stack ; E2E et contrats dans `tests/`. |

## Lancer l’application

1. **Backend** : `cd backend && go run ./cmd/server` (API sur http://localhost:8080)
2. **Frontend** : `cd frontend && npm install && npm run dev` (app sur http://localhost:5173, proxy `/api` vers le backend)
