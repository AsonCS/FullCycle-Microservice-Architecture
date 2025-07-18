package br.com.wallet.database

import br.com.wallet.entity.Account
import br.com.wallet.gateway.AccountGateway
import kotlinx.coroutines.Dispatchers.IO
import org.jetbrains.exposed.sql.*
import org.jetbrains.exposed.sql.SqlExpressionBuilder.eq
import org.jetbrains.exposed.sql.transactions.experimental.newSuspendedTransaction
import org.jetbrains.exposed.sql.transactions.transaction

class AccountService(
    database: Database
) : AccountGateway {
    object Accounts : Table() {
        val id = varchar(
            name = "id",
            length = 255
        ).uniqueIndex()
        val balance = float(
            name = "balance"
        )
    }

    init {
        transaction(database) {
            SchemaUtils.create(Accounts)
        }
    }

    override suspend fun findAll(): List<Account> = dbQuery {
        Accounts.selectAll()
            .map {
                Account(
                    id = it[Accounts.id],
                    balance = it[Accounts.balance]
                )
            }
    }

    override suspend fun findById(
        id: String
    ): Account? = dbQuery {
        Accounts
            .selectAll()
            /*.select(
                Accounts.balance
            )*/
            .where(
                Accounts.id eq id
            )
            .map {
                Account(
                    id = it[Accounts.id],
                    balance = it[Accounts.balance]
                )
            }.singleOrNull()
    }

    override suspend fun upsert(
        vararg accounts: Account
    ) = dbQuery {
        accounts.forEach { account ->
            Accounts.upsert(
                Accounts.balance,
                where = { Accounts.id eq account.id },
            ) {
                it[id] = account.id
                it[balance] = account.balance
            }
        }
    }

    private suspend fun <T> dbQuery(block: suspend () -> T): T =
        newSuspendedTransaction(IO) { block() }
}
