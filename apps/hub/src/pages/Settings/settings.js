import React, { Component } from 'react'
import { ToastContainer, toast } from 'react-toastify';
import MenuBar from '../../components/MenuBar/menuBar';
import {VersionBar} from '../../components/VersionBar/versionBar';

import './settings.css'

class Settings extends Component {

    constructor(props) {
        super(props);
    
        this.state = {
          version: process.env.REACT_APP_VERSION ? process.env.REACT_APP_VERSION : "0.0.1",
          fields: []
        };
      }

    render() {

        return(
        <div>
            <ToastContainer />
            <VersionBar version={this.state.version}/>
            <MenuBar />
            <h1>Settings</h1>
        </div>
        )
    }
}

export default Settings