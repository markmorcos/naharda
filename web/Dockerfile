# Multi-stage → distroless Node (project.md §6, §11). Astro standalone SSR server.
FROM node:24-bookworm-slim AS build
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm install
COPY . .
RUN npm run build

FROM gcr.io/distroless/nodejs24-debian12
WORKDIR /app
COPY --from=build /app/dist ./dist
COPY --from=build /app/node_modules ./node_modules
COPY --from=build /app/package.json ./package.json
ENV HOST=0.0.0.0
ENV PORT=8080
EXPOSE 8080
CMD ["dist/server/entry.mjs"]
