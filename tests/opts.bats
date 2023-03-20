setup() {
    load 'test_helper/bats-support/load'
    load 'test_helper/bats-assert/load'
    export NO_COLOR=1
}

@test "help" {
run nuv sub opts
# just one as it is a cat of a message
assert_line "Usage:"
run nuv sub opts -h
assert_line "Usage:"
run nuv sub opts --help
assert_line "Usage:"
# do not check the actual version but ensure the output is not the help test
run nuv sub opts --version
refute_output "Usage:"
}

@test "cmd" {
    run nuv sub opts hello
    assert_line "hello!"
}

@test "args" { 
    run nuv sub opts args mike
    assert_line "name: mike"
    assert_line "-c: no"

    run nuv sub opts args mike miri -c
    assert_line "name: mike"
    assert_line "name: miri"
    assert_line "-c: yes"
}

@test "arg1 arg3" {
    run nuv sub opts arg1 aaa arg2 1 2 --fl=ag
    assert_line "arg1 name=('aaa') arg2 x=1 y=2 --fl=ag"
    run nuv sub opts arg3 opt1 10 20 --fa
    assert_line "arg3=true opt1=true opt2=false x=10 y=20 --fa=true --fb=false"
}

@test "errors" {
    run nuv sub opts arg1
    assert_line "Usage:"
    run nuv sub opts arg opt4
    assert_line "Usage:"
}