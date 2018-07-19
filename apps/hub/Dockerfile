FROM node:8
WORKDIR /usr/src/app
COPY ./package.json ./
RUN npm install
COPY ./src/ ./src
COPY ./public/ ./public
RUN npm run build

FROM node:8
WORKDIR /usr/src/app/
RUN npm install -g serve
COPY --from=0 /usr/src/app/build ./build
COPY --from=0 /usr/src/app/node_modules ./node_modules

ENTRYPOINT serve -s build

EXPOSE 5000

