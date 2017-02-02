import React, {Component} from 'react';
import ReactDOM from 'react-dom';

class DeployingView extends Component {
  render() {
    return <p>Deploying</p>;
  }
}

ReactDOM.render(React.createElement(DeployingView), document.getElementById('react-root'));
