<?php
trait A {}
trait B {}
class C {
    use A;
}
class D {
    use A, B;
}
