console.log("sup");

$(document).ready(function() {
	// $('.dropdown-menu').on('click', 'a', function(){

	//   console.log("hit callback");
	//   $(this).parents(".dropdown").find('.selection').text($(this).text());
	//   $(this).parents(".dropdown").find('.selection').val($(this).text());

	// });

        $(".dropdown-menu li a").click(function(){
                console.log("in handler");
                $(this).parents(".dropdown").find('.btn').html($(this).text() + ' <span class="caret"></span>');
                $(this).parents(".dropdown").find('.btn').val($(this).data('value'));
        });
});