// tools.js
// http://www.runoob.com/w3cnote/js-get-url-param.html

function getQueryVariable(variable)
{
    var query = window.location.search.substring(1);
    var vars = query.split("&");
    for (var i = 0; i < vars.length; i++) {
        var pair = vars[i].split("=");
        if(pair[0] == encodeURIComponent(variable)) {
            var ret = decodeURIComponent(pair[1])
            return ret;
        }
    }
    return(undefined);
}
