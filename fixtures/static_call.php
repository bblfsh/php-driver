<?php

A::b();
A::{'b'}();
A::$b();
A::$b['c']();
A::$b['c']['d']();

A::b()['c'];

static::b();
$a::b();
${'a'}::b();
$a['b']::c();
