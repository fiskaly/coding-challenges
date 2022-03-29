import React from 'react';
import config from './config';
import axios from 'axios';


// example customer id: 171c0f84-0b77-4cfc-96b1-368ddba2eb52

function fetchCustomer(customer_id) {
  const url = `${config.BACKEND_URL}:${config.BACKEND_PORT}/customer`;
  axios.post(url, {customer_id: customer_id})
    .then((res) => {
      console.log(res.data);
    });
}


export default function CustomerTable(props) {
  // fetchCustomer('171c0f84-0b77-4cfc-96b1-368ddba2eb52');
  return <p>Replace this with a table displaying the customers</p>
}