<?php

namespace AstExtractor;

use AstExtractor\Exception\BaseFailure;
use PhpParser\Node;

class Response
{
    public $id;
    public $name;
    public $driver = self::DRIVER;
    public $language;
    public $language_version;
    public $ast;
    public $status = self::STATUS_PENDING;
    public $errors = [];

    public const DRIVER = "bblfish PHP ast v1";

    public const STATUS_OK = "ok";
    public const STATUS_ERROR = "error";
    public const STATUS_FATAL = "fatal";
    public const STATUS_PENDING = "pending";

    public function __construct(
        $id,
        string $name,
        string $lang,
        $version,
        array $ast
    ) {
        $this->id = $id;
        $this->name = $name;
        $this->language = $lang;
        $this->language_version = $version;
        $this->ast = $ast;
    }

    public static function fromRequest(Request $request, array $ast, ...$err)
    {
        $response = new self(
            $request->id,
            $request->name,
            $request->language,
            $request->language_version,
            $ast
        );

        if (count($err) > 0) {
            $response->errors = $err;
        }

        return $response;
    }

    public static function fromError(\Exception $error)
    {
        $response = new self(null, '', '', null, []);
        $response->errors = [$error->getMessage()];

        return $response;
    }

    public function toArray()
    {
        return [
            'id' => $this->id,
            'name' => $this->name,
            'driver' => $this->driver,
            'language' => $this->language,
            'language_version' => $this->language_version,
            'ast' => $this->ast,
            'status' => $this->status,
            'errors' => $this->errors
        ];
    }
}
