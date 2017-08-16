import {Component} from 'react';

export class FormComponent extends Component {
  linkField(field, attribute = 'value') {
    return {
      value: this.state[field] || this.props[field],
      onChange: ({target}) => {
        this.setState({[field]: target[attribute]});
      }
    };
  }
}
