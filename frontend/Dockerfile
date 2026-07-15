# Multi-stage Dockerfile for React Vite Frontend

# Base node image for dependency installation and development
FROM node:20-alpine AS base
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .

# Development stage: runs the Vite dev server
FROM base AS development
EXPOSE 5173
CMD ["npm", "run", "dev", "--", "--host"]

# Build stage: compiles typescript and builds production assets
FROM base AS builder
RUN npm run build

# Production stage: serves the built assets using Nginx
FROM nginx:alpine AS production
COPY --from=builder /app/dist /usr/share/nginx/html
# Copy default nginx config if needed, otherwise use default nginx config
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
