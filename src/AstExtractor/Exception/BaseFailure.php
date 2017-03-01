<?php

namespace AstExtractor\Exception;

abstract class BaseFailure extends \Exception
{
    public const ERROR = 1;
    public const FATAL = 2;

    public function __construct($type, string $msg, ...$values)
    {
        $msg = count($values) == 0 ? $msg : sprintf($msg, ...$values);
        parent::__construct($msg, $type, null);
    }
}
