<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2018/8/9
 * Time: 17:51
 */
if(!file_exists("./log/AppLog.log")){
    echo json_encode([]);
    return;
}
$file = fopen("./log/AppLog.log", "r");
$log=[];
$i=0;
if(!file_exists("lastnumber.txt")){
    $lastnumber=0;
    file_put_contents("lastnumber.txt",0);
}else{
    $lastnumber=file_get_contents('lastnumber.txt',$i);
}
//输出文本中所有的行，直到文件结束为止。
while(!feof($file))
{
    $content= fgets($file);
    if($content){
        $log[$i]=$content;
    }
    $i++;
}
fclose($file);
if($i>$lastnumber){
    //说明有新的数据 就输出
    echo json_encode(array_slice($log,$lastnumber));
}
file_put_contents('lastnumber.txt',$i);

