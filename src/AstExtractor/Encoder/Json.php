<?php

namespace AstExtractor\Encoder;

class Json implements Interfaces\EncoderDecoder
{
    private const ENCODING_OPTS = JSON_FORCE_OBJECT | JSON_UNESCAPED_UNICODE;
    private const MAX_DEPTH = 1024;

    private $reader;

    /**
     * Json constructor.
     * The passed $reader will be set as "blocker", so reads will block the process till a new value is read
     * @param $reader
     * @throws \Exception
     */
    public function __construct($reader)
    {
        if ($reader === false || !stream_set_blocking($reader, true)) {
            throw new \Exception('No proper reader passed');
        }

        $this->reader = $reader;
    }

    /**
     * encode returns a string json representation of the passed $input
     * @param array $input
     * @return string json representation
     */
    public static function encode(array $input)
    {
        $encode = json_encode($input, self::ENCODING_OPTS, self::MAX_DEPTH);
        if (!$encode && json_last_error() === JSON_ERROR_UTF8) {
            try {
                self::utf8_encode_recursive($input);
            } catch (\Exception $e) {
                return sprintf(
                    '{"error": "Error#%s, %s. Recursive utf8 encoding failed."}',
                    json_last_error(), json_last_error_msg()
                );
            }

            $encode = json_encode($input, self::ENCODING_OPTS);
        }

        if (!$encode) {
            return sprintf('{"error": "Error#%s, %s"}', json_last_error(), json_last_error_msg());
        }

        return $encode;
    }

    /**
     * decode returns the decode value given a string json representation
     * @param string $input
     * @return mixed|string
     */
    public static function decode(string $input)
    {
        $decoded = json_decode($input, true);
        if (!$decoded) {
            return sprintf('{"error": "Error#%s, %s"}', json_last_error(), json_last_error_msg());
        }

        return $decoded;
    }

    /**
     * utf8_encode_recursive encodes the string contents of the passed $input as valid utf8
     * All inner string contents are scanned, and all array and public object values are converted to utf8
     *
     * @param $input Input data to convert
     * @throws \Exception If something goes wrong
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
                throw new \Exception('Unexpected error during recursive array utf8 encoding');
            }
        }
    }

    /**
     * next reads the next json from the reader, and returns its string representation
     * @throws \Exception
     * @return array|bool|null
     */
    public function next()
    {
        while (!feof($this->reader) && $read = fgets($this->reader)) {
            if ($read === false) {
                throw new \Exception('Error reading from passed stream');
            }

            if (trim($read)==="") continue;

            return [self::decode($read)];
        }
    }
}
