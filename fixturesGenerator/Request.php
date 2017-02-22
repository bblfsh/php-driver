<?php

namespace FixturesGenerator;

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
}