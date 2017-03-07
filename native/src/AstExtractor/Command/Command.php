<?php

namespace AstExtractor\Command;

use AstExtractor\Exception\BaseFailure;
use AstExtractor\Formatter\BaseFormatter;
use AstExtractor\Formatter\Json;
use AstExtractor\Formatter\Msgpack;
use AstExtractor\AstExtractor;

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
        $command->init(new Json($stdin), $stdin, $stdout);
    }

    private function init(BaseFormatter $formatter, $stdin, $stdout)
    {
        while (!feof($stdin)) {
            $requests = [];
            try {
                $requests = $formatter->readNext();
            } catch (BaseFailure $e) {
                self::writeErr(null, $e, $stdout, $formatter);
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
                    //TODO: catch different exceptions
                    //  Request::fromArray -> wrong request
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
            $response = $request->answer([]);
            $response->errors = [$e];
            $response->status = Response::getStatus($e->getCode());
        }

        self::write($response, $stdout, $encoder);
    }
}
