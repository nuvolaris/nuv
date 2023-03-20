setup() {
    load 'test_helper/bats-support/load'
    load 'test_helper/bats-assert/load'
    export NO_COLOR=1
}

@test "welcome" {
    run nuv
    assert_line '* sub:           sub command'
    assert_line '* testcmd:       test nuv commands'
}

@test "testcmd" {
    run nuv testcmd
    assert_line "24"
}

@test "sub" {
    run nuv sub
    assert_line '* opts:         opts test'
    assert_line '* simple:       simple'
}

@test "sub simple" {
    run nuv sub simple
    assert_line simple
}
