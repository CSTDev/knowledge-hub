FROM node:8
WORKDIR /usr/src/app
COPY ./package.json ./
RUN npm install
COPY ./src/ ./src
COPY ./public/ ./public
RUN npm run build

FROM node:8
ENV NPM_CONFIG_PREFIX=/home/node/.npm-global
USER node
WORKDIR /usr/src/app/
RUN npm install -g @staticdeploy/app-server
COPY --from=0 /usr/src/app/build ./build
COPY --from=0 /usr/src/app/node_modules ./node_modules
ENV PATH="{$PATH}:/home/node/.npm-global/bin/"

ENTRYPOINT app-server

EXPOSE 3000

