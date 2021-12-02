FROM node:14-alpine

ENV NODE_ENV=production

WORKDIR /app
COPY . /app/
RUN npm install --save-dev
RUN npm run build

EXPOSE 3001/tcp

CMD ["npm", "run", "start"]