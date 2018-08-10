/**
 *模型分析
 */
$json={};
$json['ExportDb']={};
$json['ImportDb']={};
$json['Fieldrule']={};
$json['RuleField']={};
$json['Tablerule']={};
dragtable=[];
dbobject=[];
function getdbinfo() {
    var alias =($("#adddbobject input[name=alias]").val());
    var host =($("#adddbobject input[name=host]").val());
    var username =($("#adddbobject input[name=username]").val());
    var passwd =($("#adddbobject input[name=passwd]").val());
    var dbname =($("#adddbobject input[name=dbname]").val());
    var tablename =($("#adddbobject input[name=tablename]").val());
    var is_export =($("#adddbobject select[name=is_export]").val());
    //后台获取数据库对象
	if(dbobject[alias]!=undefined){
		alert('数据库别名已经使用,请更换');
		return;
	}
    $.post("getDbinfo.php",{host:host,username:username,password:passwd,dbname:dbname,tablename:tablename},function(res) {
        if (!res.status) {
            alert(res.msg);
            return;
        }
        var tempobject={};
        tempobject.alias=alias;
        tempobject.name=alias+"."+tablename;
        tempobject.field=[];
        for(var i in res.msg){
            tempobject.field.push(res.msg[i])
		}
        tempobject.is_export=is_export;
        setLeftMenu(tempobject);
        dbobject[alias]=tempobject;
		//todo 添加到exportdb里面你去
		if(is_export==1){
            $json['ExportDb'][alias]={
                'host': host,
                "port": '3306',
                'username': username,
                'passwd': passwd,
                'dbname': dbname,
                'tablename': tablename,
            };
		}else{
            $json['ImportDb'][alias]={
                'host': host,
                "port": '3306',
                'username': username,
                'passwd': passwd,
                'dbname': dbname,
                'tablename': tablename,
            };
		}
    },'json');
	console.log($json);
}


/**模型计数器*/
var modelCounter = 0;
/**
 * 初始化一个jsPlumb实例
 */
var instance = jsPlumb.getInstance({
	DragOptions: { cursor: "pointer", zIndex: 2000 },
	ConnectionOverlays: [
			[ "Arrow", {
			    location: 1,
			    visible:true,
			    width:11,
			    length:11,
			    direction:1,
			    id:"arrow_forwards"
			} ],
			[ "Arrow", {
			    location: 0,
			    visible:true,
			    width:11,
			    length:11,
			    direction:-1,
			    id:"arrow_backwards"
			} ],
            [ "Label", {
                location: 0.5,
                id: "label",
                cssClass: "aLabel"
            }]
        ],
	Container: "container"
});
instance.importDefaults({
	 ConnectionsDetachable:true,
	 ReattachConnections:true
});
/**
 * 设置左边菜单
 * @param Data
 */
function setLeftMenu(dbobject)
{
	var  classname=dbobject.is_export==1?'dbexport':'dbimport';
    var element_str = '<li class="'+classname+'" id="' + dbobject.alias + '" model_type="' + dbobject.alias + '">' + dbobject.name + '</li>';
    $("#leftMenu").append(element_str);
	$("#leftMenu li").draggable({
		helper: "clone",
		scope: "plant"
	});
	$("#container").droppable({
		scope: "plant",
		drop: function(event, ui){
			CreateModel(ui, $(this));
		}
	});
}
/**
 * 添加模型
 * @param ui
 * @param selector
 */
function CreateModel(ui, selector)
{
	var modelId = $(ui.draggable).attr("id");
	var model_id = modelId + "_wojiubuxinnengyourenhewoyiyang_" + modelCounter++;
	var type = $(ui.draggable).attr("model_type");
    //如果已经拖拽过 那么就不能再次拖拽了
    if(dragtable[type]!=undefined){
        alert('该表已经添加')
        return;
    }
	$(selector).append('<div id="'+model_id+'_zoom" class="modezoom"><div class="panzoom"><div class="model" id="' + model_id
			+ '" modelType="'+ type +'">' 
			+ getModelElementStr(type) + '</div></div></div>');
	var left = parseInt(ui.offset.left - $(selector).offset().left);
	var top = parseInt(ui.offset.top - $(selector).offset().top);
    $('#'+model_id+"_zoom").css("position","absolute").css("left",left).css("top",top);


    //todo锚点放大和缩小问题
    var $section = $('#'+model_id+"_zoom");
    var $panzoom = $section.find('.model').panzoom();
    $panzoom.parent().on('mousewheel.focal', function( e ) {
        e.preventDefault();
        var delta = e.delta || e.originalEvent.wheelDelta;
        var zoomOut = delta ? delta < 0 : e.originalEvent.deltaY > 0;
        $panzoom.panzoom('zoom', zoomOut, {
            animate: false,
            focal: e
        });
        instance.repaintEverything();
    });
	if(dbobject[type].is_export==1){
        var x_position=1;
	}else{
        var x_position=0;
	}
	//添加一个表级锚点
    instance.addEndpoint(model_id, { anchors: "TopCenter" }, hollowCircle);
	for(var i in dbobject[type]['field']){
		var id=type+'.'+dbobject[type]['field'][i];
        instance.addEndpoint(id, { anchors: [[x_position,0.5]] }, hollowCircle);
        if(dbobject[type].is_export==1){//如果是导出表就两侧加上锚点
            instance.addEndpoint(id, { anchors: [[0,0.5]] }, hollowCircle);
        }
	}
    dragtable[type]=1;
	//注册实体可draggable
	$("#" + model_id).draggable({
		containment: "parent",
		drag: function (event, ui) {
			instance.repaintEverything();
		},
		stop: function () {
			instance.repaintEverything();
		}
	});
}
//基本连接线样式
var connectorPaintStyle = {
    stroke: "#62A8D1",
    strokeWidth: 2
};
// 鼠标悬浮在连接线上的样式
var connectorHoverStyle = {
    strokeWidth: 4,
    lineWidth: 5,
    stroke: "#e41263",
};
//端点样式设置
var hollowCircle = {
	endpoint: ["Dot",{ cssClass: "endpointcssClass"}], //端点形状
	connectorStyle: connectorPaintStyle,
    connectorHoverStyle: connectorHoverStyle,
	paintStyle: {
		fill: "#62A8D1",
		radius: 6
	},		//端点的颜色样式
	isSource: true, //是否可拖动（作为连接线起点）
	connector: ["Flowchart", {stub: 30, gap: 0, coenerRadius: 0, alwaysRespectStubs: true, midpoint: 0.5 }],
	isTarget: true, //是否可以放置（连接终点）
	maxConnections: -1,
};

/**
 * 创建模型内部元素
 * @param type
 * @returns {String}
 */
function getModelElementStr(type)
{
    var classname=dbobject[type].is_export==1?'dbexport':'dbimport';
	var list='';
    list += '<h4 class="'+classname+'"><span index="'
        + dbobject[type].alias + '">'
        + dbobject[type].name
        + '</span><a href="javascript:void(0)" class="pull-right" onclick="removeElement(this);">X</a>'
        + '</h4>';
    list += '<ul>'
    var properties = dbobject[type].field;
    list += parseProperties(type,properties);
    list += '</ul>';
	return list;
}
/**
 * 循环遍历properties
 * @param obj
 * @returns {String}
 */
function parseProperties(type,obj)
{
	var str = "";
	for(var v in obj){
        str += '<li id="'+type+'.'+obj[v]+'"><input type="checkbox" name="'
            + type+'.'+obj[v] + '" value="'
            + type+'.'+obj[v] + '">'
            + type+'.'+obj[v] + '</li>';	}
	return str;
}
//设置连接Label的label
function init(conn)
{
	$("#select_sourceList").empty();
	$("#select_targetList").empty();
	// console.log(conn);
	var sourceName = conn.sourceId;
	var targetName = conn.targetId;
    sourceName_arr=sourceName.split('.');
    // console.log(sourceName_arr);
    // console.log($json['ExportDb']);
	//如果连接的是导出表 那么就报错
	$("#submit_label").unbind("click");
    $("#submit_tablelabel").unbind("click");
	$("#submit_label").on("click",function(){
		setlabel(conn);
	});
    $("#submit_tablelabel").on("click",function(){
        setlabel(conn);
    });
    var sourceName = conn.sourceId;
    if(sourceName.indexOf("wojiubuxinnengyourenhewoyiyang")!=-1){
        $("#mytableModal").modal();
	}else{
        $("#myModal").modal();
	}
}
/**
 * 获取option
 * @param obj
 * @returns {String}
 */
function getOptions(obj,head)
{
	var str = "";
	for(var v in obj){
		if(obj[v].properties == undefined){
			var val = head + '.' + obj[v].des;
			str += '<option value="' + val + '">'
					+val
					+'</option>';
		}else{
			str += arguments.callee(obj[v].properties,head);
		}
	}
	return str;
}
//setlabel
function setlabel(conn)
{
    var sourceName = conn.sourceId;
    var targetName = conn.targetId;
	if(sourceName.indexOf("wojiubuxinnengyourenhewoyiyang")!=-1){
        var condition=$("#condition").val();
        var sourcetable=sourceName.split('_')[0];
        var targettable=targetName.split('_')[0];
        if($json['ExportDb'][sourcetable]!=undefined){//如果是从来源表导出 那么说明连接是正确的
            var scource=sourcetable;//来源
            var target=targettable;//目标
        }else{//如果来源是从导入表来的 那么 就说明连接是反着的
            var scource=targettable;//来源
            var target=sourcetable;//目标
        }
        if(target==scource){
        	alert('不可同表连接');
			instance.detach(conn);
		}
        conn.getOverlay("label").setLabel(scource
            + '=>' +target+':'+condition
        );
        var tempdata=['table',scource,target,condition];
        conn.setData(tempdata)
	}else {
        var sourcetable=sourceName.split('.')[0];
        var targettable=sourceName.split('.')[0];
        if($json['ExportDb'][sourcetable]!=undefined){//如果是从来源表导出 那么说明连接是正确的
            var scource=sourceName;//来源
            var target=targetName;//目标
        }else{//如果来源是从导入表来的 那么 就说明连接是反着的
            var scource=targetName;//来源
            var target=sourceName;//目标
        }
        if(sourcetable==targettable){
            alert('不可同表连接');
            instance.detach(conn);
        }
        if($json['Fieldrule'][target]!=undefined){
            alert('导入字段唯一,不可重复连接');
            instance.detach(conn);
            return;
        }
        var mode=$("#modename").val();
        var extra=$("#extra_data").val();
        if(mode=='OnetoOne'&&extra==''){
            alert('OnetoOne模式请设置对应数据');
            instance.detach(conn);
        }
        //存储一下现在的对应值
        if(mode=='OnetoOne'){
            $json['Fieldrule'][target]=[scource,mode,extra];
        }else{
            $json['Fieldrule'][target]=[scource,mode];
        }
        conn.getOverlay("label").setLabel(scource
            + '=>' +target
            + ':'+ mode
            + ':'+extra
        );
        var tempdata=['field',target,scource,mode,extra];
        conn.setData(tempdata)
	}
}
//删除节点
function removeElement(obj)
{
	var element = $(obj).parents(".model");
	//删除leftmenu 里面的元素
	var leftmenuid=element.attr('id').split('_')[0];
	$('#'+leftmenuid).remove();
	//删除掉json中的元素
	if($json['ExportDb'][leftmenuid]!=undefined){
		delete $json['ExportDb'][leftmenuid];
	}
    if($json['ImportDb'][leftmenuid]!=undefined){
        delete $json['ImportDb'][leftmenuid];
    }
    delete dbobject[leftmenuid];
    delete dragtable[leftmenuid];

	if(confirm("确定删除该模型？"))
		instance.remove(element);
}
function saveallConn() {
    var connections = instance.getConnections();
    $json['Fieldrule']={};
    $json['RuleField']={};
    for (var i in connections){
    	var conndata=connections[i].getData();
    	if(conndata[0]=='field'){
            if(conndata[3]=='OnetoOne'){
                $json['Fieldrule'][conndata[1]]=[conndata[2],conndata[3],conndata[4]]
            }else{
                $json['Fieldrule'][conndata[1]]=[conndata[2],conndata[3]]
            }
		}else{
            var reg=/\$(\w+)/g;
            var condtion=conndata[3];
            var ruleField_arr=condtion.match(reg);
            for (var i in ruleField_arr){
                if($json['RuleField'][conndata[1]]){
                    $json['RuleField'][conndata[1]].push(ruleField_arr[i].replace('\$',''))
                }else{
                    $json['RuleField'][conndata[1]]=[];
                    $json['RuleField'][conndata[1]].push(ruleField_arr[i].replace('\$',''))
                }
            }
            $json['Tablerule'][conndata[1]+conndata[2]]=conndata[3];
		}
	}
    //把jsonRuleField 去一下重 然后在json
    for(var tablename in $json['RuleField']){
        var temp = [];
        for(var i = 0; i < $json['RuleField'][tablename].length; i++) {
            //如果当前数组的第i项在当前数组中第一次出现的位置是i，才存入数组；否则代表是重复的
            if($json['RuleField'][tablename].indexOf($json['RuleField'][tablename][i]) == i){
                temp.push($json['RuleField'][tablename][i])
            }
        }
        $json['RuleField'][tablename]=temp;
    }
    jsonstr=JSON.stringify($json);
    //传到后台执行
    console.log(jsonstr);
    funDownload(jsonstr, 'config.json');
}
	function funDownload(content, filename) {
		var eleLink = document.createElement('a');
		eleLink.download = filename;
		eleLink.style.display = 'none';
		// 字符内容转变成blob地址
		var blob = new Blob([content]);
		eleLink.href = URL.createObjectURL(blob);
		// 触发点击
		document.body.appendChild(eleLink);
		eleLink.click();
		// 然后移除
		document.body.removeChild(eleLink);
	}



