import React, {Component} from 'react';
import {Row, Tabs, Tab, Button, Table, ButtonGroup, Modal, FormGroup, ControlLabel, FormControl, HelpBlock} from 'react-bootstrap';

import {requestJson} from '../lib/request'
import LayoutView from './layout'

export default class ConsoleView extends Component {
  constructor(props) {
    super(props)

    this.state = {
      apps: []
    }
  }

  componentDidMount() {
    return requestJson('/apps').then( apps => {
      this.setState({apps});
    }).catch( err => {
      alert(err.message);
    });
  }

  render() {
    return <LayoutView>
      <Row>
        <Tabs defaultActiveKey={1} id='pages'>
          <Tab eventKey={1} title='Applications'>
            <ApplicationsTab apps={this.state.apps} onAppCreated={this.onAppCreated.bind(this)} />
          </Tab>
        </Tabs>
      </Row>
    </LayoutView>;
  }

  onAppCreated(app) {
    this.setState({
      apps: this.state.apps.concat(app)
    });
  }
}

class ApplicationsTab extends Component {
  constructor(props) {
    super(props)

    this.state = {
      creatingApp: false,
      creatingAppName: ''
    }
  }

  render() {
    return <div>
      <ButtonGroup>
        <Button bsStyle='success' onClick={this.onCreatingApp.bind(this)}>Create App</Button>
      </ButtonGroup>
      <Table responsive>
        <thead>
          <tr>
            <th>Domain</th>
            <th>Version</th>
            <th>Instances</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {this.props.apps.map( app => {
            return <tr key={app.name}>
              <td>{app.name}</td>
              <td>{app.version}</td>
              <td>{app.instances}</td>
              <td>
                <ButtonGroup>
                  <Button bsStyle='danger'>Delete</Button>
                </ButtonGroup>
              </td>
            </tr>;
          })}
        </tbody>
      </Table>

      <Modal show={this.state.creatingApp} onHide={this.onCreateAppModalClose.bind(this)}>
        <Modal.Header closeButton>
          <Modal.Title>Create new app</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <FormGroup controlId='app-name'>
            <ControlLabel>Name</ControlLabel>
            <FormControl type='text' onChange={this.onCreatingAppNameEdited.bind(this)} />
            <HelpBlock>Used as your domain, must be globally unique</HelpBlock>
          </FormGroup>
        </Modal.Body>
        <Modal.Footer>
          <Button bsStyle='success' onClick={this.onCreateApp.bind(this)}>Create</Button>
        </Modal.Footer>
      </Modal>;
    </div>;
  }

  onCreatingApp() {
    this.setState({creatingApp: true});
  }

  onCreateAppModalClose() {
    this.setState({creatingApp: false});
  }

  onCreatingAppNameEdited({target: {value}}) {
    this.setState({creatingAppName: value});
  }

  onCreateApp() {
    return requestJson('/apps', {
      method: 'POST',
      body: {
        name: this.state.creatingAppName
      }
    }).then( app => {
      this.setState({
        creatingApp: false,
        creatingAppName: ''
      });

      this.props.onAppCreated(app);
    }).catch( err => {
      alert(err.message);
    });
  }
}
