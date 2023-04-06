#!/usr/bin/env bash

check_board_response() {
    board="${1}"
    expected_response="${2}"
    response_code=$(curl -s -o /dev/null -w "%{http_code}" -X "GET" "http://localhost:4433/get_boards?board=${board}")

    if [[ "${response_code}" == "${expected_response}" ]]; then
        return 0
    else
        echo "Unexpected response: ${response_code}"
        return 1
    fi
}

run_test() {
    board=${1:?"You must specify a board"}
    expected_response=${2:?"You must specify an expected response"}

    check_board_response "${board}" "${expected_response}" \
        && echo "Got expected ${expected_response} from ${board}" \
        || { echo "Didn't get expected ${expected_response} from ${board}"; ((num_fails++)); }
    ((num_tests_ran++))
    sleep 1
}

num_fails=0
num_tests_ran=0
# Run tests
run_test "g" "200"
run_test "biz" "200"
run_test "hut" "400"
run_test "gdx" "400"
run_test "pol" "200"
run_test "b" "200"
run_test "test" "400"
run_test "jx" "400"

num_passes=$(echo "${num_tests_ran} - ${num_fails}" | bc)

echo "Tests completed: ${num_tests_ran}"
echo "Tests passed: ${num_passes}"
echo "Tests failed: ${num_fails}"