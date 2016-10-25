console.log("sup");

$(document).ready(function() {
	$('.dropdown-menu').on('click', 'a', function(){

	  console.log("hit callback");
	  $(this).parents(".btn-group").find('.selection').text($(this).text());
	  $(this).parents(".btn-group").find('.selection').val($(this).text());

	});
});