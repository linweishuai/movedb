<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2018/8/8
 * Time: 8:11
 */
    ini_set('display_errors','off');
    $res=['status'=>false,'msg'=>''];
    $host=$_POST['host'];
    $username=$_POST['username'];
    $password=$_POST['password'];
    $dbname=$_POST['dbname'];
    $tablename=$_POST['tablename'];
    $mysql_connection = mysqli_connect($host,$username,$password,$dbname);
    if (!$mysql_connection)
    {
        $res['msg']="链接错误".mysqli_connect_error();
        echo json_encode($res);
    }
    $sql="SHOW COLUMNS FROM `$tablename`";
    $result=$mysql_connection->query($sql);
    if(!$result){
        $res['msg']="查询错误".mysqli_error($mysql_connection);
        echo json_encode($res);
    }
    $rows=$result->fetch_all(MYSQLI_ASSOC);
    $tableinfo_arr=[];
    foreach ($rows as $row){
        $tableinfo_arr[]=$row['Field'];
    }
    $res['status']=true;
    $res['msg']=$tableinfo_arr;
    echo json_encode($res);
