import React, { Component } from 'react';
import {  Button } from 'react-md';
import { withRouter } from 'react-router-dom';


import './menuBar.css'

class MenuBar extends Component {

    nextPath = (path) => {
        if (this.props.homeAction != null){
            this.props.homeAction()
        }
        this.props.history.push(path)
    }

    render() {
        return(
            <div className="menuBar">
                <Button className="homeIcon" id="Home" icon onClick={() => this.nextPath('/')}>
                   home
                </Button>
                <Button className="settingsIcon" id="Settings" icon onClick={() => this.nextPath('/settings')}>
                   settings
                </Button>
            </div>
        )
    }
}

export default withRouter(MenuBar);