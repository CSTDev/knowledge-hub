import React from 'react'
import { Button, Paper, List, ListItem, Subheader } from 'react-md';
import { formatLatitude, formatLongitude } from 'latlon-formatter';

import FAEye from 'react-icons/lib/fa/eye';

import './detailsPane.css';

const uuidv4 = require('uuid/v4');

export class DetailsPane extends React.Component {
  constructor(props) {
    super(props)
  }

  renderItems(reports, showReport) {
    return reports.map(report => {
      const lat = formatLatitude(report.location.lat, { degrees: true });
      const lon = formatLongitude(report.location.lng, { degrees: true });
      const country = report.location.country ? report.location.country : "";
      const filterTerm = this.props.filterTerm.toUpperCase();
      let relevantFacilities = []

      if (report.facilities) {
        report.facilities.map(facility => {
          if (facility.toUpperCase().startsWith(filterTerm)) {
            relevantFacilities.push(facility)
          }
        })
        if (relevantFacilities.length == 0){
          relevantFacilities = report.facilities
        }
      }

      return <ListItem primaryText={report.title}
        secondaryText={`${country}\n${lat} ${lon}`}
        threeLines
        rightIcon={relevantFacilities.join(", ")}
        key={uuidv4()}
        leftAvatar={<Button floating secondary onClick={(e)=>this.props.viewButtonAction(report, e)}><FAEye/></Button>}
        onClick={() => showReport(report)} />;
    });
  }


  render() {
    return <div>
      <Paper zdepth={1} className="fillParent">
        <List>
          <Subheader primaryText="Locations" />
          {this.renderItems(this.props.reports, this.props.showReport)}
        </List>
      </Paper>
    </div>
  }
}
