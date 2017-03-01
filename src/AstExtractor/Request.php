<?php

namespace AstExtractor;

use AstExtractor\Exception\Fatal;

class Request
{
    public $id;
    public $name;
    public $action;
    public $language;
    public $language_version;
    public $content;

    public const ACTION_PARSE_AST = "ParseAST";

    public const LANG_PHP = "PHP";

    public const PHP_5 = 5;
    public const PHP_6 = 6;
    public const PHP_7 = 7;

    public function __construct(
        $id,
        string $name,
        string $action,
        string $lang,
        $version,
        string $content
    ) {
        $this->id = $id;
        $this->name = $name;
        $this->action = $action;
        $this->language = $lang;
        $this->language_version = $version;
        $this->content = $content;
    }

    public static function fromArray($request) {
        if (!is_array($request) ||
            !isset($request['id']) ||
            !isset($request['name']) ||
            !isset($request['action']) ||
            !isset($request['content'])
        ) {
            throw new Fatal('Wrong request format');
        }

        return new self(
            $request['id'],
            $request['name'],
            $request['action'],
            $request['language'] ? $request['language'] : null,
            $request['language_version'] ? $request['language_version'] : null,
            $request['content']
        );
    }

    public function toArray()
    {
        return [
            "id" => $this->id,
            "name" => $this->name,
            "action" => $this->action,
            "language" => $this->language,
            "language_version" => $this->language_version,
            "content" => $this->content
        ];
    }

    public function answer(array $ast)
    {
        $response = Response::fromRequest($this, $ast);
        $response->ast = $ast;

        return $response;
    }
}
