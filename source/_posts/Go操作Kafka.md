---
type: Post
title: Go操作Kafka
tags: Go
category: 开发
category_bar: true
abbrlink: 25901
date: 2024-04-20 23:24:35
---

Kafka是一种高吞吐量的分布式发布订阅消息系统，它可以处理消费者规模的网站中的所有动作流数据，具有高性能、持久化、多副本备份、横向扩展等特点。

首先来看几个概念：

1. **消息队列**: Kafka 通过消息队列的方式来处理数据流。生产者将消息发布到 Kafka 集群中的主题（topic）中，消费者订阅这些主题并处理消息。这种解耦的模式使得生产者和消费者之间可以独立操作，从而提高了系统的可伸缩性和灵活性。
2. **分布式存储**: Kafka 使用分布式存储来保存消息。消息被分成多个分区（partition），并分布在 Kafka 集群的多个节点上，以实现水平扩展和高可用性。
3. **流处理**: Kafka 提供了一套流处理 API，允许开发人员在数据流中进行实时处理和转换。这使得用户能够构建复杂的流处理应用程序，例如实时数据分析、事件驱动的应用程序等。
4. **持久性**: Kafka 的消息被持久化在磁盘上，因此即使消费者下线或发生故障，消息仍然可以被保留和重新处理
5. **Broker**: Kafka 集群中的每个服务器节点称为 Broker。每个 Broker 存储着一个或多个主题（topics）的消息数据，并且负责消息的存储和转发。
6. **Topic**: 主题是 Kafka 中的基本数据单元。它是一个逻辑上的概念，用于分类消息。生产者（Producers）发布消息到主题，而消费者（Consumers）从主题订阅消息。
7. **Partition**: 主题可以分成多个分区。每个分区是一个有序的消息队列，其中的消息被分配到特定的顺序中。分区使得 Kafka 集群能够水平扩展，因为每个分区可以分布在不同的 Broker 上，从而实现负载均衡和高可用性。
8. **Producer**: 生产者是负责将消息发布到 Kafka 主题的应用程序。生产者将消息发送到指定的主题，然后 Kafka 集群将消息存储在相应的分区中。
9. **Consumer**: 消费者是订阅 Kafka 主题并处理消息的应用程序。消费者从指定的主题中读取消息，并根据业务逻辑进行处理。消费者可以以不同的方式组织，例如消费者组（Consumer Group），它们可以并行地处理消息以实现负载均衡和容错性。

Go社区中目前有三个比较常用的kafka客户端库 , 它们各有特点。首先是[IBM/sarama](https://github.com/IBM/sarama)（这个库已经由Shopify转给了IBM）。相较于sarama， [kafka-go](https://github.com/segmentio/kafka-go) 更简单、更易用。[segmentio/kafka-go](https://github.com/segmentio/kafka-go) 是纯Go实现，提供了与kafka交互的低级别和高级别两套API，同时也支持Context。此外社区中另一个比较常用的[confluentinc/confluent-kafka-go](https://github.com/confluentinc/confluent-kafka-go)，它是一个基于cgo的[librdkafka](https://github.com/edenhill/librdkafka)包装，在项目中使用它会引入对C库的依赖。

本文主要介绍sarama的使用。

## Sarama

go语言中连接kafka使用第三方库：[github.com/IBM/sarama](https://github.com/IBM/sarama) 。

### 下载及安装

```go
go get github.com/IBM/sarama
```

### 注意事项

`sarama` v1.20之后的版本加入了`zstd`压缩算法，需要用到cgo，在Windows平台编译时会提示类似如下错误：

```go
exec: "gcc":executable file not found in %PATH%
```

所以在Windows平台请使用v1.19版本的sarama。

## 连接kafka发送消息

```go
package main

import (
    "fmt"

    "github.com/IBM/sarama"
)

// 基于sarama第三方库开发的kafka client

func main() {
    config := sarama.NewConfig()
    config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
    config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
    config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

    // 构造一个消息
    msg := &sarama.ProducerMessage{}
    msg.Topic = "web_log"
    msg.Value = sarama.StringEncoder("this is a test log")
    // 连接kafka
    client, err := sarama.NewSyncProducer([]string{"192.168.1.7:9092"}, config)
    if err != nil {
        fmt.Println("producer closed, err:", err)
        return
    }
    defer client.Close()
    // 发送消息
    pid, offset, err := client.SendMessage(msg)
    if err != nil {
        fmt.Println("send msg failed, err:", err)
        return
    }
    fmt.Printf("pid:%v offset:%v\n", pid, offset)
}
```

## 连接kafka消费信息

```go
package main

import (
    "fmt"

    "github.com/IBM/sarama"
)

// kafka consumer

func main() {
    consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
    if err != nil {
        fmt.Printf("fail to start consumer, err:%v\n", err)
        return
    }
    partitionList, err := consumer.Partitions("web_log") // 根据topic取到所有的分区
    if err != nil {
        fmt.Printf("fail to get list of partition:err%v\n", err)
        return
    }
    fmt.Println(partitionList)
    for partition := range partitionList { // 遍历所有的分区
        // 针对每个分区创建一个对应的分区消费者
        pc, err := consumer.ConsumePartition("web_log", int32(partition), sarama.OffsetNewest)
        if err != nil {
            fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
            return
        }
        defer pc.AsyncClose()
        // 异步从每个分区消费信息
        go func(sarama.PartitionConsumer) {
            for msg := range pc.Messages() {
                fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, msg.Value)
            }
        }(pc)
    }
}
```
