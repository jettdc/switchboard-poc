$('#all-the-buttons button').click(function (e) {
    if ($(this).is(':disabled')) e.preventDefault();
    else letterPress($(this));
});

var greyOutButton = function (button) {
   button.prop('disabled', true);
}