package org.example

import kotlinx.coroutines.Dispatchers.IO
import kotlinx.coroutines.withContext
import org.apache.kafka.clients.consumer.KafkaConsumer
import org.apache.log4j.BasicConfigurator
import org.apache.log4j.Level
import org.apache.log4j.LogManager
import org.apache.log4j.Logger
import org.apache.log4j.varia.NullAppender
import java.time.Duration.ofMillis

private val consumerProps = mapOf(
    "bootstrap.servers" to "localhost:9092",
    "auto.offset.reset" to "earliest",
    "key.deserializer" to "org.apache.kafka.common.serialization.StringDeserializer",
    "value.deserializer" to "org.apache.kafka.common.serialization.ByteArrayDeserializer",
    "group.id" to "wallet",
    "security.protocol" to "PLAINTEXT"
)

suspend fun startKafka(
    onMessage: (String) -> Unit
) = withContext(IO) {
    BasicConfigurator.configure()
    LogManager.getRootLogger().level = Level.OFF
    Logger.getRootLogger().addAppender(NullAppender())

    val consumer = KafkaConsumer<String, ByteArray>(
        consumerProps
    ).apply {
        subscribe(
            listOf(
                "balances",
                "transactions"
            )
        )
    }

    consumer.use {
        while (true) {
            consumer
                .poll(ofMillis(400))
                .forEach { record ->
                    onMessage(
                        String(
                            record.value()
                        )
                    )
                }
        }
    }
}
