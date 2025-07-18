@file:OptIn(ExperimentalUuidApi::class)

package br.com.wallet

import br.com.wallet.api.configureBalances
import br.com.wallet.api.configureRouting
import br.com.wallet.api.configureSerialization
import br.com.wallet.database.AccountService
import br.com.wallet.database.getDatabase
import br.com.wallet.entity.Account
import br.com.wallet.gateway.AccountGateway
import br.com.wallet.kafka.startKafka
import io.ktor.server.application.*
import org.apache.log4j.BasicConfigurator
import org.apache.log4j.Level
import org.apache.log4j.LogManager
import org.apache.log4j.Logger
import org.apache.log4j.varia.NullAppender
import kotlin.uuid.ExperimentalUuidApi
import kotlin.uuid.Uuid

// https://ktor.io/docs/server-integrate-database.html
fun main(args: Array<String>) {
    println(args.joinToString())
    io.ktor.server.netty.EngineMain.main(args)
}

fun Application.module() {
    BasicConfigurator.configure()
    LogManager.getRootLogger().level = Level.OFF
    Logger.getRootLogger().addAppender(NullAppender())

    configureRouting()
    configureSerialization()

    val gateway: AccountGateway = try {
        AccountService(getDatabase())
    } catch (t: Throwable) {
        println("Database connection failed")
        t.printStackTrace()
        object : AccountGateway {
            val accounts = listOf(
                Uuid.random().toString(),
                Uuid.random().toString()
            ).let {
                mutableMapOf(
                    it[0] to Account(
                        balance = 1000f,
                        id = it[0]
                    ),
                    it[1] to Account(
                        balance = 0f,
                        id = it[1]
                    )
                )
            }

            override suspend fun delete(
                accountId: String
            ) {
                accounts.remove(accountId)
            }

            override suspend fun findAll() = accounts
                .values
                .toList()

            override suspend fun findById(
                id: String
            ) = accounts
                .values
                .find { it.id == id }

            override suspend fun upsert(
                vararg accounts: Account
            ) {
                accounts.forEach {
                    this.accounts[it.id] = it
                }
            }
        }
    }

    startKafka { message ->
        gateway.upsert(
            message.originAccount,
            message.destinationAccount
        )
        println(message)
    }

    configureBalances(
        gateway
    )

    println("Listening on port ${environment.config.property("ktor.deployment.port").getString()}")
}
