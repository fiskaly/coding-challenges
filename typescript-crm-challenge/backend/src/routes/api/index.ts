import { FastifyInstance } from 'fastify';

declare module 'fastify' {
    interface FastifyInstance {
    }
}

export default async function(fastify: FastifyInstance) {

};
