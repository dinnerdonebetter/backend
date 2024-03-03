FROM node:lts-bullseye

WORKDIR /src/github.com/dinnerdonebetter/frontend

COPY . .

RUN yarn install

ENTRYPOINT [ "yarn", "run", "dev" ]
