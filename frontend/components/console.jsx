import _ from 'lodash';
import {Row, Tabs, Tab} from 'react-bootstrap';
import React, {Component} from 'react';

import {requestJson} from '../lib/request';
import LayoutView from './layout';
import ApplicationsTab from './console/application';

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
            <ApplicationsTab apps={this.state.apps} onAppEdited={::this.onAppEdited} />
          </Tab>
        </Tabs>
      </Row>
    </LayoutView>;
  }

  onAppEdited(app) {
    if (app) {
      this.setState({
        apps: [app].concat(_.reject(this.state.apps, {name: app.name}))
      });
    }
  }
}
