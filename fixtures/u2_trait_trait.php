<?php
trait A {}
trait B {}
trait C { use A; }
trait D { use A, B; }
