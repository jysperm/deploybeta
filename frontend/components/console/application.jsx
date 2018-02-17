import _ from 'lodash';
import 'event-source-polyfill';
import {Button, Table, ButtonGroup, Modal, FormGroup, ControlLabel, FormControl, HelpBlock} from 'react-bootstrap';
import {Label, DropdownButton, MenuItem} from 'react-bootstrap';
import React, {Component} from 'react';

import {alertError} from '../../lib/error';
import {FormComponent} from '../../lib/components';
import {requestJson} from '../../lib/request';

export default class ApplicationsTab extends Component {
  constructor(props) {
    super(props);

    this.state = {
      editingApp: null,
      buildingVersionApp: null,
      buildingProgressVersion: null
    };
  }

  render() {
    return <div>
      <ButtonGroup>
        <Button bsStyle='success' onClick={::this.onCreateApp}>Create App</Button>
      </ButtonGroup>
      <Table responsive>
        <thead>
          <tr>
            <th>Domain</th>
            <th>Instances</th>
            <th>Version</th>
            <th>State</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {this.props.apps.map( app => {
            return <tr key={app.name}>
              <td>{app.name}</td>
              <td>{app.instances}</td>
              <td className='application-version'>
                <div>
                  <Label bsStyle='primary'>{app.version || 'N/A'}</Label>
                </div>
                <div>
                  <ButtonGroup>
                    <Button onClick={this.onBuildVersion.bind(this, app.name)}>Build</Button>
                    <DropdownButton title='Deploy...' id='deploy-dropdown'>
                      {_.map(app.versions, 'tag').map( versionTag => {
                        return <MenuItem key={versionTag} eventKey={versionTag} onClick={this.onDeployVersion.bind(this, app.name, versionTag)}>{versionTag}</MenuItem>;
                      })}
                    </DropdownButton>
                  </ButtonGroup>
                </div>
              </td>
              <td>
                <ul className='application-state'>
                  {app.nodes && app.nodes.map( ({createdAt, state, versionTag}) => {
                    return <li key={createdAt}>
                      <Label bsStyle='primary'>{versionTag}</Label>
                      <strong>{state}</strong>
                    </li>;
                  })}
                </ul>
              </td>
              <td>
                <ButtonGroup>
                  <Button onClick={this.onEditApp.bind(this, app.name)}>Edit</Button>
                  <Button bsStyle='danger' onClick={this.onDeleteApp.bind(this, app.name)}>Delete</Button>
                </ButtonGroup>
              </td>
            </tr>;
          })}
        </tbody>
      </Table>

      {this.state.editingApp && <EditAppModal {...this.state.editingApp} onClose={::this.onAppEdited} />}
      {this.state.buildingVersionApp && <BuildVersionModal {...this.state.buildingVersionApp} onClose={::this.onBuildStarted} />}
      {this.state.buildingProgressVersion && <BuildProgressModal {...this.state.buildingProgressVersion} onClose={::this.onBuildFinished} />}

    </div>;
  }

  onCreateApp() {
    this.setState({
      editingApp: {}
    });
  }

  onEditApp(name) {
    this.setState({
      editingApp: _.find(this.props.apps, {name}),
    });
  }

  onBuildVersion(name) {
    this.setState({
      buildingVersionApp: _.find(this.props.apps, {name})
    });
  }

  onDeployVersion(name, versionTag) {
    return requestJson(`/apps/${name}/version`, {
      method: 'PUT',
      body: {
        tag: versionTag
      }
    }).catch(alertError);
  }

  onDeleteApp(name) {
    return requestJson(`/apps/${name}`, {
      method: 'DELETE'
    }).then( () => {
      this.props.onAppDeleted({name});
    }).catch(alertError);
  }

  onAppEdited(app) {
    this.setState({
      editingApp: null
    });

    this.props.onAppEdited(app);
  }

  onBuildStarted(version) {
    if (version) {
      this.setState({
        buildingProgressVersion: _.extend(version, {
          appName: this.state.buildingVersionApp.name
        })
      });
    }

    this.setState({
      buildingVersionApp: null
    });
  }

  onBuildFinished() {
    this.setState({
      buildingProgressVersion: null
    });
  }
}

class EditAppModal extends FormComponent {
  constructor(props) {
    super(props);

    this.state = {};
  }

  render() {
    const title = this.props.name ? `Edit ${this.props.name}` : 'Create new app';

    return <Modal show={true} onHide={this.props.onClose.bind(this, null)}>
      <Modal.Header closeButton>
        <Modal.Title>{title}</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <FormGroup controlId='app-name'>
          <ControlLabel>Name</ControlLabel>
          <FormControl disabled={this.props.name} type='text' {...this.linkField('name')} />
          <HelpBlock>Used as your domain, must be globally unique</HelpBlock>
        </FormGroup>
        <FormGroup controlId='git-repository'>
          <ControlLabel>Git Repository</ControlLabel>
          <FormControl type='text' {...this.linkField('gitRepository')} />
          <HelpBlock>Build version from this git repository</HelpBlock>
        </FormGroup>
      </Modal.Body>
      <Modal.Footer>
        {this.props.name ? null : <Button bsStyle='success' onClick={::this.onCreateApp}>Create</Button>}
        {this.props.name && <Button bsStyle='success' onClick={::this.onEditApp}>Edit</Button>}
      </Modal.Footer>
    </Modal>;
  }

  onCreateApp() {
    return requestJson('/apps', {
      method: 'POST',
      body: this.state
    }).then( app => {
      this.props.onClose(app);
    }).catch(alertError);
  }

  onEditApp() {
    return requestJson(`/apps/${this.props.name}`, {
      method: 'PATCH',
      body: this.state
    }).then( app => {
      this.props.onClose(app);
    }).catch(alertError);
  }
}

class BuildVersionModal extends FormComponent {
  constructor(props) {
    super(props);

    this.state = {};
  }

  render() {
    return <Modal show={true} onHide={this.props.onClose.bind(this, null)}>
      <Modal.Header closeButton>
        <Modal.Title>Build version for {this.props.name}</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <FormGroup controlId='git-tag'>
          <ControlLabel>Git tag</ControlLabel>
          <FormControl type='text' {...this.linkField('gitTag')} />
          <HelpBlock>Git branch, tag or commit hash</HelpBlock>
        </FormGroup>
      </Modal.Body>
      <Modal.Footer>
        <Button bsStyle='success' onClick={::this.onBuildVersion}>Build</Button>
      </Modal.Footer>
    </Modal>;
  }

  onBuildVersion() {
    return requestJson(`/apps/${this.props.name}/versions`, {
      method: 'POST',
      body: {
        gitTag: this.state.gitTag
      }
    }).then( version => {
      this.props.onClose(version);
    }).catch(alertError);
  }
}

class BuildProgressModal extends Component {
  constructor(props) {
    super(props);

    this.state = {
      events: []
    };
  }

  componentDidMount() {
    const {appName, tag} = this.props;

    const events = new EventSourcePolyfill(`/apps/${appName}/versions/${tag}/progress`, {
      headers: {
        Authorization: localStorage.getItem('sessionToken')
      }
    });

    events.addEventListener('message', ({data}) => {
      const log = JSON.parse(data);

      if (log.payload === 'Deploying: Building finished.') {
        events.close();
        return;
      }

      this.setState({
        events: _.uniqBy(this.state.events.concat(log), 'id')
      });
    });
  }

  render() {
    return <Modal show={true} onHide={this.props.onClose}>
      <Modal.Header closeButton>
        <Modal.Title>Building {this.props.tag}</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        {this.state.events.map( ({id, payload}) => {
          return <p key={id}>{JSON.parse(payload).stream}</p>;
        })}
      </Modal.Body>
    </Modal>;
  }
}
