import React from 'react'
import {Avatar, Paper, List, ListItem, Subheader} from 'react-md';
import {formatLatitude, formatLongitude} from 'latlon-formatter';

const uuidv4 = require('uuid/v4');

export class DetailsPane extends React.Component {
  constructor(props) {
    super(props)
  }

  renderItems(reports, showReport) {
    return reports.map(report => {
      const lat = formatLatitude(report.location.lat, {degrees: true});
      const lon = formatLongitude(report.location.lng, {degrees: true});

      return <ListItem primaryText={report.title}
                       secondaryText={`${lat} ${lon}`}
                       key={uuidv4()}
                       leftAvatar={<Avatar suffix="deep-purple"></Avatar>}
                       onClick={() => showReport(report)}/>;
    });
  }

  render() {
    return <div>
      <Paper zdepth={1} className="fillParent">
        <List>
          <Subheader primaryText="Locations"/>
          {this.renderItems(this.props.reports, this.props.showReport)}
        </List>
      </Paper>
    </div>
  }
}
