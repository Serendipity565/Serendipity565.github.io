---
type: Post
title: 深入理解C语言指针
tags: C/C++
category: 开发
category_bar: true
abbrlink: 36659
date: 2024-06-15 23:23:32
---

指针是C语言中一个非常强大且重要的概念，它不仅能够提供直接的内存访问，还能用于实现许多高级的数据结构和算法。然而，指针的概念相对复杂，新手程序员常常感到困惑。本文将深入探讨C语言中的指针，从基本概念到高级应用，帮助你全面理解指针的使用。

## 引入、什么是内存地址

内存地址是计算机系统用来访问内存中某个特定存储单元的标识符。每个存储单元都有一个唯一的地址，就像每个房子都有一个唯一的门牌号码。

内存地址的范围取决于系统的位数：

- 在32位系统中，内存地址通常是32位长，地址范围为0到`0xFFFFFFFF`。
- 在64位系统中，内存地址通常是64位长，地址范围为0到`0xFFFFFFFFFFFFFFFF`。

一个内存地址代表一个字节（8bit）的存储空间。

例如：假设’不’的二进制位`0010101110011001`

!["不期而遇"再计算机内存中的地址](/img/blog/srjlcyyzz/1.png)

!['不'的真值与地址的关系](/img/blog/srjlcyyzz/2.png)

## 一、指针的基本概念

### 1.1 什么是指针

指针是一个变量，它存储另一个变量的内存地址。通过指针，我们可以直接访问和修改内存中的数据。

```c
int a = 10;
int *p = &a; // p是一个指向int类型的指针，它存储了变量a的地址
```

在上面的代码中，`&a`表示取变量`a`的地址，`p`是一个指针变量，它保存了这个地址。

### 1.2 指针的声明和初始化

指针变量的声明需要指定它所指向的数据类型，并使用`*`符号。例如：

```c
int *p; // 声明一个指向int类型的指针
```

指针的初始化可以通过将一个变量的地址赋值给它：

```c
int a = 10;
int *p = &a; // p现在指向变量a的地址
```

### 1.3 指针的解引用

通过解引用操作，我们可以访问指针所指向的变量。解引用操作符是`*` ，即取值符，取一个地址下储存的值。

```c
int a = 10;
int *p = &a;
printf("%d\\n", *p); // 输出10
*p = 20; // 修改a的值为20
printf("%d\\n", a); // 输出20
```

### 1.4 指针的加减运算

指针加减运算指的是对指针进行算术运算，从而改变指针所指向的内存地址。这种运算的结果依赖于指针指向的数据类型，因为指针的加减运算是按数据类型的大小来进行的。

指针加法运算是将一个指针加上一个整数。加法运算后，指针会移动到相应位置。

例如，如果`p`是一个指向`int`类型的指针，则`p + 1`会使指针向前移动一个`int`的大小（通常是4个字节）。

```c
int main() {
    int arr[5] = {1, 2, 3, 4, 5};
    int *p = arr; // 指向数组的第一个元素

    printf("Address of p: %p\n", (void *)p); // 输出p的地址
    printf("Value at p: %d\n", *p); // 输出p指向的值，即1

    p = p + 1; // 指针加法运算，p移动到下一个元素
    printf("Address of p after p + 1: %p\n", (void *)p); // 输出p的地址
    printf("Value at p after p + 1: %d\n", *p); // 输出p指向的值，即2

    return 0;
}
```

在这个例子中，`p + 1`使指针`p`移动到数组的下一个元素，即`arr[1]`。

指针减法运算是将一个指针减去一个整数。减法运算后，指针会移动到相应位置。与加法类似。

**注意**：指针的加减运算是按指针指向的数据类型大小进行的。即`p + 1`移动的字节数取决于指针指向的数据类型大小。

### 1.5 指针的差值

指针的差值运算可以计算两个指针之间的距离（即它们之间的元素个数）。这种运算通常用于数组操作。

```c
int main() {
    int arr[5] = {1, 2, 3, 4, 5};
    int *p1 = &arr[0]; // 指向数组的第一个元素
    int *p2 = &arr[4]; // 指向数组的第五个元素

    ptrdiff_t diff = p2 - p1; // 计算指针之间的距离
    printf("Difference between p2 and p1: %td\n", diff); // 输出差值，即4

    return 0;
}
```

在这个例子中，`p2 - p1`计算两个指针之间的距离，即数组元素的个数。

**注意：指针相减合法性，** 只有指向同一数组或同一块内存区域的指针才能进行差值运算。

## 二、指针的高级用法

### 2.1 指向指针的指针

指针不仅可以指向基本数据类型，还可以指向另一个指针。这种指针称为“二级指针”或“指向指针的指针”。

```c
int a = 10;
int *p = &a;
int **pp = &p; // pp是一个指向指针p的指针
printf("%d\\n", **pp); // 输出10
```

### 2.2 数组指针

数组指针是指向数组的指针。它存储的是数组的起始地址，可以通过该指针访问数组中的元素

通过数组指针可以访问数组中的元素。由于数组指针指向整个数组，我们需要先解引用指针，然后再访问具体的元素。

```c
int main() {
    int arr[5] = {1, 2, 3, 4, 5};
    int (*p)[5] = &arr; // 声明并初始化数组指针

    // 访问数组中的元素
    for (int i = 0; i < 5; i++) {
        printf("%d ", (*p)[i]); // 通过数组指针访问元素
    }
    return 0;
}
```

数组指针常用于函数参数，使得函数可以处理多维数组。

指向多维数组的情况下，指针的声明和使用方式略有不同。

```c
int main() {
    int arr[3][4] = {
        {1, 2, 3, 4},
        {5, 6, 7, 8},
        {9, 10, 11, 12}
    };
    int (*p)[4] = arr; // 声明并初始化指向二维数组的指针

    // 访问二维数组中的元素
    for (int i = 0; i < 3; i++) {
        for (int j = 0; j < 4; j++) {
            printf("%d ", p[i][j]); // 通过数组指针访问元素
        }
        printf("\n");
    }
    return 0;
}
```

在这个示例中，`p`是一个指向具有4个`int`元素的数组的指针，`p[i][j]`用来访问二维数组`arr`中的元素。

### 2.3 指针数组

指针数组是一种特殊的数组，它的每个元素都是一个指针。指针数组常用于处理字符串数组或二维数组。

下面是一个包含三个字符串的指针数组的声明和初始化示例：

```c
int main() {
    const char *arr[3] = {"Hello", "World", "C programming"};

    // 访问指针数组中的元素
    for (int i = 0; i < 3; i++) {
        printf("%s\n", arr[i]);
    }
    return 0;
}
```

在这个示例中，`arr`是一个包含三个指向`const char`字符串的指针数组。每个指针都指向一个字符串字面量。

### 2.4 指针函数

指针函数是返回指针的函数。它是一个函数，其返回值是一个指针。指针函数可以返回任何类型的指针，比如指向整数、字符、结构体等的指针。

```c
int* getPointer() {
    static int value = 10; // 使用static确保返回的指针在函数外部有效
    return &value;
}

int main() {
    int *ptr = getPointer(); // 调用指针函数，获取指向整数的指针
    printf("Value: %d\n", *ptr); // 解引用指针，输出值

    return 0;
}
```

指针函数常用于**动态内存分配**和**数据结构操作**（例如返回链表、树等数据结构中的某个节点的指针），将在第三四部分介绍。

### 2.5 函数指针

函数指针是指向函数的指针，允许我们动态地调用函数。函数指针在实现回调函数和函数表时非常有用。

```c
void printHello() {
    printf("Hello\\n");
}

void printWorld() {
    printf("World\\n");
}

void (*funcPtr)(); // 声明一个函数指针
funcPtr = printHello; // 将函数指针指向printHello
funcPtr(); // 调用printHello函数，输出Hello
funcPtr = printWorld; // 将函数指针指向printWorld
funcPtr(); // 调用printWorld函数，输出World
```

**注意：调用函数时被调用的那个函数不要加括号。**`funcPtr = printHello;`中`printHello` 不能加括号。**具体原因在第五部分指针的常见错误（野指针问题解决3）中解释。**

**上面示例的用法本人认为并没有什么用，硬要说的话只能说应付考试，函数指针关键的用法在于回调函数。**

#### 2.5.1 回调函数

回调函数是一种通过函数指针传递给另一个函数并在适当时候调用的函数。回调函数广泛用于**事件驱动编程**和**处理异步任务**。

```c
// 声明一个回调函数类型
typedef void (*callback_t)(int);

// 定义一个使用回调函数的函数
void process(callback_t cb, int value) {
    // 在适当时候调用回调函数
    cb(value);
}

// 定义一些回调函数
void print_value(int value) {
    printf("Value: %d\n", value);
}

void double_value(int value) {
    printf("Double value: %d\n", 2 * value);
}

int main() {
    // 使用不同的回调函数
    process(print_value, 5);  // 输出 Value: 5
    process(double_value, 5); // 输出 Double value: 10

    return 0;
}
```

在这个示例中，`callback_t` 是一个指向接受一个 `int` 参数并返回 `void` 的函数的指针类型。

`process` 函数接受一个 `callback_t` 类型的参数 `cb` 和一个 `int` 类型的参数 `value`。在函数内部，`cb(value)` 调用 `cb` 指向的函数，并将 `value` 作为参数传递给该函数。

之后定义了两个回调函数：

- `print_value`：接收一个 `int` 参数并打印它的值。
- `double_value`：接收一个 `int` 参数，计算它的两倍，并打印结果。

在 `main` 函数中，通过 `process` 函数来使用不同的回调函数：

1. `process(print_value, 5);`：调用 `process` 函数，传递 `print_value` 作为回调函数和 `5` 作为参数。`process` 内部调用 `print_value(5)`，输出 `Value: 5`。
2. `process(double_value, 5);`：调用 `process` 函数，传递 `double_value` 作为回调函数和 `5` 作为参数。`process` 内部调用 `double_value(5)`，输出 `Double value: 10`。

## 三、指针与内存管理

### 3.1 动态内存分配

C语言提供了`malloc`、`calloc`和`realloc`等函数，用于动态分配内存。分配的内存需要使用`free`函数释放。

#### 3.1.2 动态内存分配函数

##### `malloc`

`malloc`（memory allocation）函数分配指定大小的内存，并返回一个指向这块内存的指针。分配的内存未被初始化。

```c
void* malloc(size_t size);
```

- `size`：要分配的内存块的大小（以字节为单位）。
- 返回值：成功时，返回指向已分配内存块的指针；失败时，返回`NULL`。

```c
int main() {
    int *ptr;
    int n = 5;

    // 分配内存用于存储5个int类型的元素
    ptr = (int *)malloc(n * sizeof(int));

    // 检查内存分配是否成功
    if (ptr == NULL) {
        printf("Memory allocation failed\n");
        return 1;
    }

    // 使用分配的内存
    for (int i = 0; i < n; i++) {
        ptr[i] = i + 1;
    }

    // 打印数组元素
    for (int i = 0; i < n; i++) {
        printf("%d ", ptr[i]);
    }
    printf("\n");

    // 释放内存
    free(ptr);

    return 0;
}
```

##### `calloc`

`calloc`（contiguous allocation）函数分配指定数量的内存块，每块大小为指定大小，并初始化所有内存块为零。

```c
void* calloc(size_t num, size_t size);
```

- `num`：要分配的元素的数量。
- `size`：每个元素的大小（以字节为单位）。
- 返回值：成功时，返回指向已分配并初始化为零的内存块的指针；失败时，返回`NULL`。

```c
int main() {
    int *ptr;
    int n = 5;

    // 分配内存用于存储5个int类型的元素，并初始化为0
    ptr = (int *)calloc(n, sizeof(int));

    // 检查内存分配是否成功
    if (ptr == NULL) {
        printf("Memory allocation failed\n");
        return 1;
    }

    // 打印数组元素
    for (int i = 0; i < n; i++) {
        printf("%d ", ptr[i]);
    }
    printf("\n");

    // 释放内存
    free(ptr);

    return 0;
}
```

##### `realloc`

`realloc`（reallocation）函数调整之前分配的内存块的大小。它可以扩展或缩小内存块的大小，并返回一个指向新内存块的指针。

```c
void* realloc(void* ptr, size_t size);
```

- `ptr`：指向要重新分配内存的内存块的指针。如果是`NULL`，`realloc` 的行为类似于 `malloc`。
- `size`：新的内存块的大小（以字节为单位）。
- 返回值：成功时，返回指向新内存块的指针；失败时，返回`NULL`，原内存块保持不变。

```c
int main() {
    int *ptr;
    int n = 5;

    // 分配内存用于存储5个int类型的元素
    ptr = (int *)malloc(n * sizeof(int));

    // 检查内存分配是否成功
    if (ptr == NULL) {
        printf("Memory allocation failed\n");
        return 1;
    }

    // 使用分配的内存
    for (int i = 0; i < n; i++) {
        ptr[i] = i + 1;
    }

    // 打印原数组元素
    for (int i = 0; i < n; i++) {
        printf("%d ", ptr[i]);
    }
    printf("\n");

    // 重新分配内存块，扩展为10个int类型的元素
    ptr = (int *)realloc(ptr, 10 * sizeof(int));

    // 检查内存重新分配是否成功
    if (ptr == NULL) {
        printf("Memory reallocation failed\n");
        return 1;
    }

    // 使用扩展的内存
    for (int i = n; i < 10; i++) {
        ptr[i] = i + 1;
    }

    // 打印扩展后的数组元素
    for (int i = 0; i < 10; i++) {
        printf("%d ", ptr[i]);
    }
    printf("\n");

    // 释放内存
    free(ptr);

    return 0;
}
```

##### `free`

`free`函数释放之前分配的动态内存。释放内存后，指针仍然存在，但它指向的内存不再有效，因此通常将指针设为`NULL`以避免悬空指针。

```c
void free(void* ptr);
```

- `ptr`：指向要释放的内存块的指针。如果是`NULL`，`free` 不进行任何操作。

```c
int main() {
    int *ptr;
    int n = 5;

    // 分配内存用于存储5个int类型的元素
    ptr = (int *)malloc(n * sizeof(int));

    // 检查内存分配是否成功
    if (ptr == NULL) {
        printf("Memory allocation failed\n");
        return 1;
    }

    // 使用分配的内存
    for (int i = 0; i < n; i++) {
        ptr[i] = i + 1;
    }

    // 打印数组元素
    for (int i = 0; i < n; i++) {
        printf("%d ", ptr[i]);
    }
    printf("\n");

    // 释放内存
    free(ptr);
    ptr = NULL; // 避免悬空指针

    return 0;
}
```

### 3.1.2 动态内存分配的注意事项

1. **检查内存分配是否成功**：动态内存分配函数在分配失败时返回`NULL`，必须检查返回值以确保内存分配成功。
2. **避免内存泄漏**：确保每个动态分配的内存都使用`free`函数释放，否则会导致内存泄漏。
3. **避免悬空指针**：释放内存后，将指针设为`NULL`，以避免使用已释放的内存。

## 四、指针在数据结构中的应用

指针在实现链表、树、图等数据结构中起着关键作用。以下是一个简单的单链表示例：

```c
#include <stdio.h>
#include <stdlib.h>

struct Node {
    int data;
    struct Node *next;
};

// 创建一个新节点
struct Node* createNode(int data) {
    struct Node *newNode = (struct Node*)malloc(sizeof(struct Node));
    newNode->data = data;
    newNode->next = NULL;
    return newNode;
}

// 打印链表
void printList(struct Node *head) {
    struct Node *temp = head;
    while (temp != NULL) {
        printf("%d -> ", temp->data);
        temp = temp->next;
    }
    printf("NULL\\n");
}

// 主函数
int main() {
    struct Node *head = createNode(1);
    head->next = createNode(2);
    head->next->next = createNode(3);

    printList(head);

    // 释放链表内存
    struct Node *temp;
    while (head != NULL) {
        temp = head;
        head = head->next;
        free(temp);
    }

    return 0;
}
```

## 五、指针的常见错误

### 5.1 野指针

野指针（Dangling Pointer）是指向已经被释放或未分配的内存的指针。使用野指针会导致不可预测的行为，包括程序崩溃和数据损坏。

#### 野指针的常见原因

1. **未初始化的指针**：指针在声明时没有被初始化。
2. **释放后的指针继续使用**：内存释放后，指针依然被使用。
3. **超出作用域的指针**：指针指向的内存在作用域结束后被回收。

##### 1. 未初始化的指针

```c
int main() {
    int *p; // p未初始化
    *p = 10; // 未定义行为，p是野指针
    printf("%d\n", *p);

    return 0;
}
```

##### 2. 释放后的指针继续使用

```c
int main() {
    int *p = (int *)malloc(sizeof(int));
    *p = 10;
    free(p); // 释放内存
    printf("%d\n", *p); // 未定义行为，p是野指针

    return 0;
}
```

##### 3. 超出作用域的指针

局部变量在函数执行完毕后，其内存会被自动回收。如果在函数中返回局部变量的地址，那么该地址在函数返回后就指向一块已经被回收的内存区域，这就导致了野指针的产生。

```c
int* getPointer() {
    int x = 10;
    return &x; // 返回局部变量的地址
}

int main() {
    int *p = getPointer(); // p指向一个已经回收的内存
    printf("%d\n", *p); // 未定义行为，p是野指针

    return 0;
}
```

在上述代码中，函数 `getPointer` 返回了局部变量 `x` 的地址。但在 `getPointer` 函数执行完毕后，局部变量 `x` 的内存就已经被回收，因此 `p` 变成了一个野指针。

#### 如何避免野指针

1. **初始化指针**：声明指针时进行初始化。

    ```c
    int *p = NULL;
    ```

2. **释放后置空**：释放内存后，将指针置空。置空之后可以继续使用（解决野指针问题2）。

    ```c
    free(p);
    p = NULL;
    ```

3. **避免返回局部变量地址（使用动态内存分配）**：函数中不要返回局部变量的地址。（解决野指针问题3）

    ```c
    int* getPointer() {
        int *x = (int *)malloc(sizeof(int));
        *x = 10;
        return x;
    }
    ```

4. **使用静态变量**：使用静态变量可以确保变量在函数结束后依然存在，但要注意静态变量在全局范围内是共享的。（解决野指针问题3）

    ```c
    int* getPointer() {
        static int x = 10; // 使用静态变量
        return &x; // 返回静态变量的地址
    }
    
    int main() {
        int *p = getPointer();
        printf("%d\n", *p); // 合法使用静态变量的地址
    
        return 0;
    }
    ```

5. **使用智能指针**：在C++中使用智能指针（如`std::unique_ptr`和`std::shared_ptr`）来管理动态内存。

## 六、总结

常见指针定义与相关含义速查表

| 定义 | 含义 |
| --- | --- |
| int i ; | 定义整型变量i |
| int *p ; | p为指向整型数据的指针变量 |
| int a[n] ; | 定义含n个元素的整型数组a |
| int *p[n] ; | n个指向整型数据的指针变量组成的指针数组p |
| int (*p) [n] ; | p为指向含n个元素的一维整型数组的指针变量 |
| int f( ) ; | f为返回整型数的丽数 |
| int *p( ) ; | p为返回指针的函数，该指针指向一个整型数据 |
| int (*p)( ) ; | p为指向数的指针变量，该函数返回整型数 |
| int **p ; | p为指针变量，它指向一个指向整型数据的指针变量 |

指针是C语言中一个强大而复杂的特性，通过学习和理解指针的基本概念和高级用法，我们可以更有效地操作内存和实现复杂的数据结构。尽管指针的使用可能会带来一些问题，如空指针引用和内存泄漏，但通过良好的编程习惯和仔细的代码检查，这些问题是可以避免的。希望本文能帮助你更深入地理解C语言中的指针，并在实际编程中灵活运用它们。
