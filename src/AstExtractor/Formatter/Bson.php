<?php

namespace AstExtractor\Formatter;

use AstExtractor\Exception\Fatal;

class Bson extends BaseFormatter
{
    /**
     * @inheritdoc
     */
    public function encode(array $input)
    {
        return bson_encode($input);
    }

    /**
     * @inheritdoc
     */
    public function decode(string $input)
    {
        return bson_decode($input);
    }
}
