# toy-project

toy projects are projects for considered stable and will used successfully in production environments.

- k8s-resource-manager
- gitops-golang


golang 언어 작성 시 고려할 사항
1. 함수 내에서 사용하는 변수는 소문자로 시작
2. 각 기능마다 주석 달기
3. 보통 주석이나 로그는 <verb> <target> 으로 통일
4. 파일명은 명사, 함수명은 명시적
5. 자주 사용하는 것들을 struct 로 구성
6. 기능마다 파일을 쪼갤 필요는 없음
7. import 시 다음과 같이 하는 것이 좋음

import (
  
  내장 Package
  
  외부 Repository Package
  
  내부 Repository 

)
