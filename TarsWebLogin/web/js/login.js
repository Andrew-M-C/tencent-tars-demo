// login.js
document.write("<script type='text/javascript' src='/js/sha256.js'></script>");
document.write("<script type='text/javascript' src='/js/tools.js'></script>");

$(document).ready(function(){

    let btn_submit = $("#btn_submit");
    let input_user = $("#input_user");
    let input_pass = $("#input_pass");

    // read cookie
    {
        let uid = $.cookie('uid')
        // console.info("cookie: ", uid)
        if (uid && uid != "") {
            input_user.val(uid);
        }
    }

    // Enter key
    document.onkeydown = function (e) {
        let theEvent = window.event || e;
        let code = theEvent.keyCode || theEvent.which;
        if (code == 13) {
            $("#btn_submit").click();
        }
    }

    // submit button
    $("#btn_submit").click(function() {

        let uid = input_user.val();
        let pass = input_pass.val();

        if (uid == "" || pass == "") {
            return
        } else {
            $(this).prop('disabled', true);
            input_user.prop('disabled', true);
            input_pass.prop('disabled', true);

            let session = sha256(uid + '_' + pass);
            let param = {
                uid: uid,
                hash: session,
            };

            $.ajax({
                type: "POST",
                async: true,
                dataType: "json",
                url: "/cgi-bin/login",
                contentType: 'application/json;charset=UTF-8',
                data:JSON.stringify(param),
                success: function (msg) {
                    btn_submit.prop('disabled', false);
                    input_user.prop('disabled', false);
                    input_pass.prop('disabled', false);

                    if (msg.code != 0) {
                        alert(msg.msg)
                    } else {
                        redirectUrl = getQueryVariable("redirect_url");
                        if (redirectUrl) {
                            // console.log("redirect_url: " + redirectUrl)
                            redirectUrl = redirectUrl.replace("http://", "https://")
                            window.location.href = redirectUrl;
                        } else {
                            window.location.href="/html/success.html";
                        }
                    }
                },
                error:function(XMLHttpRequest, textStatus){
                    console.log(XMLHttpRequest);  //XMLHttpRequest.responseText    XMLHttpRequest.status   XMLHttpRequest.readyState
                    console.log(textStatus);
                    alert("request error");
                }
            });
        }
    });

}); // jQuery done
