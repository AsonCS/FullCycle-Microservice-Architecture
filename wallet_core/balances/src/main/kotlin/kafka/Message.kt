package br.com.wallet.kafka

import br.com.wallet.entity.Account

data class Message(
    val originAccount: Account,
    val destinationAccount: Account
)
