import React from 'react'
import {TextField} from 'react-md';

export class SearchBar extends React.Component {
  constructor(props) {
    super(props)
  }

  render() {
    const props = this.props;

    return <TextField id="searchBar" label="Filter" placeholder="Enter criteria" value={props.value} onChange={props.onChange}/>
  }
}
