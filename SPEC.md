# Duet 언어 명세 (간략화 버전)

## 1. 핵심 개념

Duet의 모든 연산은 **함수**를 통해 이루어집니다. 함수는 데이터의 흐름을 정의하며, 파이프라인을 통해 서로 연결될 수 있습니다.

## 2. 함수 종류

함수는 입력과 출력의 유무에 따라 세 가지 기본 형태로 나뉩니다.

*   **Processor (`proc`):** 입력과 출력이 모두 있는 함수. 데이터를 가공합니다.
    ```duet
    proc add(a:int, b:int):int -> a + b
    ```
*   **Consumer (`cons`):** 입력만 있는 함수. 주로 외부 시스템에 영향을 줍니다. (예: 출력)
    ```duet
    cons print_message(msg:str) -> print(msg)
    ```
*   **Supplier (`supp`):** 출력만 있는 함수. 주로 데이터를 생성합니다.
    ```duet
    supp get_random_num:float -> random()
    ```

### 실패 가능한 함수 (`eproc`, `esupp`)

함수 이름 앞에 `e`를 붙여 함수가 실패할 수 있음을 명시합니다. 실패 시 오류를 반환하여 안정성을 높입니다.

```duet
eproc safe_divide(a:int b:int):int -> if b == 0 then fail("0으로 나눌 수 없습니다") else a / b
```

## 3. 데이터 타입

기본 데이터 타입은 다음과 같습니다.

*   `int`: 정수
*   `float`: 실수
*   `str`: 문자열
*   `bool`: 불리언 (`true`, `false`)
*   `list`: 순서가 있는 값의 목록
*   `map`: 키-값 쌍의 맵
*   `nil`: 값이 없음

## 4. 제어 흐름

### 조건문

`if-then-else`는 값을 반환하는 표현식입니다.

```duet
proc get_grade(score:int):string -> if score >= 90 then "A" else "B"
```

### 반복문

`for-in`은 컬렉션을 순회하며 새로운 `list`를 반환합니다.

```duet
proc double_all(numbers:list):list -> for n in numbers then n * 2
```

## 5. 파이프라이닝 (`|>`)

`|>` 연산자는 여러 함수를 연결하여 데이터의 흐름을 만듭니다. 한 함수의 출력이 다음 함수의 입력으로 전달됩니다.

```duet
supp get_input -> readln()
proc process_data(data:str) -> upper(data)
cons print_output(result:str) -> print(result)

// 실행 파이프라인
get_input |> process_data |> print_output
```