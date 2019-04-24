#########################################
# Build stage
#########################################
FROM node:lts-alpine AS builder

# Setup env
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
ENV PATH /usr/src/app/node_modules/.bin:$PATH
ENV REACT_APP_API_ROOT https://api.readflow.app

# Install dependencies
COPY package.json /usr/src/app/package.json
COPY package-lock.json /usr/src/app/package-lock.json
RUN npm install --silent

# Build website
COPY . /usr/src/app
RUN npm run build

#########################################
# Distribution stage
#########################################
FROM nginx:stable-alpine

# Install website
COPY --from=builder /usr/src/app/build /usr/share/nginx/html

# Install configuration
RUN rm -rf /etc/nginx/conf.d
COPY etc/nginx/conf.d /etc/nginx/conf.d

# Port
EXPOSE 80

# Entrypoint
CMD ["nginx", "-g", "daemon off;"]
