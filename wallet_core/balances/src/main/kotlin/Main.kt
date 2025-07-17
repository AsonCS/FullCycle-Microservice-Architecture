package org.example

import org.example.kafka.startKafka


suspend fun main() {
    println("Hello World!")
    startKafka { message ->
        println(message)
    }
}
