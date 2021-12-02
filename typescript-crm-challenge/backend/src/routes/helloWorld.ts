import { FastifyRequest, FastifyReply, FastifyInstance } from 'fastify';


// example route
export default async function hello(fastify: FastifyInstance){
    fastify.route({
        method: 'GET',
        url: '/hello',
        schema: {
            response: {
                200: {
                    type: 'string',
                }
            }
        },
        // this function is executed for every request before the handler is executed
        preHandler: (request, reply, done) => {
            // E.g. check authentication
            done();
        },
        handler: (request, reply) => {
            reply.send('Hello World!');
        }
    });
}