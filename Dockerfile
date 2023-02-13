FROM node:16-alpine AS builder
WORKDIR /app
RUN npm i -g pnpm
COPY package.json pnpm-lock.yaml tsconfig.json tsconfig.build.json ./
RUN pnpm install --frozen-lockfile
COPY src/ ./src/
COPY types/ ./types/
RUN pnpm run build

FROM nixos/nix:latest as goBuilder
WORKDIR /app
COPY flake.nix flake.lock .
COPY build ./build/
RUN echo "experimental-features = nix-command flakes" >> /etc/nix/nix.conf
RUN nix-env -iA nixpkgs.gnumake
RUN nix develop .\#build
COPY . .
RUN make go-build

FROM node:16-alpine
ARG NODE_ENV=production
ENV NODE_ENV $NODE_ENV
WORKDIR /app
RUN npm i -g pnpm
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile --prod  && pnpm store prune
COPY entrypoint.sh .
COPY migrations/ ./migrations/
COPY email-templates/ ./email-templates
COPY --from=builder ./app/dist dist/
COPY --from=goBuilder /app/result/bin/hasura-auth /usr/local/bin/go-hasura-auth
HEALTHCHECK --interval=60s --timeout=2s --retries=3 CMD wget http://localhost:${AUTH_PORT}/healthz -q -O - > /dev/null 2>&1
ENTRYPOINT ["/app/entrypoint.sh"]
