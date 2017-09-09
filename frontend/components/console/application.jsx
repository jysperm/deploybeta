import _ from 'lodash';
import {Button, Table, ButtonGroup, Modal, FormGroup, ControlLabel, FormControl, HelpBlock} from 'react-bootstrap';
import {Label, DropdownButton, MenuItem, Checkbox} from 'react-bootstrap';
import React, {Component} from 'react';

import {FormComponent} from '../../lib/components';
import {requestJson} from '../../lib/request';

export default class ApplicationsTab extends Component {
  constructor(props) {
    super(props)

    this.state = {
      editingApp: false,
      buildingVersionApp: false
    }
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
            <th>Version</th>
            <th>Instances</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {this.props.apps.map( app => {
            let versionTags = _.map(app.versions, 'tag');

            return <tr key={app.name}>
              <td>{app.name}</td>
              <td>
                <Label bsStyle='primary'>{app.version}</Label>
                <ButtonGroup>
                  <Button onClick={this.onBuildVersion.bind(this, app.name)}>Build</Button>
                  <DropdownButton title='Deploy...' id='deploy-dropdown'>
                    {_.map(app.versions, 'tag').map( versionTag => {
                      return <MenuItem eventKey={versionTag} onClick={this.onDeployVersion.bind(this, app.name, versionTag)}>{versionTag}</MenuItem>
                    })}
                  </DropdownButton>
                </ButtonGroup>
              </td>
              <td>{app.instances}</td>
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
      {this.state.buildingVersionApp && <BuildVersionModal {...this.state.buildingVersionApp} onClose={::this.onAppEdited} />}

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
    }).catch( err => {
      alert(err.message);
    });
  }

  onDeleteApp(name) {
    return requestJson(`/apps/${name}`, {
      method: 'DELETE'
    }).then( () => {
      this.props.onAppDeleted({name});
    }).catch( err => {
      alert(err.message);
    });
  }

  onAppEdited(app) {
    this.setState({
      editingApp: false,
      buildingVersionApp: false
    });

    this.props.onAppEdited(app);
  }
}

class EditAppModal extends FormComponent {
  constructor(props) {
    super(props)

    this.state = {};
  }

  render() {
    const title = this.props.name ? `Edit ${this.props.name}` : 'Create new app';

    return <Modal show={true} onHide={::this.props.onClose}>
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
    }).catch( err => {
      alert(err.message);
    });
  }

  onEditApp() {
    return requestJson(`/apps/${this.props.name}`, {
      method: 'PATCH',
      body: this.state
    }).then( app => {
      this.props.onClose(app);
    }).catch( err => {
      alert(err.message);
    });
  }
}

class BuildVersionModal extends FormComponent {
  constructor(props) {
    super(props)

    this.state = {};
  }

  render() {
    return <Modal show={true} onHide={::this.props.onClose}>
      <Modal.Header closeButton>
        <Modal.Title>Build version for {this.props.name}</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <FormGroup controlId='git-tag'>
          <ControlLabel>Git tag</ControlLabel>
          <FormControl type='text' {...this.linkField('gitTag')} />
          <HelpBlock>Git branch, tag or commit hash</HelpBlock>
        </FormGroup>
        <Checkbox {...this.linkField('buildAndDeploy', 'checked')}>
          Deploy to app after build finished
        </Checkbox>
      </Modal.Body>
      <Modal.Footer>
        <Button bsStyle='success' onClick={::this.onBuildVersion}>Build</Button>
      </Modal.Footer>
    </Modal>;
  }

  onBuildVersion() {
    return requestJson(`/apps/${this.props.name}/${this.state.buildAndDeploy ? 'version' : 'versions'}`, {
      method: 'POST',
      body: {
        gitTag: this.state.gitTag
      }
    }).then( () => {
      this.props.onClose();
    }).catch( err => {
      alert(err.message);
    });
  }
}
