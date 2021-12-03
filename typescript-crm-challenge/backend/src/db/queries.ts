import queryDb from './postgres';


export function getCustomer(customer_id: string): Promise<string[]> {
  const queryString: string = `
  SELECT customer_id, first_name, last_name, mail
  FROM customers
  WHERE customer_id = '${customer_id}'
  `; // note: basic string formatting is NOT save against SQL-injections
  return queryDb(queryString);
}