# Duet 언어 명세

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

*   `int`: 정수 (1, 42, 256...)
*   `float`: 실수 (3.14, 6.28...)
*   `str`: 문자열 ("Hello, World!")
*   `bool`: 불리언 (`true`, `false`)
*   `list`: 순서가 있는 값의 목록 ([1, 2, 3])
*   `map`: 키-값 쌍의 맵 
*   `nil`: 값이 없음
*   `fail`: 실패 (fail "에러 메시지")

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

## 6. 표준 함수

| 함수 | 설명 | 예시 |
| --- | --- | --- |
| `print(args...)` | 인자로 받은 값들을 표준 출력에 출력합니다. | `print("Hello", "Duet!")` |
| `type(arg)` | 인자로 받은 값의 데이터 타입(자료형)을 문자열로 반환합니다. | `type(123)`는 `"INTEGER"`를 반환합니다. |
| `len(arg)` | 문자열의 길이나 리스트의 요소 개수를 반환합니다. | `len([1, 2, 3])`는 `3`을 반환합니다. |
| `first(l:list)` | 리스트의 첫 번째 요소를 반환합니다. 리스트가 비어있으면 `nil`을 반환합니다. | `first([10, 20, 30])`는 `10`을 반환합니다. |
| `last(l:list)` | 리스트의 마지막 요소를 반환합니다. 리스트가 비어있으면 `nil`을 반환합니다. | `last([10, 20, 30])`는 `30`을 반환합니다. |
| `rest(l:list):list` | 리스트의 첫 번째 요소를 제외한 나머지 요소들을 새로운 리스트로 반환합니다. | `rest([10, 20, 30])`는 `[20, 30]`을 반환합니다. |
| `push(l:list, el)` | 리스트의 끝에 새로운 요소를 추가한 새 리스트를 반환합니다. | `push([10, 20], 30)`는 `[10, 20, 30]`을 반환합니다. |
| `readln()` | 표준 입력에서 한 줄을 읽어 문자열로 반환합니다. | `supp get_user_input:str -> readln()` |
| `int(arg)` | 인자를 정수로 변환합니다. 문자열, 정수, 불리언 타입을 지원합니다. | `int("123")`는 `123`을 반환합니다. |
| `string(arg)` | 인자를 문자열로 변환합니다. | `string(123)`는 `"123"`을 반환합니다. |
| `bool(arg)` | 인자를 불리언으로 변환합니다. 문자열, 정수, 불리언 타입을 지원합니다. | `bool("true")`는 `true`를 반환합니다. |
| `read(path:str):str` | 파일의 전체 내용을 문자열로 읽어 반환합니다. 파일이 없거나 오류 발생 시 `fail`을 반환합니다. | `read("my_file.txt")` |
| `write(path:str, content:str)` | 문자열을 파일에 씁니다. 성공 시 `true`를, 실패 시 `fail`을 반환합니다. | `write("log.txt", "This is a log.")` |
| `lines(path:str):list` | 파일을 줄 단위로 읽어 `list`로 반환합니다. 실패 시 `fail`을 반환합니다. | `lines("data.csv")` |
| `split(s:str, sep:str):list` | 문자열 `s`를 `sep`을 기준으로 나누어 `list`로 반환합니다. | `split("a,b,c", ",")`는 `["a", "b", "c"]`를 반환합니다. |
| `join(l:list, sep:str):str` | `list`의 요소들을 `sep`을 이용해 합쳐 하나의 문자열로 반환합니다. | `join(["a", "b"], "-")`는 `"a-b"`를 반환합니다. |
| `trim(s:str):str` | 문자열의 앞뒤 공백을 제거합니다. | `trim("  hello  ")`는 `"hello"`를 반환합니다. |
| `replace(s:str, old:str, new:str):str` | 문자열 `s`에서 모든 `old`를 `new`로 교체합니다. | `replace("a-b-c", "-", "/")`는 `"a/b/c"`를 반환합니다. |
| `contains(s:str, sub:str):bool` | 문자열 `s`가 `sub`을 포함하는지 확인합니다. | `contains("hello world", "world")`는 `true`를 반환합니다. |