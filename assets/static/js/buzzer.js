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

function repopulateTable(activeParties) {
  $('#waitlist-table tbody').remove();
  $('#waitlist-table').append('<tbody>');
  for (var i in activeParties) {
    htmlStr = "<tr activePartyID="+ activeParties[i].ID + ">";
    htmlStr += "<td>" + activeParties[i].PartyName + "</td>";
    htmlStr += "<td>" + activeParties[i].PartySize + "</td>";
    htmlStr += "<td>" + parseTimeCreated(activeParties[i].TimeCreated) + "</td>";
    if (activeParties[i].PhoneAhead) {
      htmlStr += "<td><span class=\"glyphicon glyphicon-earphone\"></span></td>";
    } else {
      htmlStr += "<td><span class=\"glyphicon glyphicon-user\"></span></td>";
    }
    htmlStr += '<td><button class="btn btn-default buzz-button" type="button">Buzz!</button><button class="btn btn-default delete-party-button" type="button">Delete</button></td>';
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

$(document).ready(function() {
  $(".add-party-button").click(function(){
    // activePartyID = $('#party-name-field').id();
    partyName = $('#party-name-field').val();
    partySize = $('.btn#party-dropdown-button').val();
    waitHours = $('.btn#hours-dropdown').val();
    waitMins = $('.btn#minutes-dropdown').val();
    phoneAhead = $('.phone-ahead-toggle .active input').attr('id') === "phone" ? true : false;
    if (partyName === "") {
        $('#alert_placeholder').html('<div class="alert alert-danger alert_place" role="alert">Missing party name</div>');
      return;
    } else if (partySize === "") {
        $('#alert_placeholder').html('<div class="alert alert-danger alert_place" role="alert">Missing party size</div>');
      return;
    } else if (waitHours === "") {
        $('#alert_placeholder').html('<div class="alert alert-danger alert_place" role="alert">Missing wait time hours</div>');
      return;
    } else if (waitMins === "") {
        $('#alert_placeholder').html('<div class="alert alert-danger alert_place" role="alert">Missing wait time minutes</div>');
      return;
    }
    $('#alert_placeholder').html('');
    waitTimeExpected = parseInt(waitHours)*60 + parseInt(waitMins);
    jsonStr = JSON.stringify({"party_name": partyName, "party_size": parseInt(partySize), "wait_time_expected": waitTimeExpected, "phone_ahead": phoneAhead});
    AjaxJSONPOST("/frontend_api/create_new_party", jsonStr, addPartyErrorCallback, repopulateWaitlistSuccessCallback, completeCallback);
  });

  registerDeletePartyClickHandlers();
  registerBuzzClickHandlers();

  $(".dropdown li a").click(function(){
    console.log("in handler");
    $(this).parents(".dropdown").find('.btn').html($(this).text() + ' <span class="caret"></span>');
    $(this).parents(".dropdown").find('.btn').val($(this).text());
  });
});
