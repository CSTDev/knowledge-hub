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
import {CreateRecord} from './data/api'
import * as _ from 'lodash';
import { ToastContainer, toast } from 'react-toastify';
import './react-toastify.css';



class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      filterText: "",
      reports: Reports,
      filteredReports: Reports,
      selectedReport: null,
      version: process.env.REACT_APP_VERSION ? process.env.REACT_APP_VERSION : "0.0.1",
      fields: []
    };
  }

  componentDidMount() {
    //TODO call api GetFields
    console.log("Got fields")
    const loadedFields = [
      "Invested",
      "Duration",
      "Start Date",
      "End Date",
      "Company"
    ]
    this.setState({fields: loadedFields})
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
      if (!reportToSave.id){
                   
        CreateRecord(reportToSave).then((response) => {
          console.dir(response)
          if(!response.ok){
            toast("Failed to Save")
            return
          }
          const id = response.json().then(json => {
            return json.ID
          }).then(id =>{
            console.log(id)
            reportToSave.id = id;

            let newReports = this.state.reports;
            newReports.push(reportToSave)
            this.setState({reports: newReports})
            console.dir(this.state)

            // const reports = _.cloneDeep(this.state.reports);
            // const reportIndex = _.findIndex(this.state.reports, (report) => {
            //   return report.title === reportToSave.title;
            // });
            // console.log
            // if (reportIndex !== -1) {
            //   reports[reportIndex].id = id
            // }
            
            // const filteredReports = FilterData(reports, this.state.filterText, null);
    
            // this.setState({reports, filteredReports});
            // console.log("State after save:")
            // console.dir(this.state)
          });
        })
        
      }
    
    }

    this.setState({
      selectedReport: null
    });
    
  };

  render() {
    return (
      <div>
        <ToastContainer />
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
        <ReportDialog report={this.state.selectedReport} onHide={this.hideReport} fields={this.state.fields} showFields={this.state.selectedReport ? !this.state.selectedReport.id : false}/>
        
      </div>
    );
  }
}

export default App;
