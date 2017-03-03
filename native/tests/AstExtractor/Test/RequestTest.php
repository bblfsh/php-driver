<?php
declare(strict_types=1);

namespace AstExtractor\Test;

use PHPUnit\Framework\TestCase;
use AstExtractor\Request;

class RequestTest extends TestCase
{
    public function testNewRequest(): void
    {
        $req = Request::fromArray([
            'metadata' => ['name' => 'name',],
            'content' => 'content',
        ]);

        $this->assertEquals('name', $req->name);
        $this->assertEquals('content', $req->content);
    }


    /**
     * @expectedException \AstExtractor\Exception\Fatal
     */
    public function testNewRequestFails(): void
    {
        Request::fromArray([]);
    }
}
