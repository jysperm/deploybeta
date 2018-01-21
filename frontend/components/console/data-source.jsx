import React, {Component} from 'react';
import {Button, Table, ButtonGroup, Modal} from 'react-bootstrap';
import {FormGroup, ControlLabel, FormControl, HelpBlock} from 'react-bootstrap';

import {alertError} from '../../lib/error';
import {FormComponent} from '../../lib/components';
import {requestJson} from '../../lib/request';

export default class DataSourcesTab extends Component {
  constructor(props) {
    super(props)

    this.state = {
      editingDataSource: null
    }
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

  }

  onDeleteDataSource(dataSourceName) {

  }
}

class EditDataSourceModal extends FormComponent {
  constructor(props) {
    super(props)

    this.state = {};
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
          <FormControl componentClass="select">
            <option value='mongodb'>MongoDB</option>
            <option value='redis'>Redis</option>
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

class EditDataSourceLinksModal extends FormComponent {

}
