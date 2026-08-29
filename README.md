# 🎯 Prep Dashboard

A placement-prep dashboard to track your daily study activity, keep notes, and
browse curated study materials for the topics that matter for cracking good
companies:

- **DSA** · **LLD** · **Projects** · **OS Fundamentals** · **DBMS + SQL** · **HR Questions**

Plus a **link-preview** feature: paste any DSA question URL (LeetCode / GFG /
Codeforces) and get a rich card with title, description, and thumbnail.

## Stack

| Layer     | Tech                     |
| --------- | ------------------------ |
| Frontend  | React (Vite)             |
| Backend   | Go + Gin                 |
| Database  | MongoDB                  |

Monorepo layout:

```
prep-dashboard/
├── backend/      # Go + Gin API
├── frontend/     # React + Vite app
└── docker-compose.yml   # local MongoDB
```

## Features

- 📊 **Overview** — total activities, completed, hours logged today, hours-by-topic bars
- 📝 **Activity Log** — log daily study with topic, hours, notes; check off as done
- 🗒️ **Notes** — topic-tagged notes with tags, edit/delete
- 📚 **Study Materials** — curated roadmaps + resource links per topic
- 🔗 **DSA Links** — save question URLs and get scraped rich previews

## Getting started

### 1. Start MongoDB

With Docker:

```bash
docker compose up -d
```

Or use a [MongoDB Atlas](https://www.mongodb.com/atlas) free cluster and set
`MONGO_URI` accordingly (see below). Or install MongoDB locally.

### 2. Backend

```bash
cd backend
cp .env.example .env      # adjust MONGO_URI if needed
go run .
```

API runs on **http://localhost:9002**.

Environment variables (`backend/.env`):

| Var         | Default                     | Notes                        |
| ----------- | --------------------------- | ---------------------------- |
| MONGO_URI   | mongodb://localhost:27017   | Atlas SRV string also works  |
| DB_NAME     | prep_dashboard              |                              |
| PORT        | 9002                        |                              |
| CORS_ORIGIN | http://localhost:9001       | frontend origin              |

### 3. Frontend

```bash
cd frontend
npm install
npm run dev
```

Dev app runs on **http://localhost:5173** and proxies `/api` to the backend
(`BACKEND_URL`, default `http://localhost:9002`).

## API reference

| Method | Path                  | Purpose                          |
| ------ | --------------------- | -------------------------------- |
| GET    | /api/health           | health check                     |
| GET    | /api/guides           | curated study materials          |
| GET    | /api/activities       | list activities (`?topic`,`?date`)|
| POST   | /api/activities       | create activity                  |
| PATCH  | /api/activities/:id   | update (toggle done, edit)       |
| DELETE | /api/activities/:id   | delete activity                  |
| GET    | /api/stats            | dashboard aggregates             |
| GET    | /api/notes            | list notes (`?topic`)            |
| POST   | /api/notes            | create note                      |
| PATCH  | /api/notes/:id        | update note                      |
| DELETE | /api/notes/:id        | delete note                      |
| POST   | /api/link-preview     | scrape OG/meta for a URL         |

## Deploy (PM2)

Runs the backend on **:9002** and the frontend on **:9001**.

```bash
npm install -g pm2      # once
./deploy.sh             # builds both, then starts/reloads under PM2
```

`deploy.sh` builds the Go binary (`backend/prep-backend`) and the Vite
production build (`frontend/dist`), then `pm2 startOrReload ecosystem.config.cjs`
(zero-downtime reload if already running). The frontend is served by
`vite preview`, which proxies `/api` → `http://localhost:9002`.

Useful PM2 commands:

```bash
pm2 status            # process list
pm2 logs prep-backend # tail logs (also in ./logs/)
pm2 restart all
pm2 startup           # once, to auto-start on boot
```

Ports/URIs (e.g. `MONGO_URI`) live in `ecosystem.config.cjs`.

## Tests

```bash
cd backend && go test ./...
```
