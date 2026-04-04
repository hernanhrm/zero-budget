# syntax=docker/dockerfile:1
# Build from repo: docker build -f frontend/Dockerfile frontend --build-arg VITE_IDENTITY_URL=https://...
FROM node:22-alpine AS build
WORKDIR /app

RUN corepack enable

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY apps ./apps
COPY packages ./packages

RUN pnpm install --frozen-lockfile

ARG VITE_IDENTITY_URL=""
ENV VITE_IDENTITY_URL=${VITE_IDENTITY_URL}

RUN pnpm --filter app build

FROM nginx:alpine AS runner
COPY --from=build /app/apps/app/dist /usr/share/nginx/html
COPY nginx/default.conf /etc/nginx/conf.d/default.conf
EXPOSE 8080
