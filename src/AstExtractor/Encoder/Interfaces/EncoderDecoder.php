<?php

namespace AstExtractor\Encoder\Interfaces;

interface EncoderDecoder {
    public static function encode(array $input);
    public static function decode(string $input);
    public function next();
}