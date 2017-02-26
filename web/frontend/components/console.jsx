import React, {Component} from 'react';
import {Row, Tabs, Tab, Button} from 'react-bootstrap';

import LayoutView from './layout'

export default class LoginView extends Component {
  render() {
    return <LayoutView>
      <Row>
        <Tabs defaultActiveKey={1} id='pages'>
          <Tab eventKey={1} title='Applications'>
            <ApplicationsTab />
          </Tab>
        </Tabs>
      </Row>
    </LayoutView>;
  }
}

class ApplicationsTab extends Component {
  render() {
    return <p>ApplicationsTab</p>;
  }
}
