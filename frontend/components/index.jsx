import React, {Component} from 'react';
import ReactDOM from 'react-dom';
import {hashHistory, Router, IndexRoute, Route} from 'react-router';
import {Row, Button} from 'react-bootstrap';

import LayoutView from './layout'
import LoginView from './login'
import ConsoleView from './console'

class RootView extends Component {
  render() {
    return <Router history={hashHistory}>
      <Router path='/'>
        <IndexRoute component={IndexView} />
        <Route path='/login' component={LoginView} />
        <Route path='/console' component={ConsoleView} />
      </Router>
    </Router>;
  }
}

class IndexView extends Component {
  render() {
    const username = localStorage.getItem('username');
    var button;

    if (username) {
      button = <Button href='#/console' bsStyle='primary' bsSize='large'>Back to console as {username}</Button>;
    } else {
      button = <Button href='#/login' bsStyle='primary' bsSize='large'>Login or Sign up</Button>;
    }

    return <LayoutView pageHeader={true}>
      <Row>{button}</Row>
    </LayoutView>;
  }
}

ReactDOM.render(React.createElement(RootView), document.getElementById('react-root'));
