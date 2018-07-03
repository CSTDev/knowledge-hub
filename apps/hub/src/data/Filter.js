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
export function FilterData(reports, filterText, filterOptions, mapBounds) {
  filterText = filterText.toUpperCase();
  const reportsMatchingTextFilter = _.filter(reports, report => {

    if (filterText === "")
      return true;
    if (filterOptions.title && report.title.toUpperCase().includes(filterText)) {
      return true;
    }

    if (filterOptions.locations && report.location.country.toUpperCase().includes(filterText)) {
      return true;
    }

    if (filterOptions.shortName && report.shortName.toUpperCase().startsWith(filterText)) {
      return true;
    }

    if (filterOptions.facilities && (_.some(report.facilities, facility => {
      return facility.toUpperCase().startsWith(filterText);
    }) )) {
      return true;
    }

    return false;
  });

  return reportsMatchingTextFilter;
}

