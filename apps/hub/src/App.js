import React, { Component } from 'react';
import './App.css';
import { ToastContainer, toast } from 'react-toastify';
import Home from './pages/Home/home'
import Settings from './pages/Settings/settings'
import './react-toastify.css';
import {
  BrowserRouter as Router,
  Route,
  Link
} from 'react-router-dom'



class App extends Component {

  render() {
    return(
      <Router>
        <div>
          <Route exact path="/" component={Home}/>
          <Route path="/settings" component={Settings}/>
        </div>
      </Router>
    )
  }
  
}

export default App;
