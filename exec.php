<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2018/8/9
 * Time: 17:51
 */
set_time_limit(0);
$configjson=$_POST['jsonstr'];
//写入文件
file_put_contents('config.json',$configjson);
sleep(1);
$a = exec("main.exe", $out, $status);
