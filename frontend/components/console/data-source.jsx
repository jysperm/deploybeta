import {Button, Table, ButtonGroup, Modal} from 'react-bootstrap';
import {FormGroup, ControlLabel, FormControl, HelpBlock, Checkbox} from 'react-bootstrap';
import React, {Component} from 'react';
import _ from 'lodash';

import {alertError} from '../../lib/error';
import {FormComponent} from '../../lib/components';
import {requestJson} from '../../lib/request';

export default class DataSourcesTab extends Component {
  constructor(props) {
    super(props);

    this.state = {
      editingDataSource: null,
      editingLinks: null
    };
  }

  render() {
    return <div>
      <ButtonGroup>
        <Button bsStyle='success' onClick={::this.onCreateDataSource}>Create Data Source</Button>
      </ButtonGroup>
      <Table responsive>
        <thead>
          <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Instances</th>
            <th>Links</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {this.props.dataSources.map( dataSource => {
            return <tr key={dataSource.name}>
              <td>{dataSource.name}</td>
              <td>{dataSource.type}</td>
              <td>{dataSource.instances}</td>
              <td>
                <Button bsStyle='info' onClick={this.onEditingLinks.bind(this, dataSource.name)}>Edit Links</Button>
              </td>
              <td>
                <ButtonGroup>
                  <Button bsStyle='danger' onClick={this.onDeleteDataSource.bind(this, dataSource.name)}>Delete</Button>
                </ButtonGroup>
              </td>
            </tr>;
          })}
        </tbody>
      </Table>

      {this.state.editingDataSource && <EditDataSourceModal {...this.state.editingDataSource} onClose={::this.onDataSourceEdited} />}
      {this.state.editingLinks && <EditLinksModal apps={this.props.apps} dataSource={this.state.editingLinks} onClose={::this.onLinksEdited} />}

    </div>;
  }

  onCreateDataSource() {
    this.setState({
      editingDataSource: {}
    });
  }

  onDataSourceEdited(dataSource) {
    this.setState({
      editingDataSource: null
    });

    this.props.onDataSourceEdited(dataSource);
  }

  onEditingLinks(dataSourceName) {
    this.setState({
      editingLinks: _.find(this.props.dataSources, {name: dataSourceName})
    });
  }

  onLinksEdited() {
    this.setState({
      editingLinks: null
    });
  }

  onDeleteDataSource(name) {
    return requestJson(`/data-sources/${name}`, {
      method: 'DELETE'
    }).then( () => {
      this.props.onDataSourceDeleted({name});
    }).catch(alertError);
  }
}

class EditDataSourceModal extends FormComponent {
  constructor(props) {
    super(props);

    this.state = {
      type: 'mongodb'
    };
  }

  render() {
    return <Modal show={true} onHide={this.props.onClose.bind(this, null)}>
      <Modal.Header closeButton>
        <Modal.Title>Create new Data Source</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <FormGroup controlId='data-source-name'>
          <ControlLabel>Name</ControlLabel>
          <FormControl disabled={this.props.name} type='text' {...this.linkField('name')} />
          <HelpBlock>Internal access domain, must be globally unique</HelpBlock>
        </FormGroup>
        <FormGroup controlId='type'>
          <ControlLabel>Database Type</ControlLabel>
          <FormControl componentClass="select" {...this.linkField('type')}>
            <option value='mongodb'>MongoDB</option>
            <option value='redis'>Redis</option>
            <option value='mysql'>MySQL</option>
          </FormControl>
          <HelpBlock>Can not be modified after creation</HelpBlock>
        </FormGroup>
      </Modal.Body>
      <Modal.Footer>
        <Button bsStyle='success' onClick={::this.onCreateDataSource}>Create</Button>
      </Modal.Footer>
    </Modal>;
  }

  onCreateDataSource() {
    return requestJson('/data-sources', {
      method: 'POST',
      body: this.state
    }).then( dataSource => {
      this.props.onClose(dataSource);
    }).catch(alertError);
  }
}

class EditLinksModal extends FormComponent {
  constructor(props) {
    super(props);

    this.state = {
      linkedApps: _.clone(this.props.dataSource.linkedApps)
    };
  }

  render() {
    return <Modal show={true} onHide={this.props.onClose.bind(this, null)}>
      <Modal.Header closeButton>
        <Modal.Title>Link {this.props.dataSource.name} to ...</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        {this.props.apps.map( app => {
          let checked = _.includes(this.state.linkedApps, app.name);

          return <Checkbox key={app.name} value={checked} onChange={this.onCheckChanged.bind(this, app.name)}>
            {app.name}
          </Checkbox>;
        })}
      </Modal.Body>
      <Modal.Footer>
        <Button onClick={::this.onCloseClicked}>Close</Button>
      </Modal.Footer>
    </Modal>;
  }

  onCheckChanged(appName, {target}) {
    return requestJson(`/data-sources/${this.props.dataSource.name}/links/${appName}`, {
      method: target.checked ? 'POST' : 'DELETE'
    }).then( () => {
      if (target.checked) {
        this.setState({
          linkedApps: _.union(this.state.linkedApps, [appName])
        });
      } else {
        this.setState({
          linkedApps: _.without(this.state.linkedApps, appName)
        });
      }
    }).catch(alertError);
  }

  onCloseClicked() {
    this.props.onClose();
  }
}
