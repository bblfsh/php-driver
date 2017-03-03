<?php

namespace AstExtractor;

use AstExtractor\Exception\BaseFailure;

class Response
{
    public $name;
    public $ast;
    public $status = self::STATUS_PENDING;
    public $errors = [];

    public const STATUS_OK = "ok";
    public const STATUS_ERROR = "error";
    public const STATUS_FATAL = "fatal";
    public const STATUS_PENDING = "pending";

    public function __construct(
        array $ast,
        $name = null
    ) {
        $this->ast = $ast;
        if ($name !== null) {
            $this->name = $name;
        }
    }

    public static function fromRequest(Request $request, array $ast, ...$err)
    {

        $response = new self(
            $ast,
            $request->name
        );

        if (count($err) > 0) {
            $response->status = SELF::STATUS_ERROR;
            $response->errors = $err;
        } else {
            $response->status = SELF::STATUS_OK;
        }

        return $response;
    }

    public static function fromError(\Exception $error)
    {
        $response = new self([], null);
        $response->errors = [$error->getMessage()];
        $response->status = SELF::STATUS_ERROR;

        return $response;
    }

    public function toArray()
    {
        return [
            'ast' => $this->ast,
            'metadata' => ['name' => $this->name,],
            'status' => $this->status,
            'errors' => $this->errors
        ];
    }
}
