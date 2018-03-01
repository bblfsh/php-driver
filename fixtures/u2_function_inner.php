<?php
function func1() {
    function func2(){}
    func2();
}

function func3() {
    function func4() {
        function func5() {}
        func5();
    }
    func4();
}
if (1) {
    function func5() {}
}
