import React from 'react'
import {DialogContainer, Toolbar, Button, Grid, Cell, Paper, TextField} from 'react-md';
import {formatLatitude, formatLongitude} from 'latlon-formatter';
import * as _ from 'lodash';

import './reportDialog.css'

export class ReportDialog extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      report: null
    }
  }

  onShow = (report) => {
    this.setState({report});
  };

  reportFieldChange = (reportID, value, event) => {
    const newReport = _.cloneDeep(this.state.report);
    newReport.reports[reportID] = {reportID, reportDetails: value};
    this.setState({report: newReport});
  };

  renderChildReports = (report) => {
    return report.reports.map(childReport => {
      return <Grid>
        <Cell size={12}>
          <TextField label={`Report ${childReport.reportID}`}
                     value={childReport.reportDetails}
                     rows={2}
                     maxRows={6}
                     onChange={(value, event) => this.reportFieldChange(childReport.reportID, value, event)}
          />
          <br/>
          <a href={childReport.url}>{childReport.url}</a>
        </Cell>
      </Grid>
    });
  };

  renderDetails = (report) => {
    return Object.entries(report.details).map(([key, value]) => {
      return <Grid>
        <Cell size={4}>{key}</Cell>
        <Cell size={8}>{value}</Cell>
      </Grid>
    });
  };

  render() {
    if (this.props.report == null) {
      return null;
    } else {
      const report = this.state.report || this.props.report;
      const lat = formatLatitude(report.location.lat, {degrees: true});
      const lon = formatLongitude(report.location.lng, {degrees: true});
      console.dir(this.props.report)

      return <DialogContainer visible={this.props.report !== null}
                              id="reportDialog"
                              className="reportDialog"
                              onHide={this.props.onHide}
                              onShow={() => this.onShow(this.props.report)}>
        <Toolbar fixed
                 colored
                 title="Report details"
                 nav={<Button icon onClick={() => this.props.onHide(null)}>close</Button>}
                 actions={<Button flat onClick={() => this.props.onHide(report)}>Save</Button>}
        />
        <Paper className="fillParent md-toolbar-relative" zDepth={1}>
          <h1>Details</h1>
          <h2>{report.title}</h2>
          <h3>{lat} {lon}</h3>
          {this.props.report.reports !== undefined &&
          <Grid>
            <Cell size={6}>
              <h2>Reports</h2>
              {this.renderChildReports(report)}
            </Cell>
            <Cell size={6}>
              <h2>Details</h2>
              {this.renderDetails(report)}
            </Cell>
          </Grid>}


        </Paper>
      </DialogContainer>
    }
  }
}
