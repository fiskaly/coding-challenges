import 'make-promises-safe';
import fastify from 'fastify';
import fastifyCors from 'fastify-cors';
import app from './app';

export const server = fastify({
    logger: true,
    pluginTimeout: 10000,
    trustProxy: true,
});

server.register(fastifyCors, {
    credentials: true,
    origin: 'http://localhost:3000'  // only allow frontend to make calls
})

server.register(app);

const opts = {
    port: Number(process.env.PORT) || 3001,
    host: process.env.HOST || 'localhost',
};
server.listen(opts, (err) => {
    if (err) {
        server.log.error(err);
        process.exit(1);
    }
});
