<?php

namespace FixturesGenerator;

use AstExtractor\Encoder\Json;

class Generator
{
    private const SLEEP_SEC = 0;

    private const EXAMPLE_1 = '<?php echo "hello";';
    private const EXAMPLE_2 = '<?php echo "bye";';
    private const EXAMPLE_3 = '
        <body>
            <?php
                echo $a = 3;
                echo "hello world" . ($a + 1);
            ?>
        </body>';
    private const EXAMPLE_4 = '
        <?php
            $iso8859notUtf8 = "\x80\x83\x8a\x8e\x9a\x9e"
            ."\x9f\xa2\xa5\xb5\xc0\xc1\xc2"
            ."\xc3\xc4\xc5\xc7\xc8\xc9\xca"
            ."\xcb\xcc\xcd\xce\xcf\xd1\xd2"
            ."\xd3\xd4\xd5\xd6\xd8\xd9\xda"
            ."\xdb\xdc\xdd\xe0\xe1\xe2\xe3"
            ."\xe4\xe5\xe7\xe8\xe9\xea\xeb"
            ."\xec\xed\xee\xef\xf1\xf2\xf3"
            ."\xf4\xf5\xf6\xf8\xf9\xfa\xfb"
            ."\xfc\xfd\xff";
        
            $double_chars_in = ["\x8c", "\x9c", "\xc6", "\xd0", "\xde", "\xdf", "\xe6", "\xf0", "\xfe"];
        ?>';
    private const EXAMPLE_5 = '
        <?php
            class strangeChars {
                public $a = "\xfe";
                public function __construct($a){$this->a = $a;}
            }
            $b = "\xfe";
            $c = ["\xfe"];
            $d = new strangeChars("\xfe");
        ?>';



    public static function generate($count = 1)
    {
        if (!(boolean)$count--) return;
        echo
            self::encode(Generator::getRequest(10000001, 'FILE_HELLO', Generator::EXAMPLE_1)) .
            self::encode(Generator::getRequest(10000002, 'FILE_BYE', Generator::EXAMPLE_2));

        if (!(boolean)$count--) return;
        sleep(Generator::SLEEP_SEC);
        echo self::encode(Generator::getRequest(10000003, 'FILE_HELLO_WORLD', Generator::EXAMPLE_3));

        if (!(boolean)$count--) return;
        sleep(Generator::SLEEP_SEC);
        echo self::encode(Generator::getRequest(10000003, 'FILE_STRANGE_CHARS', Generator::EXAMPLE_4));

        if (!(boolean)$count--) return;
        sleep(Generator::SLEEP_SEC);
        echo self::encode(Generator::getRequest(10000003, 'FILE_STRANGE_CHARS', Generator::EXAMPLE_5));

        //$globPattern = './tests/fixtures/WordPress__wp-includes__formatting.php';
        //$globPattern = './tests/fixtures/drupal__core__modules__migrate_drupal__tests__fixtures__drupal7.php';
        $globPattern = './tests/fixtures/*.php';
        foreach (glob($globPattern) as $i => $filePath) {
            if (!(boolean)$count--) return;
            sleep(Generator::SLEEP_SEC);
            echo self::encode(
                Generator::getRequest(
                    $i + 10000004,
                    'FILE_' . $filePath,
                    file_get_contents($filePath)
                )
            );
        }
    }

    private static function encode($input)
    {
        return Json::encode($input) . PHP_EOL;
        //return msgpack_pack($input) . PHP_EOL;
    }

    private static function getRequest($id, string $name, string $content)
    {
        return (new Request(
            $id,
            $name,
            Request::ACTION_PARSE_AST,
            Request::LANG_PHP,
            Request::PHP_7,
            $content
        ))->toArray();
    }
}


