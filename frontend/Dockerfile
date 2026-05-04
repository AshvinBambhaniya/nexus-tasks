FROM node:20-alpine

WORKDIR /app

# Copy package files and install dependencies
COPY package*.json ./
RUN npm install

# Install system dependencies exactly as in the open-source project
RUN apk add --no-cache nginx supervisor wget curl
RUN mkdir -p /run/nginx

# Copy configuration files
COPY ./docker/nginx.conf /etc/nginx/nginx.conf
COPY ./docker/supervisord.conf /etc/supervisord.conf

# Copy application source and build
COPY . ./
RUN npm run build

# Match the port mapping from the project
EXPOSE 80

# Match the command style from the open-source project
ENTRYPOINT ["/bin/sh", "-c"]
CMD ["supervisord -c /etc/supervisord.conf"]
