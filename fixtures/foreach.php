<?php

// foreach on variable
foreach ($a as $b)  {}
foreach ($a as &$b) {}
foreach ($a as $b => $c) {}
foreach ($a as $b => &$c) {}
foreach ($a as list($a, $b)) {}
foreach ($a as $a => list($b, , $c)) {}
foreach (array() as $b) {}
foreach ($a as $b):
endforeach;
