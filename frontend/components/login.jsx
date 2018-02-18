import _ from 'lodash';
import {Row, Col, Tabs, Tab, Form, FormGroup, FormControl, ControlLabel, Button} from 'react-bootstrap';
import React, {Component} from 'react';

import {alertError} from '../lib/error';
import {FormComponent} from '../lib/components';
import {requestJson} from '../lib/request';
import LayoutView from './layout';

export default class LoginView extends Component {
  render() {
    return <LayoutView pageHeader={true}>
      <Row>
        <Tabs animation={false} defaultActiveKey='login' id='signup-or-login'>
          <Tab eventKey='signup' title='Sign up'>
            <SignupTab onLogged={::this.onLogged} />
          </Tab>
          <Tab eventKey='login' title='Login'>
            <LoginTab onLogged={::this.onLogged} />
          </Tab>
        </Tabs>
      </Row>
    </LayoutView>;
  }

  onLogged(session) {
    localStorage.setItem('username', session.username);
    localStorage.setItem('sessionToken', session.token);
    this.props.history.push('/console');
  }
}

class SignupTab extends FormComponent {
  constructor(props) {
    super(props);

    this.state = {
      username: '',
      email: '',
      password: ''
    };
  }

  render() {
    return <Form horizontal>
      <FormGroup controlId='signup-username'>
        <Col componentClass={ControlLabel} sm={2}>
          Username
        </Col>
        <Col sm={10}>
          <FormControl type='text' {...this.linkField('username')} />
        </Col>
      </FormGroup>

      <FormGroup controlId='signup-email'>
        <Col componentClass={ControlLabel} sm={2}>
          Email
        </Col>
        <Col sm={10}>
          <FormControl type='email' {...this.linkField('email')} />
        </Col>
      </FormGroup>

      <FormGroup controlId='signup-password'>
        <Col componentClass={ControlLabel} sm={2}>
          Password
        </Col>
        <Col sm={10}>
          <FormControl type='password' {...this.linkField('password')} />
        </Col>
      </FormGroup>

      <FormGroup>
        <Col smOffset={2} sm={10}>
          <Button type='button' bsStyle='primary' onClick={::this.onSubmit}>
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
        this.props.onLogged(session);
      });
    }).catch(alertError);
  }
}

class LoginTab extends FormComponent {
  constructor(props) {
    super(props);

    this.state = {
      username: '',
      password: ''
    };
  }

  render() {
    return <Form horizontal>
      <FormGroup controlId='login-username'>
        <Col componentClass={ControlLabel} sm={2}>
          Username
        </Col>
        <Col sm={10}>
          <FormControl type='text' {...this.linkField('username')} />
        </Col>
      </FormGroup>

      <FormGroup controlId='login-password'>
        <Col componentClass={ControlLabel} sm={2}>
          Password
        </Col>
        <Col sm={10}>
          <FormControl type='password' {...this.linkField('password')} />
        </Col>
      </FormGroup>

      <FormGroup>
        <Col smOffset={2} sm={10}>
          <Button type='button' bsStyle='primary' onClick={::this.onSubmit}>
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
      this.props.onLogged(session);
    }).catch(alertError);
  }
}
