import { cashTransactionApi } from '@/api/cashTransaction.api'
import type {
  CashTransaction,
  CreateCashTransactionRequest,
  UpdateCashTransactionRequest,
} from '@/types/cashTransaction.types'
import type { BaseFilterRequest, Meta } from '@/types/common.types'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useCashTransactionStore = defineStore('cashTransaction', () => {
  const transactions = ref<CashTransaction[]>([])
  const meta = ref<Meta | null>(null)
  const selected = ref<CashTransaction | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchPaginated(params: BaseFilterRequest) {
    loading.value = true
    error.value = null
    try {
      const res = await cashTransactionApi.paginate(params)
      transactions.value = res.data
      meta.value = res.meta || null
    } catch (err: any) {
      error.value = err?.message || 'Unsuccessfully fetch cash transactions'
    } finally {
      loading.value = false
    }
  }

  async function fetchById(id: string) {
    loading.value = true
    error.value = null
    try {
      const res = await cashTransactionApi.getById(id)
      selected.value = res.data
    } catch (err: any) {
      error.value = err?.message || 'Unsuccessfully fetch cash transaction'
    } finally {
      loading.value = false
    }
  }

  async function create(payload: CreateCashTransactionRequest) {
    loading.value = true
    error.value = null
    try {
      const res = await cashTransactionApi.create(payload)
      transactions.value = [res.data, ...transactions.value]
    } catch (err: any) {
      error.value = err?.message || 'Unsuccessfully create cash transaction'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function update(id: string, payload: UpdateCashTransactionRequest) {
    loading.value = true
    error.value = null
    try {
      const res = await cashTransactionApi.update(id, payload)
      const index = transactions.value.findIndex((t) => t.id === id)
      if (index !== -1) transactions.value[index] = res.data
    } catch (err: any) {
      error.value = err?.message || 'Unsuccessfully update cash transaction'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function remove(id: string) {
    loading.value = true
    error.value = null
    try {
      await cashTransactionApi.delete(id)
      transactions.value = transactions.value.filter((t) => t.id !== id)
    } catch (err: any) {
      error.value = err?.message || 'Unsuccessfully delete cash transaction'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    transactions,
    meta,
    selected,
    loading,
    error,
    fetchPaginated,
    fetchById,
    create,
    update,
    remove,
  }
})
