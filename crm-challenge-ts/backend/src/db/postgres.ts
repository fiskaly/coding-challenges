import config from '../config';
import { Pool } from 'pg';

const v1: Pool = new Pool({
    user: <string>config.PSQL_USER,
    host: <string>config.PSQL_HOST,
    database: <string>config.PSQL_DB,
    password: <string>config.PSQL_PASSWORD,
    port: <number>config.PSQL_PORT,
  }
);

export default function queryDb(query: string):Promise<any[]> {
  return v1.query(query)
    .then((res) => {
      return res.rows;
    })
    .catch((err) => {
      console.log(err);
      console.log(query);
      throw(err);
    });
}
