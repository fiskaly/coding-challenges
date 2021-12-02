import { FastifyInstance } from 'fastify';
import routes from './routes';


declare module 'fastify' {
    interface FastifyInstance {
        config: { [key: string]: any };
    }
}

export default async function app(fastify: FastifyInstance) {

    fastify.register(routes);

    if (process.env.NODE_ENV === 'development') {
        fastify.ready(() => {
            console.log(fastify.printRoutes());
        });
    }
}