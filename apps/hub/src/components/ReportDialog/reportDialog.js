import React from 'react'
import { DialogContainer, Toolbar, Button, Paper, Divider, TextField, SelectionControl } from 'react-md';
import { formatLatitude, formatLongitude } from 'latlon-formatter';

import FacilityChip from './facilityChip';


import './reportDialog.css'

export class ReportDialog extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      report: null,
      changed: false,
      showAll: false,
      fields: [],
      shortNameError: false,
    }
  }

  componentWillReceiveProps = (props) => {
    this.setState({ report: this.props.report, showAll: props.showFields });
  }

  onShow = (report) => {
    this.setState({ report });
  };

  toggleShowEmpty = (newState) => {
    this.setState({ showAll: newState })
  }

  fieldToKey = (field) => {
    field = field.charAt(0).toLowerCase() + field.slice(1)
    field = field.replace(/\s/g, '');
    return field
  }

  onValueChange = (field, value, e) => {
    let report = this.state.report;
    field = this.fieldToKey(field);

    if (field === "shortName") {
      const shortNameRegex = /^[A-Z]{3}\s[0-9]{2}$/;
      const isValid = shortNameRegex.test(value);
      if (value !== "" && !isValid) {
        this.setState({ shortNameError: true });
        return;
      } else {
        report[field] = value;
        this.setState({ report, changed: true, shortNameError: false });
        return;
      }
    }

    if (field === "title" || field === "country" || field === "facilities") {
      report[field] = value;
    } else {
      report.details[field] = value;
    }
    this.setState({ report, changed: true });
  };

  removeFacility = (facility) => {
    let report = this.state.report;
    const facilities = report.facilities.slice();

    facilities.splice(facilities.indexOf(facility), 1);
    report.facilities = facilities;
    this.setState({ report })
  }

  addFacilityOnTab = (e) => {
    e.keyCode = 9;
    this.addFacility(e)
  }

  addFacility = (e) => {
    if (e.key === 'Enter' || e.keyCode === 9 || e.key === ' ') {
      const value = e.target.value.trim().toUpperCase();
      if (value != "") {
        e.target.value = "";

        let report = this.state.report;
        let facilities = report.facilities ? report.facilities : []
        facilities.push(value)
        report.facilities = facilities;
        this.setState({ report })
      }
    }
  }


  render() {
    if (this.props.report == null) {
      return null;
    } else {
      const report = this.props.report;
      const lat = formatLatitude(report.location.lat, { degrees: true });
      const lng = formatLongitude(report.location.lng, { degrees: true });
      const facilities = this.state.report ? this.state.report.facilities : [];
      const chips = facilities && facilities.length > 0 ? facilities.map(facility => <FacilityChip key={facility} name={facility} onClick={() => this.removeFacility(facility)} />) : ""

      return <DialogContainer visible={report !== null}
        id="locationDialog"
        className="locationDialog"
        aria-label="Location Dialog"
        onHide={this.props.onHide}
        onShow={() => this.onShow(this.props.report)}>
        <Toolbar fixed
          colored
          title="Location details"
          nav={<Button icon onClick={() => this.props.onHide(null)}>close</Button>}
          actions={<Button flat onClick={() => this.props.onHide(report)} disabled={this.state.shortNameError}>Save</Button>}
        />
        <Paper className="fillParent md-toolbar-relative" zDepth={1}>
          <SelectionControl
            id="showEmpty"
            type="switch"
            label="Show Empty Fields"
            name="showEmpty"
            className="showEmpty"
            checked={this.state.showAll}
            onChange={(newState) => this.toggleShowEmpty(newState)}
          />
          <p className="location">
            Lat: {lat} - Lng: {lng}
          </p>
          <TextField
            className="locationTitle"
            label="Name"
            id="locationTitle"
            defaultValue={report.title ? report.title : ""}
            onChange={this.onValueChange.bind(this, "title")}
          />
          <TextField
            className="shortName dataInput"
            label="Short Name"
            id="shortName"
            defaultValue={report.shortName ? report.shortName : ""}
            onChange={this.onValueChange.bind(this, "shortName")}
            error={this.state.shortNameError}
            errorText="Short name should be in format LLL NN."
          />
          <TextField
            className="country"
            label="Location/Country"
            id="country"
            defaultValue={report.country ? report.country : ""}
            onChange={this.onValueChange.bind(this, "country")}
          />
          <TextField
            className="description"
            label="Description"
            id="description"
            defaultValue={report.details.description ? report.details.description : ""}
            onChange={this.onValueChange.bind(this, "description")}
          />
          {chips}
          <TextField
            className="facilities"
            label="Facilities"
            id="facilities"
            onBlur={this.addFacilityOnTab.bind(this)}
            onKeyUp={this.addFacility.bind(this)}
          />
          <Divider />
          {[].concat(this.props.fields)
            .sort((a, b) => a.order > b.order)
            .map((field, index) => {
              var fieldKey = this.fieldToKey(field.value)
              if (this.state.showAll || report.details[fieldKey]) {
                return <TextField
                  className="dataInput"
                  key={index}
                  label={field.value}
                  id={fieldKey}
                  defaultValue={report.details[fieldKey] ? report.details[fieldKey] : ""}
                  onChange={this.onValueChange.bind(this, fieldKey)}
                />;
              }
            })

          }


        </Paper>
      </DialogContainer>
    }
  }
}
