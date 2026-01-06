# Multi-stage build for production-ready frontend
# Stage 1: Build
FROM node:20.18.2-alpine AS builder

# Set working directory
WORKDIR /app

# Install security updates and dependencies
RUN apk add --no-cache \
    libc6-compat \
    && rm -rf /var/cache/apk/*

# Copy package files
COPY package.json pnpm-lock.yaml ./

# Install pnpm
RUN npm install -g pnpm@9.15.0

# Install dependencies with frozen lockfile for reproducible builds
RUN pnpm install --frozen-lockfile --prod=false

# Declare build arguments
ARG VITE_API_URL
ENV VITE_API_URL=${VITE_API_URL}

# Copy application source code
COPY . .

# Build the application
RUN pnpm build

# Stage 2: Production
FROM nginx:1.27.3-alpine

# Install security updates
RUN apk add --no-cache \
    curl \
    && rm -rf /var/cache/apk/* \
    /var/tmp/* \
    /tmp/*

# Create non-root user
RUN addgroup -g 1001 -S nodejs && \
    adduser -S -D -H -u 1001 -h /usr/share/nginx/html -s /sbin/nologin -G nodejs nodejs

# Create nginx runtime directories with proper permissions
RUN mkdir -p /var/cache/nginx /var/log/nginx /run && \
    chown -R nodejs:nodejs /var/cache/nginx /var/log/nginx /run

# Copy custom nginx config and entrypoint
COPY docker/nginx.conf /etc/nginx/nginx.conf
COPY docker/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

# Copy built files from builder stage
COPY --from=builder --chown=nodejs:nodejs /app/dist /usr/share/nginx/html

# Set proper permissions
RUN chown -R nodejs:nodejs /usr/share/nginx/html && \
    chmod -R 755 /usr/share/nginx/html && \
    find /usr/share/nginx/html -type d -exec chmod 755 {} \; && \
    find /usr/share/nginx/html -type f -exec chmod 644 {} \;

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/ || exit 1

# Use entrypoint script
ENTRYPOINT ["/docker-entrypoint.sh"]
