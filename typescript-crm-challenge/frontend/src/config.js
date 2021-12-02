import * as envSchema from 'env-schema';

const config =  {
  VERSION: {
    type: 'string',
    default: '0.1.0',
  },
  BACKEND_URL: {
    type: 'string',
    default: 'http://0.0.0.0',
  },
  BACKEND_PORT: {
    type: 'number',
    default: 3001
  },
};

export const schema = {
  type: 'object',
  required: Object.keys(config), // all properties are required!
  properties: config,
};

export default envSchema({ schema });
