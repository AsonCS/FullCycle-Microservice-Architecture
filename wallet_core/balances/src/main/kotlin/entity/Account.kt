package br.com.wallet.entity

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class Account(
    @SerialName("balance")
    val balance: Float,
    @SerialName("id")
    val id: String,
    @SerialName("updatedAt")
    val updatedAt: String? = null
)
