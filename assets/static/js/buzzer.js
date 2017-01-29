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
  $('.spinner').show();
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
    $('.spinner').hide();
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
    console.log("asdfhkasjhdf");
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

// placeholder until Joon comments this
function registerGetHistoricalClickHandlers() {
    $(".get_parties_button").on('click', function() {
         jsonObj = JSON.stringify({"start_date": $(".form-control.startDate").val(),
             "end_date": $(".form-control.endDate").val()
         });
         AjaxJSONPOST("/analytics_api/get_historical_parties", jsonObj, function(response) { console.log(response); }, getHistoricalPartiesSuccessCallback, completeCallback);
    });
}

// Going to delete
function registerGetHistoricalPartiesClickHandler() {
    $(".get_parties_button").on('click', function() {
         jsonObj = JSON.stringify({"start_date": $(".form-control.startDate").val(),
             "end_date": $(".form-control.endDate").val()
         });
         AjaxJSONPOST("/analytics_api/get_historical_parties", jsonObj, function(response) { console.log(response); }, getHistoricalPartiesSuccessCallback, completeCallback);
    });
}

// another placeholder until Joon comments this
function registerGetAveragePartySizeClickHandler() {
    $(".get_average_party_size_button").on('click', function() {
         jsonObj = JSON.stringify({"start_date": $(".form-control.startDate").val(),
             "end_date": $(".form-control.endDate").val()
         });
         AjaxJSONPOST("/analytics_api/get_historical_parties", jsonObj, function(response) { console.log(response); }, getHistoricalPartiesSuccessCallback, completeCallback);
         AjaxJSONPOST("/analytics_api/get_average_party_size", jsonObj, function(response) { console.log(response); }, getAveragePartySizeSuccessCallback, completeCallback);
    });
}

function registerGetAverageWaitTimeClickHandler() {
    $(".get_average_wait_time_button").on('click', function() {
         jsonObj = JSON.stringify({"start_date": $(".form-control.startDate").val(),
             "end_date": $(".form-control.endDate").val()
         });
         AjaxJSONPOST("/analytics_api/get_historical_parties", jsonObj, function(response) { console.log(response); }, getHistoricalPartiesSuccessCallback, completeCallback);
         AjaxJSONPOST("/analytics_api/get_average_wait_time", jsonObj, function(response) { console.log(response); }, getAverageWaitTimeSuccessCallback, completeCallback);
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

// get party info when ADD button is selected
$(document).ready(function() {

  registerDeletePartyClickHandlers();
  registerSeatPartyClickHandlers();
  registerBuzzClickHandlers();
  registerAssignBuzzerClickHandlers();
  registerUnlinkBuzzerClickHandlers();
  registerGetHistoricalClickHandlers();
  registerGetAveragePartySizeClickHandler();
  registerGetAverageWaitTimeClickHandler();
  registerAddPartyHandlers();

  // spinner parameters
  var opts = {
    lines: 15, // The number of lines to draw
    length: 56, // The length of each line
    width: 14, // The line thickness
    radius: 72, // The radius of the inner circle
    scale: 0.50, // Scales overall size of the spinner
    corners: 1, // Corner roundness (0..1)
    color: '#9B9B9B', // #rgb or #rrggbb or array of colors
    opacity: 0, // Opacity of the lines
    rotate: 0, // The rotation offset
    direction: 1, // 1: clockwise, -1: counterclockwise
    speed: 1, // Rounds per second
    trail: 56, // Afterglow percentage
    fps: 20, // Frames per second when using setTimeout() as a fallback for CSS
    zIndex: 2e9, // The z-index (defaults to 2000000000)
    className: 'spinner', // The CSS class to assign to the spinner
    top: '50%', // Top position relative to parent
    left: '50%', // Left position relative to parent
    shadow: false, // Whether to render a shadow
    hwaccel: false, // Whether to use hardware acceleration
    position: 'absolute', // Element positioning
  };
  var target = document.getElementById('buzzer-modal');
  var spinner = new Spinner(opts).spin(target);

  setTimeout(refreshWaitlistTableRepeat, 2000);
});

function getHistoricalPartiesSuccessCallback(xhr, success) {
    if (xhr.historical_parties) {
        for (var date in xhr.historical_parties) {
            for (var party in xhr.historical_parties[date]) {
                $("#historical_parties").append("date: " + date + "---- " + party.PartyName + "\t" + "<br>");
            }
        }
    }
}

function getAveragePartySizeSuccessCallback(xhr, success) {
    if (xhr.labels) {
        xhr.labels.forEach(function(date) {
            $("#historical_parties").append("date: " + date + "---- " + "\t" + "<br>");
        });
    }

    if (xhr.data) {
        xhr.data.forEach(function(datapiece) {
            $("#historical_parties").append("averagePartySize: " + datapiece + "---- " + "\t" + "<br>");
        });
    }
}

function getAverageWaitTimeSuccessCallback(xhr, success) {
    if ("average_wait_hours" in xhr) {
        $("#average_party_size").append("average wait hours:\t" + xhr.average_wait_hours);
        $("#average_party_size").append("<br>");
    }
    if (xhr.average_wait_minutes) {
        $("#average_party_size").append("average wait minutes:\t" + xhr.average_wait_minutes);
        $("#average_party_size").append("<br>");
    }
}
