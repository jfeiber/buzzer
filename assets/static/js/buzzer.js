console.log("sup");

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

function errorAlert(errorStr) {
  $('#alert_placeholder').html('<div class="alert alert-danger alert_place" role="alert">'+errorStr+'</div>');
}

function addPartyErrorCallback(xhr, error) {
  console.debug(xhr);
  console.debug(error);
  errorAlert("Add party request failed");
}

function deletePartyErrorCallback(xhr, error) {
  console.debug(xhr);
  console.debug(error);
  errorAlert("Delete party request failed");
}

function buzzPartyErrorCallback(xhr, error) {
  console.debug(xhr);
  console.debug(error);
  errorAlert("Buzz party request failed");
}

function isPartyAssignedBuzzerErrorCallback(xhr, error) {
  $('#buzzer-party-modal').modal('hide');
  console.debug(xhr);
  console.debug(error);
  errorAlert("Could not check to see if buzzer was assigned party.");
}

function unlinkBuzzerErrorCallback(xhr, error) {
  console.debug(xhr);
  console.debug(error);
  errorAlert("Unlink buzzer request failed.");
}

function clearModalCallback() {
  $('#buzzer-party-modal').modal('hide');
  $('.spinner').show();
  $('#buzzer-modal-success-message').hide();
}

function checkIfPartyAssignedBuzzer(activePartyID) {
  jsonObj = JSON.stringify({"active_party_id": parseInt(activePartyID)});
  AjaxJSONPOST("/frontend_api/is_party_assigned_buzzer", jsonObj, isPartyAssignedBuzzerErrorCallback, isPartyAssignedBuzzerSuccessCallback, completeCallback);
}

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

function addPartySuccessCallbackBuzzer(xhr, success) {
  $('#buzzer-party-modal').modal({backdrop: 'static', keyboard: false});
  checkIfPartyAssignedBuzzer(xhr.active_party_id);
}

function addPartySuccessCallbackPA(xhr, success) {
  console.log(xhr);
  console.log(success);
  refreshWaitlistTable();
}

function repopulateWaitlistSuccessCallback(xhr, success) {
  console.debug(xhr);
  console.debug(success);
  AjaxJSONPOST("/frontend_api/get_active_parties", "", addPartyErrorCallback, updateWaitlistSuccessCallback, completeCallback);
}

function completeCallback(xhr, data) {
  console.log(data);
}

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

function parseEstimatedWait(estimatedWaitTime) {
  var hours = Math.floor(estimatedWaitTime/60);
  var minutes = estimatedWaitTime-hours*60;
  hours = hours < 10 ? '0' + hours : hours;
  minutes = minutes < 10 ? '0' + minutes : minutes;
  return hours + ":" + minutes;
}

function refreshWaitlistTableRepeat() {
  AjaxJSONPOST("/frontend_api/get_active_parties", "", addPartyErrorCallback, updateWaitlistSuccessCallback, completeCallback);
  setTimeout(refreshWaitlistTable, 30000);
}

function refreshWaitlistTable() {
  AjaxJSONPOST("/frontend_api/get_active_parties", "", addPartyErrorCallback, updateWaitlistSuccessCallback, completeCallback);
}

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
      htmlStr += '<td><div class="btn-toolbar"><button class="btn btn-default buzz-button" type="button">Assign Buzzer</button><button class="btn btn-default delete-party-button" type="button">Delete</button></div></td>';
    } else {
      htmlStr += "<td><span class=\"glyphicon glyphicon-user\"></span></td>";
      htmlStr += '<td><div class="btn-toolbar"><button class="btn btn-default buzz-button" type="button">Buzz!</button><button class="btn btn-default delete-party-button" type="button">Delete</button></div></td>';
    }
    htmlStr += "</tr>";
    $('#waitlist-table').append(htmlStr);
  }
  $('#waitlist-table').append('</tbody>');
  registerDeletePartyClickHandlers();
  registerBuzzClickHandlers();
}

function updateWaitlistSuccessCallback(xhr, data) {
  repopulateTable(xhr.waitlist_data);
}

function registerDeletePartyClickHandlers() {
  $(".delete-party-button").click(function(){
    console.log($(this).closest('tr').attr('activePartyID'));
    activePartyID = $(this).closest('tr').attr('activePartyID');
    AjaxJSONPOST('/frontend_api/delete_party', JSON.stringify({"active_party_id": activePartyID}), deletePartyErrorCallback, repopulateWaitlistSuccessCallback, completeCallback);
  });
}

function registerBuzzClickHandlers() {
  $(".buzz-button").click(function(){
    console.log($(this).closest('tr').attr('activePartyID'));
    activePartyID = $(this).closest('tr').attr('activePartyID');
    AjaxJSONPOST('/frontend_api/activate_buzzer', JSON.stringify({"active_party_id": activePartyID}), buzzPartyErrorCallback, completeCallback, completeCallback);
  });
}

function registerUnlinkBuzzerClickHandlers() {
  $(".unlink-buzzer-button").click(function(){
    console.log($(this).closest('tr').attr('activePartyID'));
    activePartyID = $(this).closest('tr').attr('activePartyID');
    AjaxJSONPOST('/frontend_api/unlink_buzzer', JSON.stringify({"active_party_id": activePartyID}), unlinkBuzzerErrorCallback, completeCallback, completeCallback);
  });
}

function registerGetHistoricalClickHandlers() {
    $(".get_parties_button").on('click', function() {
         jsonObj = JSON.stringify({"start_date": $(".form-control.startDate").val(),
             "end_date": $(".form-control.endDate").val()
         });
         AjaxJSONPOST("/analytics_api/get_historical_parties", jsonObj, function(response) { console.log(response); }, getHistoricalPartiesSuccessCallback, completeCallback);
    });
}

$(document).ready(function() {
  $(".add-party-button").click(function(){
    // activePartyID = $('#party-name-field').id();
    partyName = $('#party-name-field').val();
    partySize = $('.btn#party-dropdown-button').val();
    waitMins = $('.btn#minutes-dropdown').val();
    phoneAhead = $('.phone-ahead-toggle .active input').attr('id') === "phone" ? true : false;
    if (partyName === "") {
        $('#alert_placeholder').html('<div class="alert alert-danger alert_place" role="alert">Missing party name</div>');
      return;
    } else if (partySize === "") {
        $('#alert_placeholder').html('<div class="alert alert-danger alert_place" role="alert">Missing party size</div>');
      return;
    } else if (waitMins === "") {
        $('#alert_placeholder').html('<div class="alert alert-danger alert_place" role="alert">Missing wait time minutes</div>');
      return;
    }
    $('#alert_placeholder').html('');
    waitTimeExpected = parseInt(waitMins);
    jsonStr = JSON.stringify({"party_name": partyName, "party_size": parseInt(partySize), "wait_time_expected": waitTimeExpected, "phone_ahead": phoneAhead});
    successCallback = (phoneAhead) ? addPartySuccessCallbackPA : addPartySuccessCallbackBuzzer;
    AjaxJSONPOST("/frontend_api/create_new_party", jsonStr, addPartyErrorCallback, successCallback, completeCallback);
  });

  registerDeletePartyClickHandlers();
  registerBuzzClickHandlers();
  registerUnlinkBuzzerClickHandlers();
  registerGetHistoricalClickHandlers();

  $(".dropdown li a").click(function(){
    console.log("in handler");
    $(this).parents(".dropdown").find('.btn').html($(this).text() + ' <span class="caret"></span>');
    $(this).parents(".dropdown").find('.btn').val($(this).text());
  });

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
    console.log("hell yeah");
    if (xhr.historical_parties) {
        console.log(xhr)
        xhr.historical_parties.forEach( function (party) {
            $("#historical_parties").append("partyName:\t" + party.PartyName + "\t" + "TimeSeated:\t" + party.TimeSeated);
            $("#historical_parties").append("<br>");
        });
    }
}