import { productApi } from '@/api/product.api'
import type { BaseFilterRequest, Meta } from '@/types/common.types'
import type { Product, StoreProductRequest, UpdateProductRequest } from '@/types/product.types'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useProductStore = defineStore('product', () => {
  const products = ref<Product[]>([])
  const meta = ref<Meta | null>(null)
  const selected = ref<Product | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchPaginated(params: BaseFilterRequest) {
    loading.value = true
    error.value = null
    try {
      const res = await productApi.paginate(params)
      products.value = res.data
      meta.value = res.meta || null
    } catch (err: unknown) {
      const errorResponse = err as { message?: string }
      error.value = errorResponse.message || 'Failed to fetch products'
    } finally {
      loading.value = false
    }
  }

  async function fetchAll() {
    loading.value = true
    error.value = null
    try {
      const res = await productApi.getAll()
      products.value = res.data
    } catch (err: unknown) {
      const errorResponse = err as { message?: string }
      error.value = errorResponse.message || 'Failed to fetch all products'
    } finally {
      loading.value = false
    }
  }

  async function fetchById(id: string) {
    loading.value = true
    error.value = null
    try {
      const res = await productApi.findById(id)
      selected.value = res.data
    } catch (err: unknown) {
      const errorResponse = err as { message?: string }
      error.value = errorResponse.message || 'Failed to fetch product'
    } finally {
      loading.value = false
    }
  }

  async function create(payload: StoreProductRequest) {
    loading.value = true
    error.value = null
    try {
      const res = await productApi.store(payload)
      products.value = [res.data, ...products.value]
      return res.data
    } catch (err: unknown) {
      const errorResponse = err as { message?: string }
      error.value = errorResponse.message || 'Failed to create product'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function update(id: string, payload: UpdateProductRequest) {
    loading.value = true
    error.value = null
    try {
      const res = await productApi.update(id, payload)
      const index = products.value.findIndex((p) => p.id === id)
      if (index !== -1) products.value[index] = res.data
      if (selected.value?.id === id) selected.value = res.data
      return res.data
    } catch (err: unknown) {
      const errorResponse = err as { message?: string }
      error.value = errorResponse.message || 'Failed to update product'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteProduct(id: string) {
    loading.value = true
    error.value = null
    try {
      await productApi.delete(id)
      products.value = products.value.filter((p) => p.id !== id)
    } catch (err: unknown) {
      const errorResponse = err as { message?: string }
      error.value = errorResponse.message || 'Failed to delete product'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function restore(id: string) {
    loading.value = true
    error.value = null
    try {
      const res = await productApi.restore(id)
      const index = products.value.findIndex((p) => p.id === id)
      if (index !== -1) products.value[index] = res.data
      return res.data
    } catch (err: unknown) {
      const errorResponse = err as { message?: string }
      error.value = errorResponse.message || 'Failed to restore product'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function ensureDataLoaded() {
    if (products.value.length === 0) {
      await fetchAll()
    }
  }

  return {
    products,
    meta,
    selected,
    loading,
    error,
    fetchPaginated,
    fetchAll,
    fetchById,
    create,
    update,
    delete: deleteProduct,
    restore,
    ensureDataLoaded,
  }
})
