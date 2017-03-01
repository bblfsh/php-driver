<?php

namespace AstExtractor\Formatter;

use AstExtractor\Exception\Fatal;

class Json extends BaseFormatter
{
    /**
     * ENCODING_OPTS forces encoding to "return always an object"
     *   and "encode multibyte Unicode characters literally"
     */
    private const ENCODING_OPTS = JSON_FORCE_OBJECT | JSON_UNESCAPED_UNICODE;

    /**
     * DECODE_USING_ASSOC_ARRAY is true if the decode value will be represented by an
     *  associative array instead of by an object
     */
    private const DECODE_USING_ASSOC_ARRAY = true;

    /**
     * @inheritdoc
     */
    public function encode(array $input)
    {
        $encode = json_encode($input, self::ENCODING_OPTS);
        if (!$encode && json_last_error() === JSON_ERROR_UTF8) {
            self::utf8_encode_recursive($input);
            $encode = json_encode($input, self::ENCODING_OPTS);
        }

        if (!$encode) {
            throw new Fatal('Error#%s, %s', json_last_error(), json_last_error_msg());
        }

        return $encode;
    }

    /**
     * @inheritdoc
     */
    public function decode(string $input)
    {
        $decoded = json_decode($input, self::DECODE_USING_ASSOC_ARRAY);
        if (!$decoded) {
            throw new Fatal('Error#%s, %s', json_last_error(), json_last_error_msg());
        }

        return $decoded;
    }

    /**
     * @inheritdoc
     */
    public function readNext()
    {
        while (!feof($this->reader) && $read = fgets($this->reader)) {
            if ($read === false) {
                throw new Fatal('Error reading from passed stream');
            }

            if (trim($read) === "") {
                continue;
            }

            return [self::decode($read)];
        }
    }

    /**
     * utf8_encode_recursive encodes the string contents of the passed $input as valid utf8
     * All inner string contents are scanned, and all array and public object values are converted to utf8
     * @param $input Input data to convert
     * @throws \Exception If something went wrong
     */
    private static function utf8_encode_recursive(&$input)
    {
        if (is_string($input)) {
            $input = utf8_encode($input);
        }

        if (is_object($input)) {
            $ovs = get_object_vars($input);
            foreach ($ovs as $k => $v)    {
                self::utf8_encode_recursive($input->$k);
            }
        }

        if (is_array($input)) {
            $success = array_walk_recursive($input, function(&$v){self::utf8_encode_recursive($v);});
            if (!$success) {
                throw new Fatal('Unexpected error during recursive array utf8 encoding');
            }
        }
    }
}
