import 'whatwg-fetch';
import _ from 'lodash';

export function requestJson(url, options = {}) {
  options.headers = options.headers || {};

  if (options.body) {
    options.headers['Content-Type'] = 'application/json';
    options.body = JSON.stringify(options.body);
  }

  const sessionToken = localStorage.getItem('sessionToken');

  if (sessionToken) {
    options.headers['Authorization'] = sessionToken;
  }

  return fetch(url, options).then( res => {
    return res.text().then( body => {
      var err;

      if (res.ok) {
        try {
          return JSON.parse(body);
        } catch (err) {
          return body;
        }
      } else {
        try {
          err = new Error(`${res.status}: ${JSON.parse(body).error}`);
        } catch (e) {
          err = new Error(`${res.status}: ${res.statusText}`);
        }

        throw _.extend(err, {res, body});
      }
    });
  });
}
