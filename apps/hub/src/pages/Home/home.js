import React, { Component } from 'react';
import './home.css';
import 'react-md/dist/react-md.blue_grey-amber.min.css';
import 'material-icons/iconfont/MaterialIcons-Regular.woff2'
import 'material-icons/iconfont/material-icons.css'
import { SearchBar } from '../../components/SearchBar';
import { DetailsPane } from '../../components/DetailsPane';
import { FilterData } from '../../data/Filter';
import { ReportDialog } from "../../components/ReportDialog/reportDialog";
import { MapView } from '../../components/Map/map';
import { VersionBar } from '../../components/VersionBar/versionBar';
import { CreateRecord, GetRecords, LoadFields, UpdateRecord } from '../../data/api'
import { ToastContainer, toast } from 'react-toastify';
import '../../react-toastify.css';
import MenuBar from '../../components/MenuBar/menuBar';



class Home extends Component {
  constructor(props) {
    super(props);

    this.state = {
      filterText: "",
      reports: [],
      filteredReports: [],
      selectedReport: null,
      version: process.env.REACT_APP_VERSION ? process.env.REACT_APP_VERSION : "0.0.1",
      fields: []
    };
  }

  componentDidMount() {
    LoadFields().then((response) => {
      if (!response || response.status !== 200) {

        return
      }
      response.json().then(json => {
        return json
      }).then(fieldList => {
        this.setState({ fields: fieldList })
      })

    });
  }

  filterChange(value, event) {
    const filteredReports = FilterData(this.state.reports, value, null);

    this.setState({
      filterText: value,
      filteredReports: filteredReports
    });
  }

  viewSummary() {


  }

  getRecords = (bounds) => {
    let records = GetRecords(bounds).then(response => {
      if (!response || (response.status !== 200)) {
        if (response.message && response.message == "404") {
          this.setState({ reports: [], filteredReports: [] });
          return
        }
        toast("Failed to load records");
        return;
      }
      response.json().then(json => {
        return json;
      }).then(records => {
        this.setState({ reports: records });

        this.filterChange(this.state.filterText)

      });
    });
  }

  showReport = (report) => {
    if (!report.details) {
      report.details = {}
    }
    this.setState({
      selectedReport: report
    });
  };

  hideReport = (reportToSave) => {
    if (reportToSave) {
      if (!reportToSave.id) {
        CreateRecord(reportToSave).then((response) => {
          if (!response || response.status !== 201) {
            toast("Failed to Save")
            return
          }
          response.json().then(json => {
            return json.ID
          }).then(id => {
            reportToSave.id = id;

            let newReports = this.state.reports;
            newReports.push(reportToSave)
            this.setState({ reports: newReports })
          });

          this.setState({
            selectedReport: null
          });
        });

      } else {
        UpdateRecord(reportToSave).then((response) => {
          if (!response || response.status !== 200) {
            toast("Failed to Save")
            return
          }
          toast("Saved")

          this.setState({
            selectedReport: null
          });
        });
      }

    } else {
      this.setState({
        selectedReport: null
      });
    }

  };

  render() {
    return (
      <div>
        <ToastContainer />
        <VersionBar version={this.state.version} />
        <MenuBar />
        <div className="mapArea">
          <div className="searchArea" style={{ display: this.state.selectedReport ? "none" : "block" }}>
            <SearchBar value={this.state.filterText} onChange={(value, event) => this.filterChange(value, event)} />
          </div>
          <div className="mapViewArea">
            <MapView className="mapView" reports={this.state.filteredReports} view={this.showReport} viewSummary={this.viewSummary} getRecords={this.getRecords} />
          </div>
        </div>
        <div className="detailsArea">
          <DetailsPane reports={this.state.filteredReports} showReport={this.showReport} filterTerm={this.state.filterText} />
        </div>
        <ReportDialog report={this.state.selectedReport} onHide={this.hideReport} fields={this.state.fields} showFields={this.state.selectedReport ? !this.state.selectedReport.id : false} />

      </div>
    );
  }
}

export default Home;
