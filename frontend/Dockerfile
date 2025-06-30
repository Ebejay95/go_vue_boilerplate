# Multi-stage build for Vue.js application
FROM node:18-alpine AS build-stage

WORKDIR /app

# Copy package files first
COPY frontend/package*.json ./

# Install all dependencies (including devDependencies for build)
RUN npm install

# Copy proto files
COPY proto ../proto

# Copy ALL frontend source code (this will overwrite package.json but that's ok after npm install)
COPY frontend/ ./

# Debug: List files to verify structure
RUN echo "=== Frontend files ===" && ls -la
RUN echo "=== Source files ===" && ls -la src/

# Build Vue.js application
RUN npm run build:vue

# Production stage
FROM node:18-alpine AS production-stage

WORKDIR /app

# Copy package files first
COPY frontend/package*.json ./

# Install only production dependencies
RUN npm install --omit=dev

# Copy built Vue files and server
COPY --from=build-stage /app/dist ./dist
COPY --from=build-stage /app/server.js ./

# Copy proto files
COPY proto ../proto

# Set production environment
ENV NODE_ENV=production

EXPOSE 3000

CMD ["npm", "start"]
