// PM2 process definitions for Prep Dashboard.
// Backend (Go binary) on :9002, Frontend (vite preview of built dist) on :9001.
// Deploy with ./deploy.sh — it builds both then (re)starts these.
const path = require('path')

const root = __dirname

module.exports = {
  apps: [
    {
      name: 'prep-backend',
      // Compiled by deploy.sh; cwd is backend/ so .env + godotenv resolve.
      script: path.join(root, 'backend', 'prep-backend'),
      cwd: path.join(root, 'backend'),
      env: {
        PORT: '9002',
        DB_NAME: 'prep_dashboard',
        // MONGO_URI is intentionally NOT set here: godotenv.Load() won't
        // override env vars already present in the process, so hardcoding it
        // would shadow backend/.env (which holds the real Atlas URI/secret).
        // Allow the deployed frontend origin through CORS.
        CORS_ORIGIN: 'http://localhost:9001',
      },
      autorestart: true,
      max_restarts: 10,
      out_file: path.join(root, 'logs', 'backend-out.log'),
      error_file: path.join(root, 'logs', 'backend-err.log'),
    },
    {
      name: 'prep-frontend',
      // `npm run preview` serves the built dist/ (see vite.config.js preview block).
      script: 'npm',
      args: 'run preview',
      cwd: path.join(root, 'frontend'),
      interpreter: 'none', // let PM2 exec npm directly
      env: {
        FRONTEND_PORT: '9001',
        BACKEND_URL: 'http://localhost:9002',
        // Hostnames nginx may forward; vite preview rejects anything else.
        ALLOWED_HOSTS: 'prep.aryan-tech.xyz,localhost',
      },
      autorestart: true,
      max_restarts: 10,
      out_file: path.join(root, 'logs', 'frontend-out.log'),
      error_file: path.join(root, 'logs', 'frontend-err.log'),
    },
  ],
}
