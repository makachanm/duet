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
*   
### 실패 가능 데이터 타입

타입 선언 뒤에 `?`를 추가하여 해당 타입이 정상 값 또는 `FAIL` 객체를 가질 수 있음을 나타낼 수 있습니다. (예: `str?`, `int?`)
이 구문은 함수의 반환 타입이나 매개변수 타입으로 사용할 수 있습니다. 이를 통해 함수가 `FAIL` 객체를 반환하거나 인자로 받을 수 있음을 명시적으로 보여줍니다.

```duet
proc handle_error(input:str?):str -> if is_fail(input) then "오류가 발생했습니다." else "정상 값: " + input
```

## 4. 제어 흐름

### 조건문

`if-then-else`는 값을 반환하는 표현식입니다.

```duet
proc get_grade(score:int):string -> if score >= 90 then "A" else "B"
```

`match`는 여러 케이스를 비교하는 표현식입니다.
```duet
proc get_grade(score:int):string -> match score { 
is score > 90 then "A"
is score > 80 then "B"
is score > 70 then "C"
is score > 60 then "D"
default "F"
}
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

### 6.1. 입출력 (Input/Output)

| 함수 | 설명 | 예시 |
| --- | --- | --- |
| `print(args...)` | 인자로 받은 값들을 표준 출력에 출력합니다. | `print("Hello", "Duet!")` |
| `readln()` | 표준 입력에서 한 줄을 읽어 문자열로 반환합니다. | `supp get_user_input:str -> readln()` |
| `read(path:str):str` | 파일의 전체 내용을 문자열로 읽어 반환합니다. | `read("my_file.txt")` |
| `write(path:str, content:str)` | 문자열을 파일에 씁니다. 성공 시 `true`를 반환합니다. | `write("log.txt", "This is a log.")` |
| `lines(path:str):list` | 파일을 줄 단위로 읽어 `list`로 반환합니다. | `lines("data.csv")` |

### 6.2. 타입 변환 (Type Conversion)

| 함수 | 설명 | 예시 |
| --- | --- | --- |
| `int(arg)` | 인자를 정수로 변환합니다. | `int("123")`는 `123`을 반환합니다. |
| `string(arg)` | 인자를 문자열로 변환합니다. | `string(123)`는 `"123"`을 반환합니다. |
| `bool(arg)` | 인자를 불리언으로 변환합니다. | `bool("true")`는 `true`를 반환합니다. |
| `type(arg)` | 인자의 데이터 타입을 문자열로 반환합니다. | `type(123)`는 `"INTEGER"`를 반환합니다. |

### 6.3. 리스트 조작 (List Manipulation)

| 함수 | 설명 | 예시 |
| --- | --- | --- |
| `len(l:list)` | 리스트의 요소 개수를 반환합니다. | `len([1, 2, 3])`는 `3`을 반환합니다. |
| `first(l:list)` | 리스트의 첫 번째 요소를 반환합니다. | `first([10, 20])`는 `10`을 반환합니다. |
| `last(l:list)` | 리스트의 마지막 요소를 반환합니다. | `last([10, 20])`는 `20`을 반환합니다. |
| `rest(l:list):list` | 첫 요소를 제외한 새 리스트를 반환합니다. | `rest([10, 20])`는 `[20]`을 반환합니다. |
| `push(l:list, el)` | 끝에 요소를 추가한 새 리스트를 반환합니다. | `push([10], 20)`는 `[10, 20]`을 반환합니다. |

### 6.4. 문자열 조작 (String Manipulation)

| 함수 | 설명 | 예시 |
| --- | --- | --- |
| `len(s:str)` | 문자열의 길이를 반환합니다. | `len("hello")`는 `5`를 반환합니다. |
| `split(s:str, sep:str):list` | 문자열을 `sep` 기준으로 나누어 리스트로 반환합니다. | `split("a,b", ",")`는 `["a", "b"]`를 반환합니다. |
| `join(l:list, sep:str):str` | 리스트 요소를 `sep`으로 합쳐 문자열로 반환합니다. | `join(["a", "b"], "-")`는 `"a-b"`를 반환합니다. |
| `trim(s:str):str` | 문자열의 앞뒤 공백을 제거합니다. | `trim("  a  ")`는 `"a"`를 반환합니다. |
| `upper(s:str):str` | 문자열을 대문자로 변환합니다. | `upper("aBc")`는 `"ABC"`를 반환합니다. |
| `lower(s:str):str` | 문자열을 소문자로 변환합니다. | `lower("aBc")`는 `"abc"`를 반환합니다. |
| `replace(s:str, old:str, new:str):str` | `old`를 `new`로 교체합니다. | `replace("a-b", "-", "/")`는 `"a/b"`를 반환합니다. |
| `contains(s:str, sub:str):bool` | `sub` 문자열 포함 여부를 확인합니다. | `contains("hello", "ell")`는 `true`를 반환합니다. |

### 6.5. 수학 (Math)

| 함수 | 설명 | 예시 |
| --- | --- | --- |
| `abs(n)` | 숫자의 절댓값을 반환합니다. (결과는 `float`) | `abs(-5)`는 `5.0`을 반환합니다. |
| `sqrt(n)` | 숫자의 제곱근을 반환합니다. (결과는 `float`) | `sqrt(16)`은 `4.0`을 반환합니다. |
| `pow(base, exp)` | `base`의 `exp` 거듭제곱을 반환합니다. | `pow(2, 3)`은 `8.0`을 반환합니다. |
| `sin(n)` | 숫자의 사인(sine) 값을 반환합니다. | `sin(0)`은 `0.0`을 반환합니다. |
| `cos(n)` | 숫자의 코사인(cosine) 값을 반환합니다. | `cos(0)`은 `1.0`을 반환합니다. |
| `tan(n)` | 숫자의 탄젠트(tangent) 값을 반환합니다. | `tan(0)`은 `0.0`을 반환합니다. |

## 7. 데모 프로그램

### 7.1. Hello World

가장 간단한 형태의 Duet 프로그램입니다.

```duet
cons main -> print("Hello, World!")

main
```

### 7.2. 파일 내용 대문자로 변환

`input.txt` 파일을 읽어 대문자로 변경한 후, `output.txt` 파일에 저장합니다.

```duet
proc to_upper(s:str):str -> upper(s) 

supp read_file:str -> read("input.txt")
cons write_file(content:str) -> write("output.txt", content)

read_file |> to_upper |> write_file
```
