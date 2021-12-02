import React, {useState} from 'react';
import axios from 'axios';
import config from './config';


export default function HelloWorld(props) {
  const [text, setText] = useState('loading...');

  const url = `${config.BACKEND_URL}:${config.BACKEND_PORT}/hello`;
  axios.get(url)
    .then((res) => {
      console.log(res.data);
      if (res.status !== 200) {
        setText(`${res.status} error when calling the backend`);
      } else {
        setText(res.data);
      }
    });

  return <p>
    {text}
  </p>
}