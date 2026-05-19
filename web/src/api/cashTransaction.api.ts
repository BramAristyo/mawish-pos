import http from '@/lib/http'
import type { BaseFilterRequest, BaseResponse } from '@/types/common.types'
import type {
  CashTransaction,
  CreateCashTransactionRequest,
  UpdateCashTransactionRequest,
} from '@/types/cashTransaction.types'

export const cashTransactionApi = {
  paginate: (params: BaseFilterRequest) =>
    http.get<any, BaseResponse<CashTransaction[]>>('/cash-transactions', { params }),

  getById: (id: string) =>
    http.get<any, BaseResponse<CashTransaction>>(`/cash-transactions/${id}`),

  create: (payload: CreateCashTransactionRequest) =>
    http.post<any, BaseResponse<CashTransaction>>('/cash-transactions', payload),

  update: (id: string, payload: UpdateCashTransactionRequest) =>
    http.put<any, BaseResponse<CashTransaction>>(`/cash-transactions/${id}`, payload),

  delete: (id: string) =>
    http.delete<any, BaseResponse<null>>(`/cash-transactions/${id}`),
}
