import { FastifyInstance } from 'fastify';
import hello from './helloWorld';
import api from './api';

export default async function(fastify: FastifyInstance) {
    fastify.register(hello);
    fastify.register(api);
};
