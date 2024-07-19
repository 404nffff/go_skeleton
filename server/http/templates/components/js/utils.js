/**
 * post 请求
 * @param {*} url 
 * @param {*} data 
 * @param function success 
 * @param {*} error 
 */
function post(url, data, success) {
    $.ajax({
        type: 'POST',
        url: url,
        data: data,
        success: success,
        error: function (xhr, type) {

            layer.msg(xhr.responseJSON.msg, {icon: 2});
        }
    });
}