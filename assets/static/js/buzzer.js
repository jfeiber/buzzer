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

function addPartyErrorCallback(xhr, error) {
  console.debug(xhr);
  console.debug(error);
  $('#alert_placeholder').html('<div class="alert alert-danger alert_place" role="alert">Add party request failed</div>');
}

function addPartySuccessCallback(xhr, success) {
  console.debug(xhr);
  console.debug(success);
  AjaxJSONPOST("/frontend_api/get_active_parties", "", addPartyErrorCallback, updateWaitlistSuccessCallback, addPartyCompleteCallback);
}

function addPartyCompleteCallback(xhr, data) {
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
  for (i in activeParties) {
    htmlStr = "<tr>";
    htmlStr += "<td>" + activeParties[i].PartyName + "</td>";
    htmlStr += "<td>" + activeParties[i].PartySize + "</td>";
    htmlStr += "<td>" + parseTimeCreated(activeParties[i].TimeCreated) + "</td>";
    if (activeParties[i].PhoneAhead) {
      htmlStr += "<td><span class=\"glyphicon glyphicon-earphone\"></span></td>";
    } else {
      htmlStr += "<td><span class=\"glyphicon glyphicon-user\"></span></td>";
    }
    htmlStr += "<td><button type=\"button\" onclick=\"alert('Something!')\">Click Me!</button><button type=\"button\" onclick=\"alert('Delete!')\">Delete</button></td>"
    htmlStr += "</tr>";
    $('#waitlist-table').append(htmlStr);
  }
  $('#waitlist-table').append('</tbody>');
}

function updateWaitlistSuccessCallback(xhr, data) {
  repopulateTable(xhr["waitlist_data"]);
}

$(document).ready(function() {
  $(".add-party-button").click(function(){
    // activePartyID = $('#party-name-field').id();
    partyName = $('#party-name-field').val();
    partySize = $('.btn#party-dropdown').val();
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
    waitTimeExpected = parseInt(waitHours)*60 + parseInt(waitMins);
    jsonStr = JSON.stringify({"party_name": partyName, "party_size": parseInt(partySize), "wait_time_expected": waitTimeExpected, "phone_ahead": phoneAhead});
    AjaxJSONPOST("/frontend_api/create_new_party", jsonStr, addPartyErrorCallback, addPartySuccessCallback, addPartyCompleteCallback);
  });

  $(".delete-party-button").click(function(){
    console.log($(this).closest('tr').attr('activePartyID'));
    activePartyID = $(this).closest('tr').attr('activePartyID');
    AjaxJSONPOST('/frontend_api/delete_party', JSON.stringify({"activePartyID": activePartyID}), addPartyErrorCallback, addPartySuccessCallback, addPartyCompleteCallback);
  });

  $(".dropdown-menu li a").click(function(){
    console.log("in handler");
    $(this).parents(".dropdown").find('.btn').html($(this).text() + ' <span class="caret"></span>');
    $(this).parents(".dropdown").find('.btn').val($(this).text());
  });
});
