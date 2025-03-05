---
type: Post
title: trie树
tags: 算法
category: 算法
category_bar: true
abbrlink: 19236
date: 2024-08-06 23:23:23
---

## 定义

字典树，英文名 trie。顾名思义，就是一个像字典一样的树。

![trie1.png](/img/blog/tries/1.png)

可以发现，这棵字典树用边来代表字母，而从根结点到树上某一结点的路径就代表了一个字符串。

举个例子，$1→4→8→12$ 表示的就是字符串 `caa`。

## 实现

```cpp
int getnum(char x)
{
    if (x >= 'A' && x <= 'Z')
    {
        return x - 'A';
    }
    else if (x >= 'a' && x <= 'z')
    {
        return x - 'a' + 26;
    }
    else
    {
        return x - '0' + 52;
    }
}

void insert(string s)
{
    int p = 0, len = s.size();
    for (int i = 0; i < len; i++)
    {
        int c = getnum(s[i]);
        if (!trie[p][c])
        {
            trie[p][c] = ++cnt;
        }
        p = trie[p][c];
        exist[p]++;
    }
}

int find(string s)
{
    int p = 0, len = s.size();
    for (int i = 0; i < len; i++)
    {
        int c = getnum(s[i]);
        if (!trie[p][c])
        {
            return 0;
        }
        p = trie[p][c];
    }
    return exist[p];
}
```

## 应用

### 检索字符串

字典树最基础的应用——查找一个字符串是否在 “字典” 中出现过。

[于是他错误的点名开始了 - 洛谷](https://www.luogu.com.cn/problem/P2580)

```cpp
#include <bits/stdc++.h>
using namespace std;
#define endl '\n'
typedef long long ll;
typedef long double ld;

int q, n, nex[500010][65], exist[500010], cnt;
string s;

int getnum(char x)
{
    if (x >= 'A' && x <= 'Z')
    {
        return x - 'A';
    }
    else if (x >= 'a' && x <= 'z')
    {
        return x - 'a' + 26;
    }
    else
    {
        return x - '0' + 52;
    }
}

void insert(string s)
{
    int p = 0, len = s.size();
    for (int i = 0; i < len; i++)
    {
        int c = getnum(s[i]);
        if (!nex[p][c])
        {
            nex[p][c] = ++cnt;
        }
        p = nex[p][c];
    }
    exist[p]++;
}

int find(string s)
{
    int p = 0, len = s.size();
    for (int i = 0; i < len; i++)
    {
        int c = getnum(s[i]);
        if (!nex[p][c])
        {
            return 0;
        }
        p = nex[p][c];
    }
    return exist[p];
}

void solve()
{
    cin >> n;
    for (int i = 1; i <= n; i++)
    {
        cin >> s;
        insert(s);
    }
    cin >> q;
    map<string, int> mp;
    for (int i = 1; i <= q; i++)
    {
        cin >> s;
        int k = find(s);
        if (mp.find(s) != mp.end())
        {
            cout << "REPEAT" << endl;
        }
        else if (k == 1)
        {
            cout << "OK" << endl;
            mp[s]++;
        }
        else if (k == 0)
        {
            cout << "WRONG" << endl;
        }
    }
    return;
}

int main()
{
    ios::sync_with_stdio(false);
    cin.tie(0);
    cout.tie(0);

    int t;
    t = 1;
    // cin >> t;
    while (t--)
    {
        solve();
    }
    return 0;
}
```

### **维护异或极值**

01-trie 是指字符集为 $\{ 0,1\}$ 的 trie。01-trie 可以用来维护一些数字的异或和，支持修改(删除，重新插入)，和全局加一 (即让其所维护所有数值递增1，本质上是一种特殊的修改操作)。

如果要维护异或和，需要按值从低位到高位建立 trie。

[最大异或对 The XOR Largest Pair - 洛谷](https://www.luogu.com.cn/problem/P10471)

```cpp
#include <bits/stdc++.h>
using namespace std;
#define endl '\n'
typedef long long ll;
typedef long double ld;

int n;
int nex[3000010][3], a[100010], cnt;

void insert(int num)
{
    int p = 0;
    for (int i = 30; i >= 0; i--)
    {
        int c = (num >> i) & 1;
        if (!nex[p][c])
        {
            nex[p][c] = ++cnt;
        }
        p = nex[p][c];
    }
    return;
}

int find(int num)
{
    int p = 0;
    int ans = 0;
    for (int i = 30; i >= 0; i--)
    {
        int c = (num >> i) & 1;
        if (nex[p][!c])
        {
            ans = ans * 2 + 1;
            p = nex[p][!c];
        }
        else
        {
            ans *= 2;
            p = nex[p][c];
        }
    }
    return ans;
}

void solve()
{
    cin >> n;
    for (int i = 1; i <= n; i++)
    {
        cin >> a[i];
        insert(a[i]);
    }
    int mymax = 0;
    for (int i = 1; i <= n; i++)
    {
        mymax = max(find(a[i]), mymax);
    }
    cout << mymax << endl;
    return;
}

int main()
{
    ios::sync_with_stdio(false);
    cin.tie(0);
    cout.tie(0);

    int t;
    t = 1;
    // cin >> t;
    while (t--)
    {
        solve();
    }
    return 0;
}
```
