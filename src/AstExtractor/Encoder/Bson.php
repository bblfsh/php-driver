<?php

namespace AstExtractor\Encoder;

class Bson implements Interfaces\EncoderDecoder
{
    public static function encode(array $input)
    {
        return bson_encode($input);
    }

    public static function decode(string $input)
    {
        return bson_decode($input);
    }

    /**
     * next reads the next bson from the reader, and returns its string representation
     * @throws \Exception
     * @return array|bool|null
     */
    public function next()
    {
        // TODO: Implement next() method.
    }
}
