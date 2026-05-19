<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import AppLayout from '@/layouts/AppLayout.vue'
import ProductForm from './ProductForm.vue'
import { useProductStore } from '@/stores/product.store'
import type { StoreProductRequest } from '@/types/product.types'
import type { ValidationError } from '@/types/common.types'

const router = useRouter()
const productStore = useProductStore()
const formRef = ref<InstanceType<typeof ProductForm> | null>(null)

const handleSubmit = async (payload: StoreProductRequest) => {
  try {
    await productStore.create(payload)
    toast.success('Product created successfully')
    router.push('/catalog/products')
  } catch (err: any) {
    if (err.error && Array.isArray(err.error)) {
      formRef.value?.setErrors(err.error as ValidationError[])
    } else {
      toast.error(err.message || 'An error occurred while creating product')
    }
  }
}
</script>

<template>
  <AppLayout>
    <div class="max-w-4xl mx-auto">
      <div class="mb-6">
        <h1 class="text-2xl font-semibold">Create Product</h1>
        <p class="text-muted-foreground">Add a new product to your catalog.</p>
      </div>

      <ProductForm
        ref="formRef"
        :loading="productStore.loading"
        @submit="handleSubmit"
      />
    </div>
  </AppLayout>
</template>
