import React, {Component} from 'react';
import {Grid, Row, PageHeader} from 'react-bootstrap';

export default class LayoutView extends Component {
  render() {
    return <Grid>
      <Row>
        <PageHeader>
          Deploying <small>Containerized platform based on Docker Swarm, Openresty and ELK</small>
        </PageHeader>
      </Row>
      {this.props.children}
    </Grid>;
  }
}
