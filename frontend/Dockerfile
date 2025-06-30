FROM node:18-alpine

WORKDIR /app

# Copy package files
COPY frontend/package*.json ./
RUN npm install

# Copy proto files
COPY proto ../proto

# Copy frontend source code
COPY frontend/ .

# Build the application
RUN npm run build

EXPOSE 3000

CMD ["npm", "start"]
