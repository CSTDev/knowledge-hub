import * as _ from 'lodash';

/**
 * Filters reports. If a report title or and report text contains the filterText
 * then the report passes the filter
 * @param reports
 * @param filterText
 * @param mapBounds
 * @returns {*}
 * @constructor
 */
export function FilterData(reports, filterText, mapBounds) {
  const reportsMatchingTextFilter = _.filter(reports, report => {
    if (report.title.includes(filterText)) {
      return true;
    }

    if (report.reports.some(childReport => childReport.reportDetails.includes(filterText))) {
      return true;
    }

    return false;
  });

  return reportsMatchingTextFilter;
  // TODO: Add map bounds filter
}

