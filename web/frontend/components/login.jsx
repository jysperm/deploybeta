import React, {Component} from 'react';
import {Row, Col, Tabs, Tab, Form, FormGroup, FormControl, ControlLabel, Button} from 'react-bootstrap';
import _ from 'lodash';

import {requestJson} from '../lib/request'
import LayoutView from './layout'

export default class LoginView extends Component {
  render() {
    return <LayoutView>
      <Row>
        <Tabs defaultActiveKey={2} id='signup-or-login'>
          <Tab eventKey={1} title='Sign up'>
            <SignupTab />
          </Tab>
          <Tab eventKey={2} title='Login'>
            <LoginTab />
          </Tab>
        </Tabs>
      </Row>
    </LayoutView>;
  }
}

class SignupTab extends Component {
  constructor(props) {
    super(props)

    this.state = {
      username: '',
      email: '',
      password: ''
    }
  }

  render() {
    const changeState = (field) => {
      return ({target: {value}}) => {
        this.setState({
          [field]: value
        });
      };
    };

    return <Form horizontal>
      <FormGroup controlId='username'>
        <Col componentClass={ControlLabel} sm={2}>
          Username
        </Col>
        <Col sm={10}>
          <FormControl type='text' value={this.state.username} onChange={changeState('username')} />
        </Col>
      </FormGroup>

      <FormGroup controlId='email'>
        <Col componentClass={ControlLabel} sm={2}>
          Email
        </Col>
        <Col sm={10}>
          <FormControl type='email' value={this.state.email} onChange={changeState('email')} />
        </Col>
      </FormGroup>

      <FormGroup controlId='password'>
        <Col componentClass={ControlLabel} sm={2}>
          Password
        </Col>
        <Col sm={10}>
          <FormControl type='password' value={this.state.password} onChange={changeState('password')} />
        </Col>
      </FormGroup>

      <FormGroup>
        <Col smOffset={2} sm={10}>
          <Button type='button' bsStyle='primary' onClick={this.onSubmit.bind(this)}>
            Sign up
          </Button>
        </Col>
      </FormGroup>
    </Form>;
  }

  onSubmit() {
    requestJson('/accounts', {
      method: 'POST',
      body: _.pick(this.state, 'username', 'email', 'password')
    }).then( () => {
      return requestJson('/sessions', {
        method: 'POST',
        body: _.pick(this.state, 'username', 'password')
      }).then( session => {
        localStorage.setItem('sessionToken', session.token);
      });
    }).catch( err => {
      alert(err.message);
    });
  }
}

class LoginTab extends Component {
  constructor(props) {
    super(props)

    this.state = {
      username: '',
      password: ''
    }
  }

  render() {
    const changeState = (field) => {
      return ({target: {value}}) => {
        this.setState({
          [field]: value
        });
      };
    };

    return <Form horizontal>
      <FormGroup controlId='username'>
        <Col componentClass={ControlLabel} sm={2}>
          Username
        </Col>
        <Col sm={10}>
          <FormControl type='text' value={this.state.username} onChange={changeState('username')} />
        </Col>
      </FormGroup>

      <FormGroup controlId='password'>
        <Col componentClass={ControlLabel} sm={2}>
          Password
        </Col>
        <Col sm={10}>
          <FormControl type='password' value={this.state.password} onChange={changeState('password')} />
        </Col>
      </FormGroup>

      <FormGroup>
        <Col smOffset={2} sm={10}>
          <Button type='button' bsStyle='primary' onClick={this.onSubmit.bind(this)}>
            Login
          </Button>
        </Col>
      </FormGroup>
    </Form>;
  }

  onSubmit() {
    return requestJson('/sessions', {
      method: 'POST',
      body: _.pick(this.state, 'username', 'password')
    }).then( session => {
      localStorage.setItem('sessionToken', session.token);
    }).catch( err => {
      alert(err.message);
    });
  }
}
