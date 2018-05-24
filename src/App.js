import React, { Component } from 'react';
import './App.css';
import 'react-md/dist/react-md.blue_grey-amber.min.css';
import 'material-icons/iconfont/MaterialIcons-Regular.woff2'
import 'material-icons/iconfont/material-icons.css'
import {SearchBar} from './components/SearchBar';
import {DetailsPane} from  './components/DetailsPane';
import {Reports} from './data/Reports';
import {FilterData} from './data/Filter';
import {ReportDialog} from "./components/ReportDialog/reportDialog";
import {MapView} from './components/Map/map';
import {VersionBar} from './components/VersionBar/versionBar';
import * as _ from 'lodash';



class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      filterText: "",
      reports: Reports,
      filteredReports: Reports,
      selectedReport: null,
      version: process.env.REACT_APP_VERSION ? process.env.REACT_APP_VERSION : "0.0.1"
    };
  }

  filterChange(value, event) {
    const filteredReports = FilterData(this.state.reports, value, null);

    this.setState({
      filterText: value,
      filteredReports: filteredReports
    });
  }

  viewSummary(){
    
    
  }

  showReport = (report) => {
    this.setState({
      selectedReport: report
    });
  };

  hideReport = (reportToSave) => {
    if (reportToSave) {
      const reports = _.cloneDeep(this.state.reports);
      const reportIndex = _.findIndex(reports, (report) => {
        return report.title === reportToSave.title;
      });
      if (reportIndex !== -1) {
        reports[reportIndex] = reportToSave;
      }
      const filteredReports = FilterData(reports, this.state.filterText, null);

      this.setState({reports, filteredReports});
    }

    this.setState({
      selectedReport: null
    });
  };

  render() {
    return (
      <div>
        <VersionBar version={this.state.version}/>
        <div className="mapArea">
          <div className="searchArea" style={{display: this.state.selectedReport ? "none" : "block"}}>
            <SearchBar value={this.state.filterText} onChange={(value, event) => this.filterChange(value, event)}/>
          </div>
          <div className="mapViewArea">
            <MapView className="mapView" reports={this.state.filteredReports} view={this.showReport} viewSummary={this.viewSummary}/>
          </div>
        </div>
        <div className="detailsArea">
          <DetailsPane reports={this.state.filteredReports} showReport={this.showReport}/>
        </div>
        <ReportDialog report={this.state.selectedReport} onHide={this.hideReport}/>
      </div>
    );
  }
}

export default App;
