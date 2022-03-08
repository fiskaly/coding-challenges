import server from '../server';

server.get('/health', (req, res) => {
  res.status(200);
  res.send(JSON.stringify({
    status: 'pass',
    version: 'v1'
  }));
});

// TODO: register services here



