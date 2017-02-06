import React, {Component} from 'react';
import ReactDOM from 'react-dom';
import {hashHistory, Router, IndexRoute, Route} from 'react-router';
import {Row, Button} from 'react-bootstrap';

import LayoutView from './layout'
import LoginView from './login'

class RootView extends Component {
  render() {
    return <Router history={hashHistory}>
      <Router path='/'>
        <IndexRoute component={IndexView} />
        <Route path='/login' component={LoginView} />
      </Router>
    </Router>;
  }
}

class IndexView extends Component {
  render() {
    return <LayoutView>
      <Row>
        <Button href='#/login' bsStyle='primary' bsSize='large'>Login or Sign up</Button>
      </Row>
    </LayoutView>;
  }
}

ReactDOM.render(React.createElement(RootView), document.getElementById('react-root'));
