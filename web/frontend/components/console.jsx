import React, {Component} from 'react';
import {Row, Tabs, Tab, Button, Table, ButtonGroup} from 'react-bootstrap';

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
            <ApplicationsTab apps={this.state.apps} />
          </Tab>
        </Tabs>
      </Row>
    </LayoutView>;
  }
}

class ApplicationsTab extends Component {
  render() {
    return <div>
      <ButtonGroup>
        <Button bsStyle='success'>Create App</Button>
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
    </div>;
  }
}
