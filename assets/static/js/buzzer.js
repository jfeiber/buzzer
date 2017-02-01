console.log("sup");

// POST payload
function AjaxJSONPOST(url, jsonStr, errorCallback, successCallback, completeCallback) {
  $.ajax({
    url: url,
    type: "POST",
    data: jsonStr,
    contentType: "application/json",
    error: errorCallback,
    success: successCallback,
    complete: completeCallback
  });
}

// creates Bootstrap alert for input errors
function errorAlert(errorStr) {
  $('#alert_placeholder').html('<div class="alert alert-danger alert_place" role="alert">'+errorStr+'</div>');
}

function getAnalyticsChartErrorCallback(xhr, error) {
  console.debug(xhr);
  console.debug(error);
  errorAlert("Chart request failed.");
  $('.datepicker-spinner').hide();
}

// error callback for add party failure
function addPartyErrorCallback(xhr, error) {
  console.debug(xhr);
  console.debug(error);
  errorAlert("Add party request failed");
}

// error callback for delete party failure
function deletePartyErrorCallback(xhr, error) {
  console.debug(xhr);
  console.debug(error);
  errorAlert("Delete party request failed");
}

// error callback for buzz party failure
function buzzPartyErrorCallback(xhr, error) {
  console.debug(xhr);
  console.debug(error);
  errorAlert("Buzz party request failed");
}

// error callback for check buzzer assignment failure
function isPartyAssignedBuzzerErrorCallback(xhr, error) {
  $('#buzzer-party-modal').modal('hide');
  console.debug(xhr);
  console.debug(error);
  errorAlert("Could not check to see if buzzer was assigned party.");
}

// error callback for unlink buzzer failure
function unlinkBuzzerErrorCallback(xhr, error) {
  console.debug(xhr);
  console.debug(error);
  errorAlert("Unlink buzzer request failed.");
}

// clear buzzer assignment modal
function clearModalCallback() {
  $('#buzzer-party-modal').modal('hide');
  $('.spinner_buzzer_modal').show();
  $('#buzzer-modal-success-message').hide();
}

// check if party with given party ID is assigned a buzzer
function checkIfPartyAssignedBuzzer(activePartyID) {
  jsonObj = JSON.stringify({"active_party_id": parseInt(activePartyID)});
  AjaxJSONPOST("/frontend_api/is_party_assigned_buzzer", jsonObj, isPartyAssignedBuzzerErrorCallback, isPartyAssignedBuzzerSuccessCallback, completeCallback);
}

// success callback for buzzer assignment check
function isPartyAssignedBuzzerSuccessCallback(xhr, success) {
  console.log(xhr);
  if (xhr.is_party_assigned_buzzer) {
    refreshWaitlistTable();
    $('.spinner_buzzer_modal').hide();
    $('#buzzer-modal-success-message').show();
    setTimeout(clearModalCallback, 2000);
  } else {
    setTimeout(checkIfPartyAssignedBuzzer, 2000, xhr.active_party_id);
  }
}

// success callback for add party
function addPartySuccessCallbackBuzzer(xhr, success) {
  $('#buzzer-party-modal').modal({backdrop: 'static', keyboard: false});
  checkIfPartyAssignedBuzzer(xhr.active_party_id);
}

// success callback logging for add party
function addPartySuccessCallbackPA(xhr, success) {
  console.log(xhr);
  console.log(success);
  refreshWaitlistTable();
}

// success callback logging for waitlist population
function repopulateWaitlistSuccessCallback(xhr, success) {
  console.debug(xhr);
  console.debug(success);
  AjaxJSONPOST("/frontend_api/get_active_parties", "", addPartyErrorCallback, updateWaitlistSuccessCallback, completeCallback);
}

// log callback to the console
function completeCallback(xhr, data) {
  console.log(data);
}

// parse elapsed time into hours and minutes
function parseTimeCreated(timeCreated) {
  var timeCreatedDate = new Date(timeCreated);
  var elapsedTime = Date.now()-timeCreatedDate;
  var hours = Math.floor(elapsedTime/3600000);
  var min = Math.floor( (elapsedTime-(hours*3600000))/60000 );
  if (min < 10) {
    min = "0" + min;
  }
  if (hours < 10) {
    hours = "0" + hours;
  }
  return hours + ":" + min;
}

// parse estimated wait time into hours and minutes
function parseEstimatedWait(estimatedWaitTime) {
  var hours = Math.floor(estimatedWaitTime/60);
  var minutes = estimatedWaitTime-hours*60;
  hours = hours < 10 ? '0' + hours : hours;
  minutes = minutes < 10 ? '0' + minutes : minutes;
  return hours + ":" + minutes;
}

// refresh waitlist table every 30 seconds
function refreshWaitlistTableRepeat() {
  AjaxJSONPOST("/frontend_api/get_active_parties", "", addPartyErrorCallback, updateWaitlistSuccessCallback, completeCallback);
  setTimeout(refreshWaitlistTable, 30000);
}

// refresh waitlist table (no built-in timeout)
function refreshWaitlistTable() {
  AjaxJSONPOST("/frontend_api/get_active_parties", "", addPartyErrorCallback, updateWaitlistSuccessCallback, completeCallback);
}

// repopulate waitlist table. This method is so jank it's crazy
function repopulateTable(activeParties) {
  $('#waitlist-table tbody').remove();
  $('#waitlist-table').append('<tbody>');
  for (var i in activeParties) {
    htmlStr = "<tr activePartyID="+ activeParties[i].ID + ">";
    htmlStr += "<td>" + activeParties[i].PartyName + "</td>";
    htmlStr += "<td>" + activeParties[i].PartySize + "</td>";
    htmlStr += "<td>" + parseTimeCreated(activeParties[i].TimeCreated) + "</td>";
    htmlStr += "<td>" + parseEstimatedWait(activeParties[i].WaitTimeExpected) + "</td>";
    if (activeParties[i].PhoneAhead) {
      htmlStr += "<td><span class=\"glyphicon glyphicon-earphone\"></span></td>";
      htmlStr += '<td><div class="btn-toolbar"><button class="btn btn-default assign-buzzer-button" type="button">Assign Buzzer</button><button class="btn btn-default seat-party-button" type="button">Seat Party</button><button class="btn btn-default delete-party-button" type="button">Delete</button></div></td>';
    } else {
      htmlStr += "<td><span class=\"glyphicon glyphicon-user\"></span></td>";
      htmlStr += '<td><div class="btn-toolbar">';
      if(activeParties[i].IsTableReady) {
        htmlStr += '<button class="btn btn-default buzz-button" disabled="disabled" type="button">Buzz!</button>';
      } else {
        if (activeParties[i].BuzzerID !== 0){
          htmlStr += '<button class="btn btn-default buzz-button" type="button">Buzz!</button>';
        } else {
          htmlStr += '<button class="btn btn-default assign-buzzer-button" type="button">Assign Buzzer</button>';
        }
        htmlStr += '<button class="btn btn-default seat-party-button" type="button">Seat Party</button><button class="btn btn-default delete-party-button" type="button">Delete</button>';
      }
      htmlStr += "</div></td>";
    }
    htmlStr += "</tr>";
    $('#waitlist-table').append(htmlStr);
  }
  $('#waitlist-table').append('</tbody>');
  registerDeletePartyClickHandlers();
  registerSeatPartyClickHandlers();
  registerAssignBuzzerClickHandlers();
  registerBuzzClickHandlers();
}

// success callback for waitlist update
function updateWaitlistSuccessCallback(xhr, data) {
  repopulateTable(xhr.waitlist_data);
}

// register click handlers for deleting a party
function registerDeletePartyClickHandlers() {
  $(".delete-party-button").click(function(){
    activePartyID = $(this).closest('tr').attr('activePartyID');
    AjaxJSONPOST('/frontend_api/delete_party', JSON.stringify({"active_party_id": activePartyID, "was_party_seated" : false}), deletePartyErrorCallback, repopulateWaitlistSuccessCallback, completeCallback);
  });
}

function registerSeatPartyClickHandlers() {
  $(".seat-party-button").click(function(){
    activePartyID = $(this).closest('tr').attr('activePartyID');
    AjaxJSONPOST('/frontend_api/delete_party', JSON.stringify({"active_party_id": activePartyID, "was_party_seated": true}), deletePartyErrorCallback, repopulateWaitlistSuccessCallback, completeCallback);
  });
}

// register click handlers for buzz button
function registerBuzzClickHandlers() {
  $(".buzz-button").click(function(){
    activePartyID = $(this).closest('tr').attr('activePartyID');
    $(this).attr('disabled', 'disabled');
    AjaxJSONPOST('/frontend_api/activate_buzzer', JSON.stringify({"active_party_id": activePartyID}), buzzPartyErrorCallback, completeCallback, completeCallback);
  });
}

// register click handlers for asign buzzer button
function registerAssignBuzzerClickHandlers() {
  $(".assign-buzzer-button").click(function(){
    console.log($(this).closest('tr').attr('activePartyID'));
    activePartyID = $(this).closest('tr').attr('activePartyID');
    AjaxJSONPOST('/frontend_api/update_phone_ahead_status', JSON.stringify({"active_party_id": activePartyID}), buzzPartyErrorCallback, addPartySuccessCallbackBuzzer, completeCallback);

  });
}

// register click handlers for unlink buzzer button
function registerUnlinkBuzzerClickHandlers() {
  $(".unlink-buzzer-button").click(function(){
    buzzerID = $(this).closest('tr').attr('buzzerID');
    AjaxJSONPOST('/frontend_api/unlink_buzzer', JSON.stringify({"buzzer_id": buzzerID}), unlinkBuzzerErrorCallback, completeCallback, completeCallback);
  });
}

// reset add party fields after ADD button is hit
function resetAddPartyFields() {
  // party name
  $('#party-name-field').html('Party Name');
  $('#party-name-field').val(null);

  // party size
  $('.btn#party-dropdown-button').html('Party Size ' + '<span class="caret"></span>');
  $('.btn#party-dropdown-button').val(null);

  // wait time in minutes
  $('.btn#minutes-dropdown').html('Minutes ' + '<span class="caret"></span>');
  $('.btn#minutes-dropdown').val(null);
}

function checkIfAddPartyFormComplete() {
  partyName = $('#party-name-field').val();
  partySize = $('.btn#party-dropdown-button').val();
  waitMins = $('.btn#minutes-dropdown').val();
  if (partyName !== "" && partySize !== "" && waitMins !== "") {
    $('.add-party-button').removeAttr('disabled');
  } else {
    $('.add-party-button').attr('disabled', 'disabled');
  }
}

// Registers click/type handlers for fields/dropdowns relating to the add party menu.
function registerAddPartyHandlers() {
  // set dropdown button value and text to reflect selected value
  $(".dropdown li a").click(function(){
    console.log("in handler");
    $(this).parents(".dropdown").find('.btn').html($(this).text() + ' <span class="caret"></span>');
    $(this).parents(".dropdown").find('.btn').val($(this).text());
    checkIfAddPartyFormComplete();
  });

  $( "#party-name-field" ).keyup(function() {
    checkIfAddPartyFormComplete();
  });

  // add party click handler
  $(".add-party-button").click(function(){
    // activePartyID = $('#party-name-field').id();
    partyName = $('#party-name-field').val();
    partySize = $('.btn#party-dropdown-button').val();
    waitMins = $('.btn#minutes-dropdown').val();
    phoneAhead = $('.phone-ahead-toggle .active input').attr('id') === "phone" ? true : false;
    $('#alert_placeholder').html('');
    waitTimeExpected = parseInt(waitMins);
    jsonStr = JSON.stringify({"party_name": partyName, "party_size": parseInt(partySize), "wait_time_expected": waitTimeExpected, "phone_ahead": phoneAhead});
    successCallback = (phoneAhead) ? addPartySuccessCallbackPA : addPartySuccessCallbackBuzzer;
    AjaxJSONPOST("/frontend_api/create_new_party", jsonStr, addPartyErrorCallback, successCallback, completeCallback);
    resetAddPartyFields();
  });
}

// Based on the selected chart type, this method calls the appropriate API endpoint and updates the
// chart.
function updateAnalyicsChartWithSelection(chartType) {
  $('.datepicker-spinner').show();
  jsonObj = JSON.stringify({"start_date": $(".form-control.startDate").val(), "end_date": $(".form-control.endDate").val()});
  if (chartType === "Avg Party Size"){
    AjaxJSONPOST("/analytics_api/get_average_party_chart", jsonObj, getAnalyticsChartErrorCallback, getAveragePartySizeChartSuccessCallback, completeCallback);
  } else if (chartType === "Total Customers") {
    AjaxJSONPOST("/analytics_api/get_total_customers_chart", jsonObj, getAnalyticsChartErrorCallback, getTotalCustomersChartSuccessCallback, completeCallback);
  } else if (chartType === "Parties Per Hour") {
    AjaxJSONPOST("/analytics_api/get_parties_hour_chart", jsonObj, getAnalyticsChartErrorCallback, getPartiesPerHourChartSuccessCallback, completeCallback);
  }
}

// Checks if all the elements that are needed to select a chart have been filled out. That would be
// the chart type and the date range. If all the elements have been filled out, then the chart is
// updated.
function checkIfChartSelectionComplete() {
  chartType = $('.btn#chart-type-dropdown').val();
  startDate = $('.form-control.startDate').val();
  endDate = $('.form-control.endDate').val();
  if (chartType !== "" && startDate !== "Start Date" && endDate !== "End Date") {
    updateAnalyicsChartWithSelection(chartType);
  }
}

// Registers click handlers for the elements associated with selecting a chart.
function registerChartTypeSelectionHandler() {
  $(".chart-type-dropdown li a").click(function(){
    $(this).parents(".chart-type-dropdown").find('.btn').html($(this).text() + ' <span class="caret"></span>');
    $(this).parents(".chart-type-dropdown").find('.btn').val($(this).text());
    checkIfChartSelectionComplete();
  });

  $('.startDate').change(function(){
    checkIfChartSelectionComplete();
  });

  $('.endDate').change(function(){
    checkIfChartSelectionComplete();
  });

}

// get party info when ADD button is selected
$(document).ready(function() {

  registerDeletePartyClickHandlers();
  registerSeatPartyClickHandlers();
  registerBuzzClickHandlers();
  registerAssignBuzzerClickHandlers();
  registerUnlinkBuzzerClickHandlers();
  registerAddPartyHandlers();
  registerChartTypeSelectionHandler();

  // spinner_buzzer_modal parameters
  var opts = {
    lines: 15, // The number of lines to draw
    length: 56, // The length of each line
    width: 14, // The line thickness
    radius: 72, // The radius of the inner circle
    scale: 0.50, // Scales overall size of the spinner_buzzer_modal
    corners: 1, // Corner roundness (0..1)
    color: '#9B9B9B', // #rgb or #rrggbb or array of colors
    opacity: 0, // Opacity of the lines
    rotate: 0, // The rotation offset
    direction: 1, // 1: clockwise, -1: counterclockwise
    speed: 1, // Rounds per second
    trail: 56, // Afterglow percentage
    fps: 20, // Frames per second when using setTimeout() as a fallback for CSS
    zIndex: 2e9, // The z-index (defaults to 2000000000)
    className: 'spinner_buzzer_modal', // The CSS class to assign to the spinner_buzzer_modal
    top: '50%', // Top position relative to parent
    left: '50%', // Left position relative to parent
    shadow: false, // Whether to render a shadow
    hwaccel: false, // Whether to use hardware acceleration
    position: 'absolute', // Element positioning
  };
  var target = document.getElementById('buzzer-modal');
  var spinner_buzzer_modal = new Spinner(opts).spin(target);

  opts = {
    lines: 11, // The number of lines to draw
    length: 34, // The length of each line
    width: 6, // The line thickness
    radius: 28, // The radius of the inner circle
    scale: 0.22, // Scales overall size of the spinner
    corners: 1, // Corner roundness (0..1)
    color: '#000', // #rgb or #rrggbb or array of colors
    opacity: 0.25, // Opacity of the lines
    rotate: 0, // The rotation offset
    direction: 1, // 1: clockwise, -1: counterclockwise
    speed: 1, // Rounds per second
    trail: 60, // Afterglow percentage
    fps: 20, // Frames per second when using setTimeout() as a fallback for CSS
    zIndex: 2e9, // The z-index (defaults to 2000000000)
    className: 'datepicker-spinner', // The CSS class to assign to the spinner
    shadow: false, // Whether to render a shadow
    hwaccel: false, // Whether to use hardware acceleration
    position: 'relative' // Element positioning
  };
  target = document.getElementById('datepicker-spinner');
  var spinner_datepicker = new Spinner(opts).spin(target);
  $('.datepicker-spinner').hide();

  setTimeout(refreshWaitlistTableRepeat, 2000);
});

//  ANALYTICS STUFF   ************************
$(document).ready(function() {
  var ctx = document.getElementById("analyticsLineChart");
  var data = {
        labels: [],
        datasets: [{
            label: '',
            data: [],
            borderWidth: 1
        }]
    };
    var options =  {
      legend: {
        display: false
      }
    };
    var analyticsLineChart = Chart.Line(ctx, {
      data:data,
      options:options
    });
  });

function getAveragePartySizeChartSuccessCallback(xhr, success) {
  console.log(xhr);
  updateAnalyticsChart(xhr.graph_data, xhr.label_data, 'Average Party Size by Date', '', 'Date', 'Avg. Customers in Party');
  $('.datepicker-spinner').hide();
}

function getTotalCustomersChartSuccessCallback(xhr, success) {
  console.log(xhr);
  updateTotalCustChart(xhr.date_data, xhr.breakfast_data, xhr.lunch_data, xhr.dinner_data);
  $('.datepicker-spinner').hide();
}

function getPartiesPerHourChartSuccessCallback(xhr, success) {
  console.log(xhr);
  updateAnalyticsChart(xhr.graph_data, xhr.label_data, 'Average Parties Per Hour', '', 'Time Party Arrived', 'Avg. Number of Parties');
  $('.datepicker-spinner').hide();
}

function updateAnalyticsChart(graphData, labelData, titleString, labelString, xAxisString, yAxisString) {
    var ctx = document.getElementById("analyticsLineChart");
    var data = {
          labels: labelData,
          datasets: [{
              label: labelString,
              data: graphData,
              backgroundColor: [
                  'rgba(66, 107, 231, 0.2)',
              ],
              borderColor: [
                  'rgba(66, 107, 231, 1)',
              ],
              borderWidth: 1
          }]
      };

      var options = {
            title: {
                display: true,
                text: titleString,
                fontSize:25,
                fontColor:'#000000',
                fontFamily: 'Lato',
                fontStyle: 'oblique'
            },
            legend: {
              display: false
            },
            scales: {

              yAxes: [{
                scaleLabel: {
                  display: true,
                  labelString: yAxisString,
                  fontSize:16,
                  fontColor:'#000000',
                  fontFamily: 'Lato'
                }
              }],
              xAxes: [{
                scaleLabel: {
                  display: true,
                  labelString: xAxisString,
                  fontSize:16,
                  fontColor:'#000000',
                  fontFamily: 'Lato'
                }
              }]
            }
          };

      var analyticsLineChart = Chart.Line(ctx, {
        data:data,
        options:options
      });
}


function updateTotalCustChart(dateData, breakfastData, lunchData, dinnerData) {
    var ctx = document.getElementById("analyticsLineChart");
    var data = {
          labels: dateData,
          datasets: [{
              label: 'Breakfast',
              data: breakfastData,
              backgroundColor: [
                  'rgba(66, 107, 231, 0.2)',
              ],
              borderColor: [
                  'rgba(66, 107, 231, 1)',
              ],
              borderWidth: 1
          },
          {
              label: 'Lunch',
              data: lunchData,
              backgroundColor: [
                  'rgba(75, 192, 192, 0.2)',
              ],
              borderColor: [
                  'rgba(75, 192, 192, 1)',
              ],
              borderWidth: 1
          },
          {
              label: 'Dinner',
              data: dinnerData,
              backgroundColor: [
                  'rgba(255, 159, 64, 0.2)'
              ],
              borderColor: [
                  'rgba(255, 159, 64, 1)'
              ],
              borderWidth: 1
          }]
      };

      var options = {
            title: {
                display: true,
                text: 'Total Customers by Date',
                fontSize:25,
                fontColor:'#000000',
                fontFamily: 'Lato',
                fontStyle: 'oblique'
            },
            scales: {
              yAxes: [{
                scaleLabel: {
                  display: true,
                  labelString: 'Number of Customers',
                  fontSize:16,
                  fontColor:'#000000',
                  fontFamily: 'Lato'
                }
              }],
              xAxes: [{
                scaleLabel: {
                  display: true,
                  labelString: 'Date of Visit',
                  fontSize:16,
                  fontColor:'#000000',
                  fontFamily: 'Lato'
                }
              }]
            }
        };

      var analyticsLineChart = Chart.Line(ctx, {
        data:data,
        options:options
      });
}
