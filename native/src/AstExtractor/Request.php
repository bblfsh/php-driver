<?php

namespace AstExtractor;

use AstExtractor\Exception\Fatal;

class Request
{
    public $name;
    public $content;

    public const LANG_PHP = "PHP";

    public const PHP_5 = 5;
    public const PHP_6 = 6;
    public const PHP_7 = 7;

    public function __construct(
        string $content,
        $name = null
    ) {
        $this->content = $content;
        if ($name !== null) {
            $this->name = $name;
        }
    }

    public static function fromArray($request) {
        if (!is_array($request) || !isset($request['content'])) {
            throw new Fatal('Wrong request format');
        }

        return new self(
            $request['content'],
            $request['metadata']['name'] ??  null
        );
    }

    public function toArray()
    {
        return [
            "content" => $this->content,
            "metadata" => ["name" => $this->name,],
        ];
    }

    public function answer(array $ast)
    {
        $response = Response::fromRequest($this, $ast);
        $response->ast = $ast;

        return $response;
    }
}
