package org.example.extensions

import com.google.gson.Gson
import org.example.message.BalanceUpdated
import org.example.message.Message
import org.example.message.TransactionCreated

fun String.toMessage(): Message = Gson().fromJson(
    this,
    when {
        this.contains("BalanceUpdated") ->
            BalanceUpdated::class.java

        this.contains("TransactionCreated") ->
            TransactionCreated::class.java

        else -> throw UnsupportedOperationException()
    }
)
