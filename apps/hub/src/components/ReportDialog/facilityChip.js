import React, {PureComponent} from 'react';
import { Avatar, Chip } from 'react-md';

export default class FacilityChip extends PureComponent {
    handleRemove = () => {
        this.props.onClick(this.props.state);
      };
    
      render() {
        return (
          <Chip
            className="facilities__chip"
            onClick={this.handleRemove}
            removable
            label={this.props.name}
            avatar={<Avatar random></Avatar>}
          />
        );
      }
}