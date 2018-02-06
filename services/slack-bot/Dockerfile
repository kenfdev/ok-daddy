FROM node:8.9.1-alpine

RUN addgroup -S app && adduser -S -g app app

WORKDIR /root/

# Turn down the verbosity to default level.
ENV NPM_CONFIG_LOGLEVEL warn

RUN mkdir -p /home/app

# Wrapper/boot-strapper
WORKDIR /home/app
COPY . ./

# This ordering means the npm installation is cached for the outer function handler.
RUN npm i

# Set correct permissions to use non root user
WORKDIR /home/app/

# chmod for tmp is for a buildkit issue (@alexellis)
RUN chown app:app -R /home/app \
  && chmod 777 /tmp

USER app

CMD ["node", "main.js"]

