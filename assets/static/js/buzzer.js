console.log("sup");

$(document).ready(function() {
	// $('.dropdown-menu').on('click', 'a', function(){

	//   console.log("hit callback");
	//   $(this).parents(".dropdown").find('.selection').text($(this).text());
	//   $(this).parents(".dropdown").find('.selection').val($(this).text());

	// });

  $(".add-party-button").click(function(){
    // console.log("add party button handler.");
    activePartyID = $('#party-name-field').id();
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
    jsonObj = JSON.stringify({"party_name": partyName, "party_size": parseInt(partySize), "wait_time_expected": waitTimeExpected, "phone_ahead": phoneAhead});
    console.log(jsonObj);
    $.ajax({
      url: "/frontend_api/create_new_party",
      type: "POST",
      data: jsonObj,
      contentType: "application/json",
      error: function(xhr, error){
        console.debug(xhr);
        console.debug(error);
        $('#alert_placeholder').html('<div class="alert alert-danger alert_place" role="alert">Add party request failed</div>')      },
      success: function(xhr, success){
        console.debug(xhr);
        console.debug(success);
      },
      complete: function(data) {
        console.log(data);
      }
    });

    //in the future this will load this via an AJAX call. For now I am lazy.
  });

  $(".dropdown-menu li a").click(function(){
    console.log("in handler");
    $(this).parents(".dropdown").find('.btn').html($(this).text() + ' <span class="caret"></span>');
    $(this).parents(".dropdown").find('.btn').val($(this).text());
  });
});
