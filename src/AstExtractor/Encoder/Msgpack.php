<?php

namespace AstExtractor\Encoder;

use MessagePack\BufferUnpacker;

class Msgpack implements Interfaces\EncoderDecoder
{
    private $reader;
    private $unpacker;

    private const BUFFER_SIZE = 1024;

    /**
     * Msgpack constructor.
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
        $this->unpacker = new BufferUnpacker();
    }

    /**
     * encode returns a string msgpack representation of the passed $input
     * @param array $input
     * @return string msgpack representation
     */
    public static function encode(array $input)
    {
        return msgpack_pack($input);
    }

    /**
     * decode returns the decode value given a string msgpack representation
     * @param string $input
     * @return mixed|string
     */
    public static function decode(string $input)
    {
        return msgpack_unpack($input);
    }

    /**
     * next reads the next msgpack from the reader, and returns its string representation
     * @throws \Exception
     * @return array|bool|null
     */
    public function next()
    {
        $request = null;
        while ($request == null && !feof($this->reader) && $read = fread($this->reader, self::BUFFER_SIZE)) {
            if ($read === false) {
                throw new \Exception('Error reading from passed stream');
            }

            $this->unpacker->append($read);
            $request = $this->unpacker->tryUnpack();
        }

        return $request;
    }
}
