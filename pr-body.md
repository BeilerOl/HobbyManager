Implémentation initiale selon le plan HobbyManager.

## Contenu
- **api/** : Spécification OpenAPI 3.0 avec CRUD sur `/api/v1/works`
- **backend/** : Serveur Go, repository SQLite, handlers HTTP, CORS, tests unitaires et intégration
- **frontend/** : Vue 3 + Vite, liste/détail/formulaire, filtres type et vu, client API
- **tests/** : Structure pour e2e et tests d'intégration/contrat
- **AGENTS.md** : Périmètres des agents et référence au contrat API

## Démarrage
1. Backend: `cd backend && go mod tidy && go run ./cmd/server`
2. Frontend: `cd frontend && npm install && npm run dev`
