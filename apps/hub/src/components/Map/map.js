import React from 'react'
import { Marker, Map, Popup, TileLayer} from 'react-leaflet';
import 'leaflet/dist/leaflet.css'
import MarkerClusterGroup from 'react-leaflet-markercluster';
import './map.css'
import * as _ from 'lodash';
import { Button } from 'react-md';


// Webpack/leaflet fix
// See https://github.com/Leaflet/Leaflet/issues/4968
import L from 'leaflet';
_.once(() => {
  delete L.Icon.Default.prototype._getIconUrl;

  L.Icon.Default.mergeOptions({
    iconRetinaUrl: require('leaflet/dist/images/marker-icon-2x.png'),
    iconUrl: require('leaflet/dist/images/marker-icon.png'),
    shadowUrl: require('leaflet/dist/images/marker-shadow.png'),
  });
})();


export class MapView extends React.Component {
   constructor (props) {
    super(props);
    this.state = {
      lat: 53.8008,
      lng: -1.5491,
      zoom: 11
    }
  }

  newPoint = (e) => {
    var newReport = {
      "title": "New",
      "location": {
        "lng": e.latlng.lng,
        "lat": e.latlng.lat
      },
    }
    this.props.view(newReport)
  }

  render() {
    const position = [this.state.lat, this.state.lng]

    return (
      <div className="mapContainer">
          <Map center={position} zoom={this.state.zoom} oncontextmenu={this.newPoint}>
            <TileLayer
              url='http://{s}.tile.osm.org/{z}/{x}/{y}.png'
              attribution='&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
            />
            <MarkerGroup reports={this.props.reports} view={this.props.view} viewSummary={this.props.viewSummary}/>
          </Map>
      </div>
    );
  }
}

function MarkerGroup(props){
  const reports = props.reports;
  const markers = reports.map((report, index)=>{
    return (<Marker position={report.location} key={index}>
        <Popup >
        <div className="summary">
          <h4>{report.title}</h4>
            Description: {report.details ? report.details.description :""}<br/>
            <div class="buttonHolder">
              <Button raised onClick={()=>props.view(report)}>View</Button>
            </div>
          </div>     
        </Popup>
      </Marker>)
  }
  );
  return (
    <MarkerClusterGroup zoomToBoundsOnClick={false} onClusterClick={()=>props.viewSummary()}>{markers}</MarkerClusterGroup>
  )
}
