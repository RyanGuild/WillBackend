var picArray;
var len = $('.profPhotoContainer').length;
console.log(len);
picArray = new Array(len);
$('.profPhotoContainer').each(function (i,e) {
    var CardNum = String(e.id).slice(9);
    CardNum = parseInt(CardNum);
    var slice = [];
    $(e).children('img').each(function (e) {
        slice.push(e);
    });
    if (slice != []){
        console.log("img show:", CardNum);
        $(slice[0]).show();
    }
    picArray[CardNum] = slice;
    $('#cardNext'+CardNum).click(nextPic);
});

function nextPic (e) {
    var cardNum = parseInt(String(this.id).slice(8));
    
}