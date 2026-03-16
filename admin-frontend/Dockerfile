# ── Stage 1: Build ────────────────────────────────────────────────────────────
FROM ghcr.io/cirruslabs/flutter:stable AS builder

WORKDIR /app

# Cache pub dependencies separately from the source code.
COPY pubspec.yaml pubspec.lock ./
RUN flutter pub get

# Copy the rest of the source and build for the web.
COPY . .
RUN flutter build web --release --no-tree-shake-icons

# ── Stage 2: Serve ────────────────────────────────────────────────────────────
FROM nginx:alpine AS runner

# Remove the default nginx static content.
RUN rm -rf /usr/share/nginx/html/*

# Copy the compiled Flutter web output.
COPY --from=builder /app/build/web /usr/share/nginx/html

# Custom nginx config: serve the SPA correctly (all paths → index.html).
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
