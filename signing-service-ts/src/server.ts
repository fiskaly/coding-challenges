import express from 'express';

const server = express();
const port = 3000;

server.listen(port, () => {
  console.log(`Running signature service on port ${port}`)
});

export default server;
