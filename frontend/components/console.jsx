import React, {Component} from 'react';
import {Row, Tabs, Tab, Button, Table, ButtonGroup, Modal, FormGroup, ControlLabel, FormControl, HelpBlock} from 'react-bootstrap';
import {Label, DropdownButton, MenuItem, Checkbox} from 'react-bootstrap';

import {requestJson} from '../lib/request';
import LayoutView from './layout';

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
            <ApplicationsTab apps={this.state.apps} onAppCreated={::this.onAppCreated} />
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
      creatingAppName: '',

      buildingVersion: false,
      buildingVersionDeploy: false,
      buildingVersionTag: 'master'
    }
  }

  render() {
    return <div>
      <ButtonGroup>
        <Button bsStyle='success' onClick={::this.onCreatingApp}>Create App</Button>
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
              <td>
                <Label bsStyle='primary'>{app.version}</Label>
                <ButtonGroup>
                  <Button onClick={this.onBuildingVersion.bind(this, app.name)}>Build</Button>
                  <DropdownButton title='Deploy...' id='deploy-dropdown'>
                    <MenuItem eventKey={1}>20170219-205301</MenuItem>
                  </DropdownButton>
                </ButtonGroup>
              </td>
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

      <Modal show={this.state.creatingApp} onHide={::this.onCreateAppModalClose}>
        <Modal.Header closeButton>
          <Modal.Title>Create new app</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <FormGroup controlId='app-name'>
            <ControlLabel>Name</ControlLabel>
            <FormControl type='text' onChange={::this.onCreatingAppNameEdited} value={this.state.creatingAppName} />
            <HelpBlock>Used as your domain, must be globally unique</HelpBlock>
          </FormGroup>
        </Modal.Body>
        <Modal.Footer>
          <Button bsStyle='success' onClick={::this.onCreateApp}>Create</Button>
        </Modal.Footer>
      </Modal>

      <Modal show={this.state.buildingVersion} onHide={::this.onBuildVersionModalClose}>
        <Modal.Header closeButton>
          <Modal.Title>Build new version</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <FormGroup controlId='git-tag'>
            <ControlLabel>App</ControlLabel>
            <FormControl type='text' disabled={true} value={this.state.buildingVersion}  />
          </FormGroup>
          <FormGroup controlId='git-tag'>
            <ControlLabel>Git tag</ControlLabel>
            <FormControl type='text' onChange={::this.onBuildingVersionTagEdited} value={this.state.buildingVersionTag}  />
            <HelpBlock>Git branch, tag or commit hash</HelpBlock>
          </FormGroup>
          <Checkbox onChange={::this.onBuildVersionDeployEdited} value={this.state.buildingVersionDeploy}>Deploy to app after build finished</Checkbox>
        </Modal.Body>
        <Modal.Footer>
          <Button bsStyle='success' onClick={::this.onBuildVersion}>Build</Button>
        </Modal.Footer>
      </Modal>
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

  onBuildingVersion(appName) {
    this.setState({buildingVersion: appName});
  }

  onBuildVersionModalClose() {
    this.setState({buildingVersion: false});
  }

  onBuildingVersionTagEdited({target: {value}}) {
    this.setState({buildingVersionTag: value});
  }

  onBuildVersionDeployEdited({target: {checked}}) {
    this.setState({buildingVersionDeploy: checked});
  }

  onBuildVersion() {
    const {buildingVersion, buildingVersionDeploy} = this.state;

    return requestJson(`/apps/${buildingVersion}/${buildingVersionDeploy ? 'version' : 'versions'}`, {
      method: 'POST',
      body: {
        gitTag: this.state.buildingVersionTag
      }
    }).then( app => {
      this.setState({
        buildingVersion: false
      });
    }).catch( err => {
      alert(err.message);
    });
  }
}
