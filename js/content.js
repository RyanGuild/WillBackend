/*global $*/
var wScrollCurrent = window.scrollY;
var wScrollDiff = 0;

function throttle(wait, func, options) {
    var context, args, result;
    var timeout = null;
    var previous = 0;
    if (!options) {options = {}; }
    var later = function () {
        previous = options.leading === false ? 0 : Date.now();
        timeout = null;
        result = func.apply(context, args);
        if (!timeout) {context = args = null; }
    };
    return function () {
        var now = Date.now();
        if (!previous && options.leading === false) {previous = now; }
        var remaining = wait - (now - previous);
        context = this;
        args = arguments;
        if (remaining <= 0 || remaining > wait) {
            if (timeout) {
                clearTimeout(timeout);
                timeout = null;
            }
            previous = now;
            result = func.apply(context, args);
            if (!timeout) {context = args = null; }
        } else if (!timeout && options.trailing !== false) {
            timeout = setTimeout(later, remaining);
        }
        return result;
    };
}
window.addEventListener('scroll', throttle(250, function () {
    wScrollDiff = wScrollCurrent - window.scrollY;
    wScrollCurrent = window.scrollY;
    if (wScrollCurrent <= 0) {
        parent.removeClassByJquery('#Header', 'header--hidden');
    } else if (wScrollDiff > 0 && parent.hasClassByJquery('#Header', 'header--hidden')) {
        parent.removeClassByJquery('#Header', 'header--hidden');
    } else if (wScrollDiff < 0) {
        if (wScrollCurrent + window.height >= document.height -200 && parent.hasClassByJquery('#Header', 'header--hidden')) {
            parent.removeClassByJquery('#Header', 'header--hidden');
        } else {
            parent.addClassByJquery('#Header', 'header--hidden');
        }
    }
}));