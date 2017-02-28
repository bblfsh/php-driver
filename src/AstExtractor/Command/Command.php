<?php

namespace AstExtractor\Command;

use AstExtractor\Encoder\Json;
use AstExtractor\Encoder\Msgpack;
use AstExtractor\Extractor\Extractor;

class Command
{
    private $extractor;

    public function __construct()
    {
        $this->extractor = new Extractor();
    }

    public static function run()
    {
        $command = new Command();
        $command->init();
    }

    private function init()
    {
        $stdin = fopen('php://stdin', 'rb');

        //$unpacker = new Msgpack($stdin);
        $unpacker = new Json($stdin);
        while (!feof($stdin)) {
            //echo PHP_EOL . PHP_EOL;
            $requests = self::logTime("readNext", function () use ($unpacker) {return $unpacker->next();});
            if ($requests === null) {
                continue;
            }

            $count = count($requests);
            foreach ($requests as $i => $request) {
                if ($request == 10) continue;

                if (!isset($request['content']) || !isset($request['name'])) {
                    echo sprintf("wrong request%s", PHP_EOL);
                    //var_dump($request);
                    continue;
                }

                echo sprintf("REQUEST %d/%d: '%s'%s", $i, $count, $request['name'], PHP_EOL);
                if (!$this->process($request['content'])) {
                    echo sprintf("Error processing!!%s", PHP_EOL);
                    return false;
                }
            }
        }

        fclose($stdin);

        //echo sprintf("%sFINISH!!%s", PHP_EOL, PHP_EOL);
        return true;
    }

    private function process(string $code)
    {
        $ast = $this->extractor->getAst($code);
        $msgPack = self::logTime("Msgpack encode", function () use ($ast) {return Msgpack::encode($ast);});
        $msgJson = self::logTime("Json encode", function () use ($ast) {return Json::encode($ast);});

        //echo PHP_EOL;
        //var_dump("AST:", $ast);
        //echo PHP_EOL;
        //var_dump("Msgpack:", strlen($msgPack));
        //var_dump($msgPack);
        echo $msgPack;
        echo PHP_EOL;
        //var_dump("Json:", strlen($msgJson));
        //var_dump($msgJson);
        //echo $msgJson.PHP_EOL;

        return true;
    }

    private const NANOSECONDS_MILISECOND = 1000000;
    private static function logTime(string $txt, callable $func)
    {
        exec('date +%s%N', $time0); //microtime(true);
        $res = $func();
        exec('date +%s%N', $time1); //microtime(true);
        //echo sprintf("Time: '%s', %s ms%s", $txt, round(($time1[0] - $time0[0]) / self::NANOSECONDS_MILISECOND), PHP_EOL);
        return $res;
    }

    private static function round($v)
    {
        return round($v, 2);
    }
}
