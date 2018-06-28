import React from 'react'
import { DialogContainer, Toolbar, Button, Paper, Divider, TextField, SelectionControl } from 'react-md';
import { formatLatitude, formatLongitude } from 'latlon-formatter';


import './reportDialog.css'

export class ReportDialog extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      report: null,
      changed: false,
      showAll: false,
      fields: []
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
    if (field === "title" || field === "country") {
      report[field] = value
    } else {
      report.details[field] = value;
    }
    this.setState({ report, changed: true });
  };


  render() {
    if (this.props.report == null) {
      return null;
    } else {
      const report = this.props.report;
      const lat = formatLatitude(report.location.lat, { degrees: true });
      const lng = formatLongitude(report.location.lng, { degrees: true });

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
          actions={<Button flat onClick={() => this.props.onHide(report)}>Save</Button>}
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
