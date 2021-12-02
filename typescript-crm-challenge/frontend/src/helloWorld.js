import React from 'react';
import axios from 'axios';
import config from './config';


export default function HelloWorld(props) {
  let text = 'loading...';

  const url = `${config.BACKEND_URL}:${config.BACKEND_PORT}/hello`;
  axios.get(url)
    .then((res) => {
      console.log(res.data);
      if (res.status !== 200) {
        text = `${res.status} error when calling the backend`;
      } else {
        text = res.data;
      }
    });

  return <p>
    {text}
  </p>
}