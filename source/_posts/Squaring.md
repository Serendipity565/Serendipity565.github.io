---
type: Post
title: Squaring
tags: 题解
category: 算法
category_bar: true
abbrlink: 27956
date: 2024-07-24 18:32:43
---

题目链接：

[Problem - C - Codeforces](https://codeforces.com/contest/1995/problem/C)

---

## 题目描述

ikrpprpp 发现了一个由整数组成的数组 $a$ 。他喜欢公平，所以想让 $a$ 变得公平，也就是让它不递减

为此，他可以对数组中的索引 $1≤i≤n$ 进行公正操作，将 $a_i$ 替换为 $a_i^2$ （位置 $i$ 的元素及其平方）。例如，如果 $a=[2,4,3,3,5,3]$ ，ikrpprpp 选择对 $i=4$ 执行正义行动，那么 $a$ 就会变成 $[2,4,3,9,5,3]$ 。

要使数组不递减，最少需要多少次正义行动？

### 输入格式

第一行包含一个整数 $t ( 1≤t≤1000 )$ - 测试用例的数量。随后是测试用例的描述。

对于每个测试用例，第一行包含一个整数 $n$ - 数组 $a$ 的大小。第二行包含 $n ( 1≤n≤2*10^5 )$ 个整数 $a1,a2,…,an（1≤a_i≤10^6）$。

所有测试用例中 $n$ 的总和不超过 $2*10^5$ 。

### 输出格式

对于每个测试用例，打印一个整数--使数组 $a$ 不递减所需的最小正义行为数。如果无法做到，则打印$−1$ 。

### 样例 #1

#### 样例输入 #1

```Plain text
7
3
1 2 3
2
3 2
3
3 1 5
4
1 1 2 3
3
4 3 2
9
16 2 4 2 256 2 4 2 8
11
10010 10009 10008 10007 10006 10005 10004 10003 10002 10001 10000
```

#### 样例输出 #1

```Plain text
0
1
-1
0
3
15
55
```

### 注释

在第一个测试案例中，无需执行正义行为。阵列本身就是公平的！

在第三个测试案例中，可以证明数组不可能非递减。

在第五个测试用例中，ikrpprppp 可以在索引 $3$ 上执行一次正义行动，然后在索引 $2$ 上执行一次正义行动，最后在索引 $3$ 上执行另一次正义行动。之后， $a$ 将变成 $[4,9,16]$ 。

---

## 题解

观察数据范围，我们可以暴力求解次数。但是这样会面临一个问题，爆long long而导致答案错误。

我们先来观察相邻的两个数。我们假定 $a,b$ 是两个相邻的数，如果 $a>b$，我们需要对 $b$ 操作。假设 $a$ 是数组的第一个元素，我们只需要找到一个 $n$ ，使得 $b^{n-1}<a≤b^n$ 。那么如果 $a$ 并不是数组的第一个元素呢，我们假设到 $a$ 时要满足条件需要 $a$ 进行了 $m$ 次操作，即 $a^m$，此时，由于幂函数的单调性可知 $b^{(n-1)+m}<a^m≤b^{n+m}$ 依然成立。也就是说到 $a,b$ 时，$b$ 需要操作 $m+n$ 次。

所以，我们只需要遍历数组，记录下前 $i-1$ 个数时的操作次数，即可得到第 $i$ 个数的操作次数，然后累加即可。

下面是 ac 代码：

```cpp
#include <bits/stdc++.h>
using namespace std;
#define endl '\n'
typedef long long ll;

ll a[200005] = {0};
ll b[200005] = {0};

void debug()
{
    return;
}

void solve()
{
    int n;
    cin >> n;
    ll mymin = 0;
    for (int i = 1; i <= n; i++)
    {
        cin >> a[i];
    }
    for (int i = 2; i <= n; i++)
    {
        if (a[i] == 1 && a[i - 1] != 1)
        {
            cout << -1 << endl;
            return;
        }
        else
        {
            ll p = a[i - 1];
            ll now = a[i];
            ll extra = 0;
            while (p * p <= now)
            {
                extra -= 1, p *= p;
            }
            while (now < p)
            {
                extra++, now *= now;
            }

            b[i] = max(0ll, b[i - 1] + extra);
        }
    }
    ll ans = 0;
    for (int i = 2; i <= n; i++)
    {
        ans += b[i];
    }
    cout << ans << endl;
    return;
}

int main()
{
    ios::sync_with_stdio(false);
    cin.tie(0);
    cout.tie(0);

    int t;
    t = 1;
    cin >> t;
    while (t--)
    {
        solve();
    }
    return 0;
}
```
