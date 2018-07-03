import React from 'react'
import {TextField} from 'react-md';
import './searchBar.css';

export class SearchBar extends React.Component {
  constructor(props) {
    super(props)
  }

  render() {
    const props = this.props;

    return <TextField className="search-bar" id="searchBar" label="Filter" placeholder="Enter criteria" value={props.value} onChange={props.onChange}/>
  }
}
