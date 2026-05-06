export type CashTransactionType = 'in' | 'out'

export interface CashTransaction {
  id: string
  coaName: string
  type: CashTransactionType
  amount: number
  description: string
  date: string
  createdAt: string
  updatedAt: string
}

export interface CreateCashTransactionRequest {
  coaId: string
  type: CashTransactionType
  amount: number
  description: string
  date: string
}

export interface UpdateCashTransactionRequest {
  coaId: string
  type: CashTransactionType
  amount: number
  description: string
  date: string
}
