import React from 'react';
import './versionBar.css'

export class VersionBar extends React.Component {
    constructor(props) {
      super(props)
    }
  
    render() {  
      return (
      <div className="versionBar">
        <p className="version">{this.props.version}</p>
      </div>
    )
    }
  }
  