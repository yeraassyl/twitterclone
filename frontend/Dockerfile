FROM node:16.13.0-alpine as build

WORKDIR /usr/local/app

COPY ./ /usr/local/app/

RUN npm install

RUN npm run build

FROM nginx:latest

COPY --from=build /usr/local/app/dist/twitter-front /usr/share/nginx/html

EXPOSE 80
