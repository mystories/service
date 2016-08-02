<?php
$client = new Yar_Client('tcp://localhost:8900');
$client->SetOpt(YAR_OPT_PACKAGER, "msgpack");
$client->setOpt(YAR_OPT_TIMEOUT, 10000);
$client->setOpt(YAR_OPT_CONNECT_TIMEOUT, 10000);
$client->SetOpt(YAR_OPT_PERSISTENT , true);

try {
    print_r($client->__call("Article.Get", ['id'=>1]));
} catch (Exception $ex) {
    var_dump($ex);
}
