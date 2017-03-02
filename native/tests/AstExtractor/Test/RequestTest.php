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
            'id' => 33,
            'name' => 'name',
            'action' => 'action',
            'content' => 'content'
        ]);

        $this->assertEquals(33, $req->id);
        $this->assertEquals('name', $req->name);
        $this->assertEquals('action', $req->action);
        $this->assertEquals('content', $req->content);
        $this->assertEquals('', $req->language);
        $this->assertNull($req->language_version);
    }


    /**
     * @expectedException \AstExtractor\Exception\Fatal
     */
    public function testNewRequestFails(): void
    {
        Request::fromArray([]);
    }
}
