import React, { Component } from 'react';
import { Button } from 'react-md';
import { withRouter } from 'react-router-dom';


import './menuBar.css'

class MenuBar extends Component {

    nextPath = async (path) => {
        if (path == '/') {
            if (this.props.location.pathname === '/') {
                this.props.resetMap()
            }
            if (this.props.location.pathname === '/settings') {
                if (await this.props.homeAction()) {
                    this.props.history.push(path)
                } else {
                    return
                }

            }
        } else {
            this.props.history.push(path)
        }

    }

    render() {
        return (
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