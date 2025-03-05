---
type: Post
title: mutsumi的排列连通
tags: 题解
category: 算法
category_bar: true
abbrlink: 8656
date: 2024-03-08 08:23:23
---

链接：<https://ac.nowcoder.com/acm/contest/67745/M>

---

## 题目描述

mutsumi有两个排列，放置在一个$2×n$的矩形中，每次你可以选择一个数字 $x$，将两个排列内的$x$所在的单元格删除。

mutsumi想删除尽可能少的数字，使得矩形至少被分成两个连通块（不一定是矩形），请输出最小的删除次数。若无法通过删除使得矩形被分成至少两个连通块，则输出 -1。

连通块：块内任意两点（可以是同一点）都可以找到至少一条只由上下左右组成的路径相连。

长度为$n$的排列为：$1$$-$$n$中，每个数字恰好出现一次。

### **输入描述:**

第一行输入一个整数$T(1 \leq T \leq 10^5)$表示测试数据组数。

每组测试数据的第一行输入一个整数$n(1 \leq n \leq 10^5)$表示排列长度。

每组测试数据的第二行输入$n$个整数表示第一个排列$a(1 \leq a_i \leq n)$。

每组测试数据的第三行输入$n$个整数表示第二个排列$b(1 \leq b_i \leq n)$。

$n$的总和不超过$10^5$。

### **输出描述:**

输出一个整数表示答案。无解则输出$-1$。

### 示例1

#### 输入

```Plain text
2
1
1
1
3
1 2 3
1 2 3
```

#### 输出

```Plain text
1
1
```

### 说明

第一组数据中，无法通过删除使得矩形至少被分成两个连通块，因此输出$-1$。第二组数据中，删除数字$2$，即可将矩形分成两个连通块。

---

## 题解

因为只有两行，所以结果只有三种，-1、1或2。按照题意，我们只需要对其进行分类讨论即可。

- 排列长度等于2：

    (1)如果$a_1= b_2$，此时结果为1。

    (2)如果$a_1\neq b_2$，此时无解，输出-1。

- 排列长度大于2：

    (1)存在$a_i=b_i$，且$i\neq 1,n$，则操作一次即可，输出1。

    (2)存在$a_i=b_{i+1}$或者$a_i=b_{i-1}$，且均不越界，则操作一次即可，输出1。

    (3)不符合上述情况，则需要操作两次，输出2。

- 排列长度等于2，即长度为1：

    (1)无解，输出-1。

完整代码如下：

```cpp
#include<bits/stdc++.h>
using namespace std;

int a[100005]={0};
int b[100005]={0};

void solve()
{
    memset(a,0,sizeof(a));
    memset(b,0,sizeof(b));
    int n;
    cin>>n;
    for (int i=1;i<=n;i++)
    {
        cin>>a[i];
    }
    for (int i=1;i<=n;i++)
    {
        cin>>b[i];
    }

    if (n==1)
    {
        cout<<"-1"<<endl;
    }
    else if (n==2)
    {
        if (a[1]==b[2])
        {
            cout<<"1"<<endl;
        }
        else
        {
            cout<<"-1"<<endl;
        }
    }
    else
    {
        for (int i=2;i<n;i++)
        {
            if (a[i]==b[i] || a[i]==b[i-1] || a[i]==b[i+1])
            {
                cout<<"1"<<endl;
                return ;
            }
        }
        if (a[1]==b[2] || a[n]==b[n-1])
        {
            cout<<"1"<<endl;
        }
        else
        {
            cout<<"2"<<endl;
        }

    }
    return ;
}

int main()
{
    ios::sync_with_stdio(false);
    cin.tie(0);
    cout.tie(0);

    int t;
    cin>>t;
    while (t)
    {
        t--;
        solve();
    }
    return 0;
}
```
