package br.com.wallet.api

import br.com.wallet.gateway.AccountGateway
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

fun Application.configureBalances(
    gateway: AccountGateway
) {
    routing {
        get("/balances") {
            call.respond(
                status = HttpStatusCode.OK,
                message = gateway.findAll()
            )
        }

        get("/balances/account_id") {
            val id = call
                .parameters["account_id"]
            if (id == null) {
                call.respond(
                    status = HttpStatusCode.BadRequest,
                    message = "Invalid Account Id"
                )
                return@get
            }
            val account = gateway
                .findById(id)

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
    }
}
