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

## Cursor Cloud (agent Frontend)

Les agents Cloud qui travaillent sur le frontend utilisent la configuration suivante.

- **Install** (exécuté à la racine du dépôt) : `cd frontend && npm install`
- **Périmètre** : uniquement le dossier `frontend/`. Ne pas modifier `backend/` ni `api/openapi.yaml` sans alignement explicite du contrat.
- **Contrat API** : [api/openapi.yaml](api/openapi.yaml) — source de vérité pour les endpoints et schémas. Les appels frontend utilisent le préfixe `/api/v1` (proxy Vite vers le backend).
- **Commandes utiles** :
  - Lancer le frontend : `cd frontend && npm run dev`
  - Linter : `cd frontend && npm run lint` ou `npm run lint:fix`
  - Tests frontend : `cd frontend && npm run test:run`
  - Pour tester avec l’API réelle, lancer aussi le backend dans un autre terminal : `cd backend && go run ./cmd/server`
- **CI** : les pull requests déclenchent le workflow Frontend CI (`.github/workflows/frontend-ci.yml`) : lint + tests. S’assurer que lint et tests passent avant de merger.
- **GitHub Pages** : la version démo (frontend seul avec mock backend et données de test) est déployée automatiquement sur chaque push sur `main` qui touche `frontend/` (workflow `.github/workflows/deploy-pages.yml`). Build avec `VITE_USE_MOCK=true` et base path `/<nom-du-repo>/`. Activer Pages dans les réglages du dépôt (Source : GitHub Actions).
