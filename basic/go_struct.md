# go 结构体

- [go 结构体](#go-结构体)
  - [定义结构体](#定义结构体)
  - [访问结构体成员变量](#访问结构体成员变量)
  - [结构体的匿名字段/嵌入类型](#结构体的匿名字段嵌入类型)
  - [匿名结构体类型](#匿名结构体类型)
  - [结构体类型的别名](#结构体类型的别名)

## 定义结构体

- `struct` 是域的集合
- 定义结构体需要使用 type 和 struct 关键字

  ```go
  package main

  import "fmt"

  type vertex struct {
      X int
      Y int
  }

  func main() {
      fmt.Println(vertex{1, 2})
  }
  ```

- 声明结构体变量：结构体字面量代表使用列举的域给新分配的结构体赋值，即匿名结构体
  - `var_name := struct_name {var1, var2...,varn}`
  - `var_name := struct_name {key1 : var1, key2 : val2..., keyn : varn}`
    - 使用 `key:` 可以仅列出部分字段，与字段名顺序无关

  ```go
  package main

  import "fmt"

  type vertex struct {
      X, Y int
  }

  var (
      v1 = vertex{1, 2}
      v2 = vertex{X : 1}
      v3 = vertex{}
      p = &vertex{2, 3}
  )

  func main() {
      fmt.Println(v1, v2, v3, p) //{1 2} {1 0} {0 0} &{2 3}
  }
  ```

## 访问结构体成员变量

- 访问结构体成员变量用 `.` 操作符
- 和 C++ 不一样，结构体指针访问结构体成员变量也用 `.` 操作符

  ```go
  package main

  import "fmt"

  type vertex struct {
      X int
      Y int
  }

  func main() {
      v := vertex{1, 2}
      fmt.Println(v)
      p := &v
      p.X = 1e9
      fmt.Println(v)
  }
  ```

- 结构体作为函数参数

## 结构体的匿名字段/嵌入类型

结构体的匿名字段指的是只有类型而没有名称的字段。匿名字段的类型必须是一个数据类型的名称或者一个与非接口类型对应的指针类型的名称。

代表匿名字段类型的非限定名称将被隐含地作为对该字段的名称。这些名称不能与所属结构体的其他字段名称重复。

```go
type A struct {
  T1
  *T2
  P.T3
  *P.T4
}

// B 包含 T1/T2/T3/T4/len 字段
type B struct {
  A

  len int
}
```

嵌入类型 A 附带的方法会与被嵌入的结构体类型 B 关联。但是调用这些方法的时候，实际上会自动转发到这个嵌入类型在 A 的值，用该值去调用调用方法。

但是，如果 B 中包含了一个和 A 同名(签名相同或者不同)的方法，那么对 B 调用同名方法调用的是 B 的方法，A 的方法被隐藏。但是可以通过 B.A 调用 A 的同名方法。即：

- 可以在被嵌入的结构体类型的值上像调用自己的字段或方法那样调用任意深度的嵌入类型值的字段或方法。唯一的前提条件是这些嵌入类型的字段或方法没有被隐藏。如果被隐藏，可以通过链式的选择表达式或调用表达式访问或调用
- 被嵌入的结构体类型的字段或方法可以隐藏任意深度的嵌入类型的同名字段或方法。这包括任何较浅层次的嵌入类型的字段或方法都会隐藏较深层次的嵌入类型包含的字段或方法。**注意：**这种隐藏是交叉的，字段可以隐藏方法，方法也可以隐藏字段，只要名称相同。因此，如果同意嵌入层次的两个嵌入类型拥有同名的字段或方法，设计它们的选择表达式或调用式会造成一个编译错误。

假设有结构体类型 B 和 非指针类型的数据类型 A：

- 如果 B 中嵌入的类型是 A：那么 B 和 \*B 的方法集合中都会包含接收者类型是 A 的方法。此外，\*B 的方法中还包含接收者类型是 \*A 的方法
- 如果 B 中嵌入的类型是 \*A：那么 B 和 \*B 的方法集合中都会包含接收者类型是 A 或 \*A 的方法

Go 中只存在嵌入而不存在继承的概念，即被嵌入的结构体的值不能赋给嵌入的结构体类型的值。

嵌入字段能够存储所有该字段相关的字段和方法。如果嵌入字段是一个接口，嵌入字段还可以存储所有实现了该接口类型的数据类型的值。

```go
package main

import "fmt"

type notifier interface {
  notify()
}

type people interface {
  sayHello()
}

type user struct {
  name, email string
}

func (u *user) notify() {
  fmt.Printf("Sending user email to %s<%s>\n", u.name, u.email)
}

func (u user) sayHello() {
  fmt.Printf("Hi, I'm user %s<%s>\n", u.name, u.email)
}

type admin struct {
  user
  level string
}

type admin2 struct {
  *user
  level string
}

func sendNotification(n notifier) {
  n.notify()
}

func peopleSayHello(p people) {
  p.sayHello()
}

func main() {
  adm := admin{
    user: user{
      name:  "kiki",
      email: "kiki@bmi-tech.cn",
    },
    level: "super",
  }
  // sendNotification(adm) //annot use adm (variable of type admin) as notifier value in argument to sendNotification: missing method notify
  peopleSayHello(adm)
  sendNotification(&adm)
  peopleSayHello(&adm)

  adm2 := admin2{
    user: &user{
      name:  "kiki",
      email: "kiki@bmi-tech.cn",
    },
    level: "super",
  }
  sendNotification(adm2)
  peopleSayHello(adm2)
  sendNotification(&adm2)
  peopleSayHello(&adm2)

  return
}
```

## 匿名结构体类型

匿名结构体类型比匿名结构体类型少了关键字 type 和 类型名称：

```go
struct {
  len int
  cap int
}
```

可以在数组、切片或 map 类型的声明中，将一个匿名结构体类型作为它们的元素类型。

也可以将匿名结构体类型作为一个变量的类型。

```go
a := struct {
  len int
  cap int
}{1, 2}
```

## 结构体类型的别名

在一个结构体类型的别名类型的值上，不能调用该结构体类型的方法，也不能调用结构体类型指针的方法。因为别名类型和源类型的内部结构相同，所以可以先转成成源类型再进行调用。

结构体类型的别名类型拥有源类型的所有字段。只是不包含源结构体关联的方法。

```go
package main

import (
  "fmt"
)

type member1 struct {
  name string
  age  int
}

type member2 struct {
  name string
  age  int
}

type member3 struct {
  age  int
  name string
}

type member4 struct {
  myname string
  age    int
}

func main() {
  var (
    m1 member1
    m2 member2
    m3 member3
    m4 member4
  )

  m1 = member1{"kiki", 28}
  fmt.Println(m1) //{kiki 28}

  // m2 = m1 //cannot use m1 (variable of type member1) as member2 value in assignment
  m2 = member2(m1)
  fmt.Println(m2) //{kiki 28}

  // m3 = m1          //cannot use m1 (variable of type member1) as member3 value in assignment
  // m3 = member3(m1) //cannot convert m1 (variable of type member1) to member3
  fmt.Println(m3) //{0 }

  // m4 = m1          //cannot use m1 (variable of type member1) as member4 value in assignment
  // m4 = member4(m1) //cannot convert m1 (variable of type member1) to member4
  fmt.Println(m4) //{ 0}

  return
}

```
