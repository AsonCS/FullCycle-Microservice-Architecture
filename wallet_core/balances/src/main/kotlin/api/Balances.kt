@file:OptIn(ExperimentalUuidApi::class)

package br.com.wallet.api

import br.com.wallet.entity.Account
import br.com.wallet.gateway.AccountGateway
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import kotlin.uuid.ExperimentalUuidApi
import kotlin.uuid.Uuid

fun Application.configureBalances(
    gateway: AccountGateway
) {
    routing {
        get("/balances") {
            val accounts = gateway.findAll()
            println("GET /balances: $accounts")

            call.respond(
                status = HttpStatusCode.OK,
                message = accounts
            )
        }

        get("/balances/{account_id}") {
            val id = call
                .parameters["account_id"]
            println("GET /balances/account_id: $id")

            if (id == null) {
                call.respond(
                    status = HttpStatusCode.BadRequest,
                    message = "Invalid Account Id"
                )
                return@get
            }

            val account = gateway
                .findById(id)
            println("GET /balances/account_id: $account")

            if (account == null) {
                call.respond(
                    HttpStatusCode.NotFound
                )
            } else {
                call.respond(
                    status = HttpStatusCode.OK,
                    message = account
                )
            }
        }

        post("/balances") {
            val body = call.receive<Map<String, String>>()
            println("POST /balances: $body")

            val balance = body["balance"]
                ?.toFloatOrNull()
                ?: 0f
            val id = body["id"]
                ?.takeIf(String::isNotBlank)
                ?: Uuid.random().toString()

            gateway.upsert(
                Account(
                    balance = balance,
                    id = id
                )
            )

            call.respond(
                HttpStatusCode.NoContent
            )
        }

        delete("/balances/{account_id}") {
            val accountId = call.parameters["account_id"]
            println("DELETE /balances/account_id: $accountId")

            if (accountId == null) {
                call.respond(
                    status = HttpStatusCode.BadRequest,
                    message = "Invalid Account Id"
                )
                return@delete
            } else {
                gateway.delete(accountId)
            }

            call.respond(
                HttpStatusCode.NoContent
            )
        }
    }
}
