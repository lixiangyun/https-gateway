String.prototype.format= function(){
    //将arguments转化为数组（ES5中并非严格的数组）
    var args = Array.prototype.slice.call(arguments);
    var count=0;
    //通过正则替换%s
    return this.replace(/%s/g,function(s,i){
        return args[count++];
    });
}

/*校验邮件地址是否合法 */
function IsEmail(str) {
    var reg=/^\w+@[a-zA-Z0-9]{2,10}(?:\.[a-z]{2,4}){1,3}$/;
    return reg.test(str);
}

function IsDomain(str) {
    var reg=/[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+\.?/;
    return reg.test(str)
}

function IsNumber(abc) {
    var reg=/[0-9]{1,5}/
    return reg.test(abc)
}

function _transferTable(list, curlist) {
    var data1 = [];
    var data2 = [];
    list.forEach(function (value, idx, array) {
        data1.push({"value": idx.toString(), "title": value});
        for ( var i = 0;i < curlist.length; i++ ) {
            if ( curlist[i] == value ) {
                data2.push(idx.toString());
                break;
            }
        }
    })
    return [data1, data2];
}

function _transferListGet(data1) {
    var data2 = [];
    if (data1.length == 0 ) {
        return data2;
    }
    data1.forEach(function (value, idx, array) {
        data2.push(value.title);
    })
    return data2;
}

function publicTransferOpen( title, alllist, curlist, success ) {
    var innerHtml = document.getElementById("transfer-div").innerHTML;
    if ( innerHtml == undefined ) {
        display("not found transferdiv")
        return;
    }
    var newfilter = "transferopen";
    var newKey = 'key123';
    var newDemo = "transfer_demo_x";

    innerHtml = innerHtml.replace("transfer-div-id", newfilter);
    innerHtml = innerHtml.replace("transfer-demo", newDemo);

    data = _transferTable(alllist, curlist);
    publicOpen(
        title,
        innerHtml,
        "550px","500px",
        function () {
            var transfer = layui.transfer;
            transfer.render({
                elem: '#transferopen'
                ,data: data[0]
                ,id: newKey
                ,title: ['未绑定列表', '已绑定列表']
                ,value: data[1]
                ,showSearch: true
            })
        },
        function (data) {
            var transfer = layui.transfer;
            var getData = transfer.getData(newKey)
            success(_transferListGet(getData))
        }
    )
}

function publicOpen(title, content, wight, height, render, success ) {
    layer.open({
        type: 1
        ,title: title
        ,area: [wight, height]
        ,content: content
        ,btn: '提交'
        ,btnAlign: 'c'
        ,shade: 0
        ,success: function(layero){
            render();
        }
        ,yes: success
    });
}

function publicOpenFormWithParm(title, divid, fromfilter, ok, wight, height, success) {
    var divbody = document.getElementById(divid).innerHTML;
    var newfilter = "_openpage_form_"+ Math.random();
    var pagebody = divbody.replace(fromfilter, newfilter);

    layer.open({
        type: 1
        ,title: title
        ,area: [wight, height]
        ,content: pagebody
        ,btn: '提交'
        ,btnAlign: 'c'
        ,shade: 0
        ,success: function(layero){
            layui.form.render();
            ok(newfilter)
        }
        ,yes: function () {
            success(layui.form.val(newfilter));
        }
    });
}

function publicOpenForm(title, divid, fromfilter, wight, height, success) {
    var divbody = document.getElementById(divid).innerHTML;
    var newfilter = "_openpage_form_"+ Math.random();
    var pagebody = divbody.replace(fromfilter, newfilter);

    layer.open({
        type: 1
        ,title: title
        ,area: [wight, height]
        ,content: pagebody
        ,btn: '提交'
        ,btnAlign: 'c'
        ,shade: 0
        ,success: function(layero){
            layui.form.render();
        }
        ,yes: function () {
            success(layui.form.val(newfilter));
        }
    });
}

function publicHttpJson(method, url, req, success ) {
    layui.jquery.ajax({
        type: method,
        url: url,
        data: req,
        timeout: 30000,
        success:function(body,status){
            success(JSON.parse(body));
        },
        error: function(XMLHttpRequest, textStatus, errorThrown) {
            display("connect timeout");
        }
    });
}

function publicPostJson(url, request, func ) {
    publicHttpJson('post', url, request, func );
}

function publicUpdateJson(url, request, func ) {
    publicHttpJson('put', url, request, func );
}

function publicDelJson(url, request, func ) {
    publicHttpJson('delete', url, request, func );
}

function publicGetJson(url, func ) {
    publicGetJsonReq(url, "", func)
}

function publicGetJsonReq(url, param, func ) {
    layui.jquery.getJSON({
        url: url,
        data: param,
        timeout: 30000,
        success: function (rsp,status,xhr) {
            func(rsp)
        },
        error: function(XMLHttpRequest, textStatus, errorThrown) {
            display("connect timeout");
        }
    })
}

function publicPageReload(param) {
    var div = document.getElementById("main-body");
    layui.jquery.ajax({
        type:'get',
        url: param,//这里是baiurl
        timeout:30000,
        success:function(body,status){
            div.innerHTML = body;
            layui.jquery('#main-body').html(body);
            layer.closeAll();
        }
    });
}

function publicOptionReload(id, array, names ) {
    var opt = document.getElementById(id);
    if (opt == undefined) {
        return;
    }
    var innerHtml = "<option value=\"\"></option>";
    var formatHtml = "<option value=\"%s\">%s</option>\r\n";
    array.forEach(function(value, index, array) {
        innerHtml += formatHtml.format(value, names[index]);
    })
    opt.innerHTML = innerHtml;
    layui.use('form', function() {
        var form = layui.form; //只有执行了这一步，部分表单元素才会自动修饰成功
        form.render();
    });
}

function publicProxyListGet(id, success) {
    publicGetJson('/proxy', function (rsp) {
        var array = [];
        var names = [];
        if ( rsp.code != undefined && rsp.code == 0 ) {
            for ( i = 0; i < rsp.count; i++ ) {
                array.push(rsp.data[i].id)
                names.push(rsp.data[i].name)
            }
        }
        success(id, array, names)
    })
}

function publicUsergroupListGet(id, success) {
    publicGetJson('/usergroup', function (rsp) {
        var array = [];
        if ( rsp.code != undefined && rsp.code == 0 ) {
            for (var i = 0; i < rsp.count; i++ ) {
                array.push(rsp.data[i].usergroup)
            }
        }
        success(id, array, array)
    })
}

function publicUserListGet(id, success) {
    publicGetJson('/user', function (rsp) {
        var array = [];
        if ( rsp.code != undefined && rsp.code == 0 ) {
            for ( i = 0; i < rsp.count; i++ ) {
                array.push(rsp.data[i].username)
            }
        }
        success(id, array, array)
    })
}

function publicAccessListGet(id, success) {
    publicGetJson('/access', function (rsp) {
        var array = [];
        if ( rsp.code != undefined && rsp.code == 0 ) {
            for ( i = 0; i < rsp.count; i++ ) {
                array.push(rsp.data[i].name)
            }
        }
        success(id, array, array)
    })
}

function publicRuleListGet(id, success) {
    publicGetJson('/rule', function (rsp) {
        var array = [];
        if ( rsp.code != undefined && rsp.code == 0 ) {
            for ( i = 0; i < rsp.count; i++ ) {
                array.push(rsp.data[i].name)
            }
        }
        success(id, array, array)
    })
}

function publicInterfaceListGet(id, success) {
    publicGetJson('/interface', function (rsp) {
        var array = [];
        var opts = [];
        if ( rsp.code != undefined && rsp.code == 0 ) {
            for ( i = 0; i < rsp.count; i++ ) {
                array.push(rsp.data[i].interface)
                opts.push(rsp.data[i].name + " [" + rsp.data[i].interface + "]")
            }
        }
        success(id, array, opts)
    })
}

function displayWaitClose(title) {
    layer.msg(title);
    setTimeout(function()
    {
        layer.closeAll();
    }, 900);
}

function display(title) {
    layer.msg(title);
}