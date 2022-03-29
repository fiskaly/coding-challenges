COPY customers(customer_id, first_name, last_name, mail, tss_id)
FROM '/postgres/customers.csv'
DELIMITER ','
CSV HEADER;