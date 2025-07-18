package br.com.wallet.gateway

import br.com.wallet.entity.Account

interface AccountGateway {
    suspend fun findAll(): List<Account>

    suspend fun findById(
        id: String
    ): Account?

    suspend fun upsert(
        vararg accounts: Account
    )
}