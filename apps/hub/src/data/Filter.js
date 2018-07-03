import * as _ from 'lodash';

/**
 * Filters reports. If a report title, short name, or facilities contains the text
 * then the report passes the filter
 * @param reports
 * @param filterText
 * @param mapBounds
 * @returns {*}
 * @constructor
 */
export function FilterData(reports, filterText, mapBounds) {
  const reportsMatchingTextFilter = _.filter(reports, report => {


    if (report.title.includes(filterText) || report.shortName.includes(filterText) || _.some(report.facilities, facility => {
      return facility.includes(filterText);
    }))  {
      return true;
    }

    // if (report.reports.some(childReport => childReport.reportDetails.includes(filterText))) {
    //   return true;
    // }

    return false;
  });

  return reportsMatchingTextFilter;
  // TODO: Add map bounds filter
}

