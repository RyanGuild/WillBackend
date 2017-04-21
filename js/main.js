/*global $*/
function removeClassByJquery (selector, className) {
    $(selector).removeClass(className);
}
function addClassByJquery (selector, className) {
    $(selector).addClass(className);
}
function hasClassByJquery (selctor, className) {
    return $(selctor).hasClass(className);
}