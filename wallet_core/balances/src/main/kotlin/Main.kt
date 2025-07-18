package br.com.wallet

import br.com.wallet.api.configureBalances
import br.com.wallet.api.configureRouting
import br.com.wallet.api.configureSerialization
import br.com.wallet.database.AccountService
import br.com.wallet.database.getDatabase
import br.com.wallet.gateway.AccountGateway
import br.com.wallet.kafka.startKafka
import io.ktor.server.application.*


fun main(args: Array<String>) {
    io.ktor.server.netty.EngineMain.main(args)
}

fun Application.module() {
    val database = getDatabase()
    val gateway: AccountGateway = AccountService(database)
    startKafka { message ->
        gateway.upsert(
            message.originAccount,
            message.destinationAccount
        )
        println(message)
    }

    configureSerialization()
    configureBalances(
        gateway
    )
    configureRouting()
}
