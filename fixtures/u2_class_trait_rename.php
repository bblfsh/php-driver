<?php
trait A {
    public function func1() {}
}
class B {
    use A { func1 as funcrename; }
}
class C {
    use A { func1 as protected rename2; }
}
