1. deployment 의 상태가 error 일 경우
- error 인 deployment 취합
- replica 0
- replica 0 && createDt >= 1 month -> delete

2. namespace : rnd-test
- createDt >= 1 month && re-deploy > 1 weeks && no traffic -> delete


go run main.go down -n "rnd-test"