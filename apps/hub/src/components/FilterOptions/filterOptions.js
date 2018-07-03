import React, { Component } from 'react';
import { Checkbox } from 'react-md';
import './filterOptions.css';

export class FilterOptions extends Component {
    constructor(props) {
        super(props)
    }

    render() {
        return <div className="filter-options">
            <p className="filter-options__header md-subheader">Filter Options</p>
            <div className="filter-options__options">
                <Checkbox
                    id="filter-options__title"
                    name="title"
                    label="Title"
                    checked={this.props.filterState.title}
                    onChange={()=>this.props.toggleFilter("title")}
                />
                <Checkbox
                    id="filter-options__shortName"
                    name="shortName"
                    label="Short Name"
                    checked={this.props.filterState.shortName}
                    onChange={()=>this.props.toggleFilter("shortName")}
                />
                <Checkbox
                    id="filter-options__facilities"
                    name="facilities"
                    label="Facilities"
                    checked={this.props.filterState.facilities}
                    onChange={()=>this.props.toggleFilter("facilities")}
                />
                <Checkbox
                    id="filter-options__location"
                    name="locations"
                    label="Locations"
                    checked={this.props.filterState.locations}
                    onChange={()=>this.props.toggleFilter("locations")}
                />
            </div>
        </div>
    }

}