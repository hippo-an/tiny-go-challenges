# Database transaction isolation level
- SQL 표준은 네 가지 수준의 트랜잭션 격리를 정의함
- 먼저 격리 레벨이란 여러 트랜잭션이 동시에 수행될때
  - 각 트랜잭션이 다른 트랜잭션에 대해 어느정도 격리되어 있는가에 대해 수준을 나누어 놓은것
  - 격리 레벨이 중요한 이유는 뭐다?
    - 데이터의 일관성을 유지하기 (데이터 무결성)
- ansi / iso 표준으로 4가지 격리 수준을 정의
- 각각의 격리 수준은 높은 격리 수준부터 낮은 격리 수준의 순서이다.
  - Serializable(직렬화 가능: 모든 트랜잭션이 고유하다)
    - 초특급 파워 격리 상태
    - 일련의 트랜잭션 집합을 어떠한 순서로 실행하든 동일한 결과를 보장하는 격리 수준
  - Repeatable read(반복 가능 읽기: 반복 읽기의 결과가 같다)
    - 동일 트랜잭션 내에서 동일 데이터를 여러번 읽어도 동일한 결과 보장하는 격리 수준
    - 단, 특정 조건을 만족하는 여러 데이터를 읽을 때, 다른 트랜잭션에서 조건을 만족하는 데이터를 추가(혹은 삭제)하는 경우 
    - 처음 쿼리의 결과 리스트와 두번째 쿼리의 결과 리스트가 다르게 나타나는 현상인
      - **Phantom Read** 발생 가능
  - Read Committed(커밋된 데이터 읽기: 다른 트랜잭션이 커밋한 놈만)
    - 다른 트랜잭션이 커밋한 데이터만 읽을 수 있다
    - 처음 데이터를 읽을때와 다른 트랜잭션이 커밋되고 다시 읽는 데이터가 동일하지 않은 현상인
      - **Non-repeatable Read** 발생 가능
  - Read Uncommitted(커밋되지 않은 데이터 읽기: 다른 트랜잭션이 커밋하지 않은 데이터도 읽는다 oh my god)
    - 다른 트랜잭션이 커밋하지 않은 데이터도 읽을 수 있다.
    - 다른 트랜잭션의 커밋 여부에 관계 없이 데이터를 읽을 수 있기 때문에 롤백이 되는 경우 완전히 잘못된 데이터를 읽는 현상인 
      - **dirty-read** 발생 가능



# Read Phenomenon
- 위에서 발생가능한 하나의 트랜잭션에 의해 다른 트랜잭션이 영향을 받아 데이터의 일관성에 문제가 생길 수 있는 현상을 말한다.
- 이는 데이터베이스의 격리 수준에 따라서 다른 읽기 현상이 발생한다.
- Dirty Read
  - 커밋되지 않은 데이터를 읽는 현상으로
  - 다른 트랜잭션이 데이터를 롤백하는 경우 읽은 데이터가 사라질 수 있어 데이터의 일관성에 문제가 생긴다.
  - 이는 Read Uncommitted 격리 수준에서 발생할 수 있는 읽기 현상으로
  - 아직 commit 되지 않은 변경된 데이터를 읽어오기 때문에 발생할 수 있는 현상이다.
  ```text
  하나의 트랜잭션에서 a 라는 값을 b 로 업데이트하고 아직 commit 하지 않은 상태에서
  다른 트랜잭션에서 해당 값을 읽으면 b 의 값을 읽어간다. (read uncommitted 격리 수준이기 때문에 커밋되지 않은 값을 읽을 수 있다.)
  그런 상황에서 업데이트를 하는 트랜잭션에 문제가 생겨 롤백이 되는 경우
  원래 값인 a 로 돌아가게 되지만 다른 트랜잭션에서 읽어간 값은 b 가 되어 드러운 읽기(dirty-read) 현상이 발생한 것이다.
  ```
  
  
- Non-Repeatable Read
  - 트랜잭션 내에서 동일한 데이터를 두 번 읽은 경우 서로 다른 결과를 얻는 현상으로
  - Read Committed 격리 수준에서 발생할 수 있는 읽기 현상이다.
  ```text
  하나의 트랜잭션에서 a 라는 값을 읽고 다른 트랜잭션에서 같은 값을 b 로 업데이트한 후 commit 을 했다고 해보자.
  업데이트가 완료된 후 값을 읽으면 하나의 트랜잭션에서 a 와 b 라는 서로 다른 값이 읽히는
  반복 불가능한 읽기(non-repeatable read) 현상이 발생한 것이다.
  ```

- Phantom Read
  - 트랜잭션이 시작된 후 다른 트랜잭션에 의해 새로 추가된 데이터를 읽는 현상으로
  - Repeatable Read 격리 수준에서 발생할 수 있다.
  ```text
  업데이트 되는 값에 대해서 반복적으로 읽기는 가능하지만
  a 라는 컬럼의 값이 10 이상인 값에 대해서 모두 조회할 때
  하나의 트랜잭션에서 처음 값을 읽고 난 다음에 다른 트랜잭션에서 특정 값을 추가(혹은 삭제)한 경우
  조건에 해당하는 조회 결과가 다르게 나타나는 유령 읽기(Phantom Read) 현상이 발생한 것이다.
  ```

- Serialization Anomaly
  - 직렬화 이상 현상은 여러 트랜잭션이 동시에 실행될때 발생할 수 있는 문제로 
  - 트랜잭션의 결과가 어떤 순서로도 순차적으로 실행한 결과와 일치하지 않는 경우의 문제이다.
  - 이는 데이터 베이스 일관성을 해치며 예기치 않은 결과를 초래할 수 있다.
  - 이는 여러 트랜잭션이 동시적으로 수행되며 데이터베이스의 각각의 격리 수준에서 발생할 수 있는 읽기 현상이 발생하여
  - 다른 트랜잭션에 영향을 끼치며 발생한다.
  - 이를 해결하기 위해 데이터베이스의 격리 수준을 높여 Serializable 로 설정하여 예방할 수 있지만
  - 시스템 성능에 영향을 줄 수 있기 때문에 적절한 균형이 필요하다.



| Isolation Level    | Dirty Read                    | Nonrepeatable Read | Phantom Read                    | Serialization Anomaly |
|--------------------|-------------------------------|---------------------|---------------------------------|-----------------------|
| Read uncommitted   | Allowed, but not in PG        | Possible            | Possible                        | Possible              |
| Read committed     | Not possible                  | Possible            | Possible                        | Possible              |
| Repeatable read    | Not possible                  | Not possible        | Allowed, but not in PG          | Possible              |
| Serializable       | Not possible                  | Not possible        | Not possible                    | Not possible          |
출처: [Transaction Isolation Levels of PostgreSQL](https://www.postgresql.org/docs/current/transaction-iso.html)
- 위는 PostgreSQL 의 공식문서 transaction isolation 내용에서 발췌한 표이다.
- Postgresql 에선 Read Uncommitted 격리 수준에서 발생하는 Dirty Read 를 허용하지 않으며,
- PostgreSQL 에선 Read Uncommitted 격리 수준은 일반적인 Read Committed 처럼 동작한다.


