<?php
callAnon(new class {});
callAnon(new class {
    public function method1() {
        echo 'hello';
    }
});
