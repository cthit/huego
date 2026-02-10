FROM node:latest AS build-step

RUN mkdir -p /usr/src/huego/frontend
RUN chown -R node /usr/src/huego/frontend

USER node

WORKDIR /usr/src/huego/frontend

COPY ./src ./src
COPY package.json .
COPY ./public ./public

RUN yarn install
RUN yarn global add react-scripts
RUN NODE_OPTIONS=--openssl-legacy-provider yarn build

FROM nginx:alpine
RUN rm /etc/nginx/conf.d/*
COPY nginx.conf /etc/nginx/conf.d/nginx.conf
COPY --from=build-step /usr/src/huego/frontend/build/ /usr/share/nginx/html/