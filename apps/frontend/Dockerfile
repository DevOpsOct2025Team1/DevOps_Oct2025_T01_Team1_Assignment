# https://bun.com/docs/guides/ecosystem/docker
FROM oven/bun:1-alpine AS base
WORKDIR /usr/src/app

FROM base AS install
RUN mkdir -p /temp/dev
COPY apps/frontend/package.json apps/frontend/bun.lock /temp/dev/
RUN cd /temp/dev && bun ci

RUN mkdir -p /temp/prod
COPY apps/frontend/package.json apps/frontend/bun.lock /temp/prod/
RUN cd /temp/prod && bun ci --production

FROM base AS prerelease
COPY --from=install /temp/dev/node_modules node_modules
COPY apps/frontend/ .

ENV NODE_ENV=production
RUN bun run build

# wtf react router is not compatible with bun? that's so lame
FROM node:24-alpine AS release
WORKDIR /usr/src/app
COPY --from=install /temp/prod/node_modules node_modules
COPY --from=prerelease /usr/src/app/build ./build
COPY --from=prerelease /usr/src/app/package.json .

USER node
EXPOSE 3000/tcp
ENTRYPOINT ["npm", "run", "start"]