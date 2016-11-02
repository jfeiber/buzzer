function AjaxJSONPOST(url, jsonObj, errorCallback, successCallback, completeCallback) {
  $.ajax({
    url: url,
    type: "POST",
    data: jsonObj,
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
}

function addPartyCompleteCallback(xhr, data) {
  console.log(data);
}

$(document).ready(function() {
  // AjaxJSONPOST("/frontend_api/get_active_parties", addPartyErrorCallback, addPartySuccessCallback, addPartyCompleteCallback);
  $(".add-party-button").click(function(){
    // console.log("add party button handler.");
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
    $('#alert_placeholder').html('');
    waitTimeExpected = parseInt(waitHours)*60 + parseInt(waitMins);
    jsonObj = JSON.stringify({"party_name": partyName, "party_size": parseInt(partySize), "wait_time_expected": waitTimeExpected, "phone_ahead": phoneAhead});
    console.log(jsonObj);
    AjaxJSONPOST("/frontend_api/create_new_party", jsonObj, addPartyErrorCallback, addPartySuccessCallback, addPartyCompleteCallback);

    //in the future this will load this via an AJAX call. For now I am lazy.
  });

  $(".dropdown-menu li a").click(function(){
    console.log("in handler");
    $(this).parents(".dropdown").find('.btn').html($(this).text() + ' <span class="caret"></span>');
    $(this).parents(".dropdown").find('.btn').val($(this).text());
  });
});
