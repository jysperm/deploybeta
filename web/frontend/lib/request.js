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
    if (res.ok) {
      return res.json();
    } else {
      return res.text().then( body => {
        var error;

        try {
          error = JSON.parse(body).error;
        } finally {
          const err = new Error(`${res.status}: ${error || res.statusText}`)
          throw _.extend(err, {res, body});
        }
      });
    }
  });
}