# Snippy Frontend - React + Vite + Tailwind CSS

This is the standalone React frontend project for the Go URL Shortener service.

## 🚀 Setup & Local Development

### Prerequisites
Make sure you have Node.js (version 18+ recommended) and `npm` installed on your system.

### Installation
1. Navigate into the frontend directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   npm install
   ```

### Running the Application
1. Run the local development server:
   ```bash
   npm run dev
   ```
2. The frontend will be served at `http://localhost:3000`.
3. The Vite server is configured with a built-in proxy. Any API requests to `/api/v1` are automatically proxied to the Go backend running on port `8080`, eliminating CORS issues during development!

## 📦 Production Build
To build the application for deployment (generating optimized static assets under `frontend/dist`):
```bash
npm run build
```
