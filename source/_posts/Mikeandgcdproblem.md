---
type: Post
title: Mike and gcd problem
tags: 题解
category: 算法
category_bar: true
abbrlink: 61177
date: 2024-05-11 22:23:45
---

题目链接：

[Problem - 798C - Codeforces](https://codeforces.com/problemset/problem/798/C)

---

## 题目描述

迈克有一个长度为$n$的序列 $A = [a_1, a_2, ..., a_n]$ 。如果序列 $B = [b_1, b_2, ..., b_n]$ 中所有元素的 *gcd* 都大于 1 ，他就认为这个序列 $B = [b_1, b_2, ..., b_n]$ 很美，即$gcd(b_1, b_2, ..., b_n)>1$。

Mike 想改变他的序列，使其更加美观。他可以选择一个索引 $i ( 1 ≤ i < n )$，删除数字 $a_i， a_{i + 1}$ ，并按照这个顺序将数字 $a_1- a_{i + 1}，a_i+a_{i+1}$ 放到它们的位置上。他希望进行尽可能少的运算。如果可能，请找出使数列 $A$ 变美的最少运算次数，或者告诉他不可能这样做。$gcd(b_1, b_2, ..., b_n)$是最大的非负数 $d$，使得 $*b_i*$ 除于 $d$ 为 $0$ $( 1 ≤ i ≤ n )$。

### 输入描述

第一行包含一个整数$n( 2 ≤n≤ 100 000 )$ - 序列 $A$  的长度。

第二行包含 $n$ 个空格分隔的整数  $a_1, a_2, ..., a_n ( 1 ≤ a_i ≤ 10^{9} )$  - 序列 A 的元素。

### 输出描述

如果可以通过执行上述操作使序列 $A$ 变美，则在第一行输出"YES"(不带引号)，否则输出"NO"(不带引号)。

如果答案是"YES"，则输出使序列 $A$ 优美所需的最小步数

### 示例

#### 输入1

```Plain text
2
1 1
```

#### 输出1

```Plain text
YES
1
```

#### 输入2

```Plain text
3
6 2 4
```

#### 输出2

```Plain text
YES
0
```

#### 输入3

```Plain text
2
1 3
```

#### 输出3

```Plain text
YES
1
```

### 说明

在第一个示例中，你只需移动一次，就可以得到序列 [0, 2] ，其中有$gcd(0,2)=2$ 。

在第二个示例中，序列中的 $gcd$ 已经大于 $1$ 。

---

## 题解

这是一道思维题，不管数据怎么样，输出结果都是"YES"。如果 $gcd(b_1, b_2, ..., b_n)>1$，直接输出0，否则 $gcd(b_1, b_2, ..., b_n)=2$。

先来解释一下为什么输出肯定是"YES"：

- 对于两相邻的数$a，b$ ，如果要进行操作，操作一次后变为 a-b，a+b，在操作一次变成 -2b，2a，也就是说，只要操作足够多的次数，最后这个数列$gcd(b_1, b_2, ..., b_n)$ 一定可以到达2，也就是输出"YES"。

下面来证明为什么要操作的话最终的 $gcd(b_1, b_2, ..., b_n)$ 一定是2：

- 将 $a_i$ 和 $a_{i+1}$ 变成 $a_i- a_{i+1}$ 和  $a_i + a_{i+1}$，若操作前其 $gcd$ 为 $x$，那么操作后其 $gcd$ 只可能为 $x$ 或 $2x$，不会再增加其他质因子。

    > 证明：设x=gcd(a_i, a_{i+1})，ai = ax, a_{i+1} = bx
    >

    > 操作后a_i=(a-b)x, a_{i+1} =(a+b)x
    >

    > 则 gcd(a_i, a_{i+1})= gcd((a-b)x, (a+b)x)=x*gcd(a -b, a+b)
    >

    > 设gcd(a-b, a+b)=k, a-b=a’k, a+b=b’k
    >

    > 那么(a-b)+(a+b)=a'k+b’k, (a+b)-(a-b)=b’k-a’k
    >

    > 所以 2a =a'k+b'k, 2b=b’k-a’k
    >

    > 因为 gcd(2a, 2b)=2×gcd(a, b)=2
    >

    > 所以gcd(a'k + b’k, b’k-a’k)=k* gcd(a’ +b’, b’ -a’)=2
    >

    > 可得k=1或k=2
    >

    > 故gcd(a_i, a_{i+1})=x*k=x或2x
    >

下面给出完整代码：

（这次代码写的比较垃圾，建议根据前面的思维证明自己写。）

```cpp
#include <bits/stdc++.h>
using namespace std;
long long a[100005] = {0};

int main()
{
    ios::sync_with_stdio(false);
    cin.tie(0);
    cout.tie(0);
    bool flag = true;
    int n;
    cin >> n;
    for (int i = 1; i <= n; i++)
    {
        cin >> a[i];
    }
    long long temp = __gcd(a[1], a[2]);
    if (temp > 1)
    {
        for (int i = 2; i < n; i++)
        {
            if (__gcd(a[i], a[i + 1]) != 1 && (__gcd(a[i], a[i + 1]) % temp == 0 || temp % __gcd(a[i], a[i + 1]) == 0))
            {
                temp = min(temp, __gcd(a[i], a[i + 1]));
                continue;
            }
            else
            {
                flag = false;
                break;
            }
        }
    }
    cout << "YES" << "\n";
    if (flag && temp > 1)
    {
        cout << 0;
        return 0;
    }
    long long ans = 0;
    for (int i = 1; i < n; i++)
    {
        if (a[i] % 2 == 1 && a[i + 1] % 2 == 1)
        {
            ans++;
            long long b = a[i];
            long long c = a[i + 1];
            a[i] = b - c;
            a[i + 1] = b + c;
        }
        else if (a[i] % 2 == 1 && a[i + 1] % 2 == 0)
        {
            ans += 2;
            long long b = a[i];
            long long c = a[i + 1];
            a[i] = -2 * c;
            a[i + 1] = 2 * b;
        }
    }
    if (a[n] % 2 == 1)
    {
        ans += 2;
    }
    cout << ans;
    return 0;
}
```
