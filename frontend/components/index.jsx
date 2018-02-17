import {HashRouter as Router, Switch, Route} from 'react-router-dom';
import {Row, Button} from 'react-bootstrap';
import React, {Component} from 'react';
import ReactDOM from 'react-dom';

import LayoutView from './layout';
import LoginView from './login';
import ConsoleView from './console';

class RootView extends Component {
  render() {
    return <Router>
      <Switch>
        <Route exact path='/' component={IndexView} />
        <Route path='/login' component={LoginView} />
        <Route path='/console' component={ConsoleView} />
      </Switch>
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
