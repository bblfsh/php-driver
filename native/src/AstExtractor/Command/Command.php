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
        $this->extractor = new AstExtractor(AstExtractor::LEXER_CONF_VERBOSE);
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
                $requests = $formatter->readNext();
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
                    $request = Request::fromArray($rawReq);
                    $ast = $this->extractor->getAst($request->content);
                    $response = $request->answer($ast);
                    self::write($response, $stdout, $formatter);
                } catch (\Exception $e) {
                    self::writeErr($request, $e, $stdout, $formatter);
                    continue;
                }
            }
        }

        fclose($stdin);

        return true;
    }

    private static function write(Response $response, $stdout, BaseFormatter $encoder)
    {
        $output = $encoder->encode($response->toArray()) . PHP_EOL /*. PHP_EOL*/;
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
