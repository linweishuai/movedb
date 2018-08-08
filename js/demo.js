/**
 *模型分析
 */
$json={};
$json['ExportDb']={};
$json['ImportDb']={};
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
		if(is_export){
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
    var element_str = '<li class="'+dbobject.is_export+'" id="' + dbobject.alias + '" model_type="' + dbobject.alias + '">' + dbobject.name + '</li>';
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
	var id = modelId + "_model_" + modelCounter++;
	var type = $(ui.draggable).attr("model_type");
	$(selector).append('<div class="model" id="' + id 
			+ '" modelType="'+ type +'">' 
			+ getModelElementStr(type) + '</div>');
	var left = parseInt(ui.offset.left - $(selector).offset().left);
	var top = parseInt(ui.offset.top - $(selector).offset().top);
	$("#"+id).css("position","absolute").css("left",left).css("top",top);
	//计算添加连接点如果是导出的话 就把锚点设置为右边 导入的话就设置为左边
	var anchor_lenth=dbobject[type]['field'].length+1;
	var jiange=parseFloat(parseFloat(1/anchor_lenth).toFixed(2));
	var jiange_half=parseFloat(parseFloat(jiange/2).toFixed(2))+0.01;
	if(dbobject[type].is_export==1){
        var x_position=1;
	}else{
        var x_position=0;
	}
	//添加一个表级锚点
    instance.addEndpoint(id, { anchors: "TopCenter" }, hollowCircle);
	for(var i=1;i<anchor_lenth;i++){
        instance.addEndpoint(id, { anchors: [[x_position,i*jiange+jiange_half]] }, hollowCircle);
	}
	//注册实体可draggable
	$("#" + id).draggable({
		containment: "parent",
		drag: function (event, ui) {
			instance.repaintEverything();
		},
		stop: function () {
			instance.repaintEverything();
		}
	});
}
//端点样式设置
var hollowCircle = {
	endpoint: ["Dot",{ cssClass: "endpointcssClass"}], //端点形状
	connectorStyle: connectorPaintStyle,
	paintStyle: {
		fill: "#62A8D1",
		radius: 6
	},		//端点的颜色样式
	isSource: true, //是否可拖动（作为连接线起点）
	connector: ["Flowchart", {stub: 30, gap: 0, coenerRadius: 0, alwaysRespectStubs: true, midpoint: 0.5 }],
	isTarget: true, //是否可以放置（连接终点）
	maxConnections: -1
};
//基本连接线样式
var connectorPaintStyle = {
	stroke: "#62A8D1",
	strokeWidth: 2
};
/**
 * 创建模型内部元素
 * @param type
 * @returns {String}
 */
function getModelElementStr(type)
{
	var list='';
    list += '<h4><span index="'
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
        str += '<li><input type="checkbox" name="'
            + type+'.'+[v] + '" value="'
            + type+'.'+obj[v] + '">'
            + type+'.'+obj[v] + '</li>';	}
	return str;
}
//设置连接Label的label
function init(conn)
{
	var label_text;
	$("#select_sourceList").empty();
	$("#select_targetList").empty();
	var sourceName = $("#" + conn.sourceId).attr("modelType");
	var targetName = $("#" + conn.targetId).attr("modelType");
	for(var i = 0; i < metadata.length; i++){
		for(var obj in metadata[i]){
			if(obj == sourceName){
				var optionStr = getOptions(metadata[i][obj].properties,metadata[i][obj].name);
				$("#select_sourceList").append(optionStr);
			}else if(obj == targetName){
				var optionStr = getOptions(metadata[i][obj].properties,metadata[i][obj].name);
				$("#select_targetList").append(optionStr);
			}
		}
	}
	$("#submit_label").unbind("click");
	$("#submit_label").on("click",function(){
		setlabel(conn);
	});
	$("#myModal").modal();
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
	conn.getOverlay("label").setLabel($("#select_sourceList").val() 
			+ ' ' 
			+ $("#select_comparison").val()
			+ ' '
			+ $("#select_targetList").val());
	if($("#twoWay").val()=="true"){
		conn.setParameter("twoWay",true);
	}else{
		conn.setParameter("twoWay",false);
		conn.hideOverlay("arrow_backwards");
	}
}
//删除节点
function removeElement(obj)
{
	var element = $(obj).parents(".model");
	if(confirm("确定删除该模型？"))
		instance.remove(element);
}

