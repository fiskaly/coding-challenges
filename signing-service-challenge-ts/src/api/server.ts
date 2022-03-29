import express from 'express';
import bodyParser from 'body-parser';

const server = express();

server.use(bodyParser.json());

server.get('/health', (req, res) => {
  res.status(200);
  res.send(JSON.stringify({
    status: 'pass',
    version: 'v1'
  }));
});

// TODO: REST endpoints ...

export default server;
