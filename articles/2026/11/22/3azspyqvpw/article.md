---
title: 2026/11/22 勉強したことメモ（Transactionとかギャップロックとか）
date: 2026-11-22T10:11:00.958812+09:00
description: ""
categories: tech
tagas:
    - MySQL
publish: true
---

## RDBMSのTransactionについて復習

DBがデータの信頼性と整合性を保証するための重要な特性であるACID（Atomicity, Consistency, Isolation, Durability）特性を担保するために、Transaction処理がある。

しかし、ACIDのすべてを担保するには、多くのユーザーを抱えるサービスでは不可欠な並列性を犠牲にする必要がある。

そこで、RDBMSには、並列性をあげることを目的として、TransactionのIsolation（独立性）を緩めるオプションが存在する。

その４種類がread uncommited, read commited, repeatable read, serializableである。

この４種類は、read uncommitdに近づくほど独立性が低いが並列性が高く、serializableに近づくほど独立性が高いが並列性が低い　という特性を持つ。

同時にupdateやinsertされるようなシステムを考える場合はIsolationを犠牲にして並列性を上げていく必要がある。そのため、データの特性ごとにどのようなIsolationを担保する必要があるかどうかを判断する必要がある。

## MySQLのギャップロックについて

例えば、以下のようなidがprimary keyのテーブルがあった時、repeatable readのTransactionだった場合は、`SELECT * FROM user WHERE id = 5 FOR UPDATE`とロックがかけられた場合、存在しないid = 5について、他のトランザクションから値をinsertされないようにするためにロックをかける必要がある。

| id | name |
| --- | --- |
| 1 | hanami |
| 2 | fujita |
| 3 | tsukimura |
| 10 | himesaki |

しかし、MySQLのInnoDBの特性上、存在しない行についてロックをかけることはできず、その代わりにrepeatable readでは`id=10`の行と`id=3`の行の間に値が入らないようにロックをかけることをする。これをギャップロックという。

InnoDBはこのギャップロックによって、ロック付きの読み取り（FOR UPDATEなど）でのファントムリード（同じトランザクション内で同じクエリを再実行したとき、別トランザクションのINSERTによって結果に行が増えてしまうこと）を防いでいるが、これを意識していないと、以下のようなケースでデッドロックが起きてしまう。

```SQL
Tx A: SELECT * FROM user WHERE id = 5 FOR UPDATE;
Tx B: SELECT * FROM user WHERE id = 6 FOR UPDATE;
Tx A: INSERT INTO user VALUES (5, ...);
Tx B: INSERT INTO user VALUES (6, ...);
```

これはMySQLのInnoDB仕様で、Postgressでは発生しない。