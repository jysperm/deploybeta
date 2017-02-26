import React, {Component} from 'react';
import {Grid, Row, PageHeader} from 'react-bootstrap';

export default class LayoutView extends Component {
  render() {
    return <Grid>
      <Row>
        <PageHeader>
          Deploying <small>A simple container platform based on reliable solutions.</small>
        </PageHeader>
      </Row>
      {this.props.children}
    </Grid>;
  }
}
