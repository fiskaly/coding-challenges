import { FastifyRequest, FastifyReply, FastifyInstance } from 'fastify';
import { getCustomer } from '../../db/queries';


export default async function customer(fastify: FastifyInstance){
  fastify.route({
    method: 'POST',
    url: '/customer',
    schema: {
      response: {
        200: {
          type: 'array',
          items: {
            customer_id: {
              type: 'string'
            },
            first_name: {
              type: 'string'
            },
            last_name: {
              type: 'string'
            },
            mail: {
              type: 'string'
            }
          }
        }
      },
      body: {
        type: 'object',
        properties: {
          customer: { type: 'string' }
        },
        required: ['customer_id']
      }
    },
    // this function is executed for every request before the handler is executed
    preHandler: (request: FastifyRequest, reply: FastifyReply, done) => {
      // E.g. check authentication
      done();
    },
    handler: async (request: FastifyRequest, reply: FastifyReply) => {
      // @ts-ignore
      const customerId: string = request.body['customer_id'];
      console.log(request.body);
      const customerResult: string[] = await getCustomer(customerId);

      reply.send(customerResult);
    }
  });
}