import React, {Component} from 'react';
import {Grid, Row, PageHeader} from 'react-bootstrap';

export default class LayoutView extends Component {
  render() {
    const pageHeader = <Row>
      <PageHeader>
        Deploybeta <small>A simple container platform based on reliable solutions.</small>
      </PageHeader>
    </Row>;

    return <Grid>
      {this.props.pageHeader && pageHeader}
      {this.props.children}
    </Grid>;
  }
}
