import React from 'react';
import './versionBar.css'

export class VersionBar extends React.Component {
    constructor(props) {
      super(props)
    }
  
    render() {
      const props = this.props;
      console.dir(this.props)
  
      return (
      <div className="versionBar">
        <p className="version">{props.version}</p>
      </div>
    )
    }
  }
  