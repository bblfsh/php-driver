<?php

namespace AstExtractor\Formatter;

use AstExtractor\Exception\Fatal;
use MessagePack\BufferUnpacker;
use MessagePack\Unpacker;
use MessagePack\Packer;

class Msgpack extends BaseFormatter
{
    private const BUFFER_SIZE = 1024; //TODO: test with smaller (example: 20)

    private const EMPTY_LINE_MSGPACK = 10;

    private $bufferUnpacker;
    private $unpacker;
    private $packer;

    /**
     * @inheritdoc
     */
    public function __construct($reader)
    {
        parent::__construct($reader);
        $this->packer = new Packer();
        $this->unpacker = new Unpacker();
        $this->bufferUnpacker = new BufferUnpacker();
    }

    /**
     * @inheritdoc
     */
    public function encode(array $input)
    {
        return msgpack_pack($input);
        //return $this->packer->pack($input); //TODO: it does not know how to encode an AST PhpParser\Node
    }

    /**
     * @inheritdoc
     */
    public function decode(string $input)
    {
        return $this->unpacker->unpack($input);
    }

    /**
     * @inheritdoc
     */
    public function readNext()
    {
        $msg = null;
        while ($msg == null && !feof($this->reader) && $read = fread($this->reader, self::BUFFER_SIZE)) {
            if ($read === false) {
                throw new Fatal('Error reading from the passed stream');
            }

            $this->bufferUnpacker->append($read);
            $msg = $this->bufferUnpacker->tryUnpack();
        }

        return array_filter((Array)$msg, function($v){
            return $v !== self::EMPTY_LINE_MSGPACK;
        });
    }
}
