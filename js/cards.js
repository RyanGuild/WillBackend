/*global $*/
var index = 0;
var count = 10;
var picArray;
$(document).ready(function (e) {
   loadCards(index, count);
});

function AjaxSuc(resp){
    var jresp = JSON.parse(resp);
    console.log('AjaxSucsess');
    $('#contentBody').html('');
    picArray = [];
    for (var i = 0; i< Object.keys(jresp['Cards']).length; i++){
        picArray.push([]);
    }
    jresp['Cards'].forEach(function (e) {
        $('#contentBody').append(e['Payload']);
    });
    console.log('ContentSucsess');
    jresp['Cards'].forEach(function (e) {
        picArray[e["Photos"]['Index']] = e["Photos"]["PhotosBucket"];
    });
    console.log('PhotoSaveSucsess');
    initPics();
    console.log('PhotoPostSucsess');
}

function loadCards(index, count){
    $('.loadIMG').show(.2);
    var sendInfo = {
            index:index,
            count:count
    };
    var jinfo = JSON.stringify(sendInfo)
    $.ajax({
        url:"http://127.0.0.1:8080/getProf",
        method:"GET",
        data: jinfo,
        dataType:"text",
        cash:false,
        success: AjaxSuc,
        error: function (x,e) {
            $('.loadIMG').hide(.2);
            console.error(x.status,",",e)
            $('#contentBody').html(x.responseText);
        }
    });
}

function initPics() {
    $('.profPhotoContainer').each(function (i,e) {
        var CardNum = String(e.id).slice(9);
        CardNum = parseInt(CardNum);
        console.log(CardNum);
        if (picArray[CardNum] == []){
            e.css("background", "url(../resourses/prof/"+picArray[CardNum][0]+")");
        }
        $('#cardNext'+CardNum).click(function (e) {
            
        });
    })
}



