<?php

/*
 * Part of bblfsh project
 */
namespace AstExtractor;

use PhpParser\ParserFactory;
use PhpParser\Lexer;
use PhpParser\Node;

/**
 * Extractor returns the AST
 *
 * @category Extractor
 * @package  AstExtractor
 * @author   bblfish <hello@bblfish.com>
 * @license  GPL https://github.com/bblfsh
 * @link     https://github.com/bblfsh
 */
class AstExtractor
{
    public const LEXER_CONF = [/*'usedAttributes' => []*/];
    public const PARSER_VERSION = ParserFactory::PREFER_PHP7;

    public $lexer;
    public $parser;

    /**
     * Returns a new Extractor
     */
    public function __construct()
    {
        $this->lexer = new Lexer(self::LEXER_CONF);
        $this->parser = (new ParserFactory)->create(self::PARSER_VERSION, $this->lexer);
    }

    /**
     * Returns the AST
     *
     * @param string $code Source code to analyze
     *
     * @throws \Exception in case it happens a parsing exception
     * @return Node[]|null Array of statements (or null if the 'throwOnError'
     *      option is disabled and the parser was unable to recover from an error).
     */
    public function getAst($code)
    {
        return $this->parser->parse($code);
    }
}

