package br.com.wallet.kafka

import br.com.wallet.entity.Account
import br.com.wallet.json
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers.IO
import kotlinx.coroutines.launch
import org.apache.kafka.clients.consumer.KafkaConsumer
import org.apache.log4j.BasicConfigurator
import org.apache.log4j.Level
import org.apache.log4j.LogManager
import org.apache.log4j.Logger
import org.apache.log4j.varia.NullAppender
import java.time.Duration.ofMillis

private val consumerProps = mapOf(
    "bootstrap.servers" to "kafka:29092",
    "auto.offset.reset" to "earliest",
    "key.deserializer" to "org.apache.kafka.common.serialization.StringDeserializer",
    "value.deserializer" to "org.apache.kafka.common.serialization.ByteArrayDeserializer",
    "group.id" to "wallet",
    "security.protocol" to "PLAINTEXT"
)

fun startKafka(
    onMessage: suspend (Message) -> Unit
) {
    CoroutineScope(IO).launch {
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
                        val value = String(record.value())
                        println("Kafka: $value")
                        if (value.contains("BalanceUpdated").not())
                            return@forEach

                        try {
                            val payload = json
                                .decodeFromString<BalanceUpdated>(
                                    value
                                ).payload
                            onMessage(
                                Message(
                                    originAccount = Account(
                                        balance = payload.originAccountBalance,
                                        id = payload.originAccountId
                                    ),
                                    destinationAccount = Account(
                                        balance = payload.destinationAccountBalance,
                                        id = payload.destinationAccountId
                                    )
                                )
                            )
                        } catch (t: Throwable) {
                            println("Kafka.error: ${t.message}")
                            t.printStackTrace()
                        }
                    }
            }
        }
    }
}
