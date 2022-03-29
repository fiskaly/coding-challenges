import server from './api/server';

const port = 3000;
server.listen(port, () => {
  console.log(`Running signature service on port ${port}`)
});


