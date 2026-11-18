# Stage 1: Build aset Vue
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# Stage 2: Serve dengan Nginx
FROM nginx:alpine

# Copy konfigurasi Nginx kustom
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy build result dari Vite
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]