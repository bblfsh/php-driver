<?php

namespace AstExtractor\Command;

use AstExtractor\Formatter\BaseFormatter;
use AstExtractor\Formatter\Json;
use AstExtractor\Formatter\Msgpack;
use AstExtractor\AstExtractor;
use AstExtractor\Request;
use AstExtractor\Response;
use AstExtractor\Exception\Fatal;

class Command
{
    private $extractor;

    public function __construct()
    {
        $this->extractor = new AstExtractor();
    }

    public static function run($argv)
    {
        $command = new Command();

        $stdin = fopen('php://stdin', 'rb');
        $stdout = fopen('php://stdout', 'ab');

        if (isset($argv[1]) && $argv[1] == BaseFormatter::MSGPACK) {
            $formatter = new Msgpack($stdin);
        } else {
            $formatter = new Json($stdin);
        }

        $command->init($formatter, $stdin, $stdout);
    }

    private function init(BaseFormatter $formatter, $stdin, $stdout)
    {
        while (!feof($stdin)) {
            $requests = [];
            try {
                //echo PHP_EOL . PHP_EOL;
                $requests = self::logTime(
                    "readNext", function () use ($formatter) {return $formatter->readNext();}
                );
            } catch (\Exception $e) {
                //TODO: encapsulate $e in this new Fatal
                self::writeErr(null, new Fatal('Wrong request format'), $stdout, $formatter);
                continue;
            }

            if (!is_array($requests)) {
                continue;
            }

            foreach ($requests as $i => $rawReq) {
                $request = null;
                try {
                    //var_dump($rawReq);
                    $request = Request::fromArray($rawReq);
                    $ast = $this->extractor->getAst($request->content);
                    $response = $request->answer($ast);
                    self::write($response, $stdout, $formatter);
                } catch (\Exception $e) {
                    //var_dump("exception");
                    //var_dump($e);
                    //var_dump($e->getTraceAsString());
                    self::writeErr($request, $e, $stdout, $formatter);
                    continue;
                }
            }
        }

        fclose($stdin);

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

    private static function write(Response $response, $stdout, BaseFormatter $encoder)
    {
        $output = self::logTime("encoding", function () use ($encoder, $response) {
            return $encoder->encode($response->toArray()) . PHP_EOL /*. PHP_EOL*/;
        });
        //var_dump("Json:", strlen($msgJson));
        //var_dump($msgJson);
        //echo $output;
        fwrite($stdout, $output);
    }

    private static function writeErr(Request $request = null, \Exception $e, $stdout, BaseFormatter $encoder) {
        if ($request === null) {
            $response = Response::fromError($e);
        } else {
            $response = $request->answer(null);
            $response->status = self::getStatus($e->getCode());
        }

        self::write($response, $stdout, $encoder);
    }

    public static function getStatus($statusCode)
    {
        switch ($statusCode) {
            case BaseFailure::ERROR:
                return Response::STATUS_ERROR;
            case BaseFailure::FATAL:
                return Response::STATUS_FATAL;
        }

        return Response::STATUS_OK;
    }
}
