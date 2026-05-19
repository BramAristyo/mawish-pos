<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { toast } from 'vue-sonner'
import AppLayout from '@/layouts/AppLayout.vue'
import ProductForm from './ProductForm.vue'
import { useProductStore } from '@/stores/product.store'
import type { UpdateProductRequest } from '@/types/product.types'
import type { ValidationError } from '@/types/common.types'

const router = useRouter()
const route = useRoute()
const productStore = useProductStore()
const formRef = ref<InstanceType<typeof ProductForm> | null>(null)
const productId = route.params.id as string

onMounted(async () => {
  try {
    await productStore.fetchById(productId)
  } catch (err: any) {
    toast.error(err.message || 'Failed to fetch product details')
    router.push('/catalog/products')
  }
})

const handleSubmit = async (payload: UpdateProductRequest) => {
  try {
    await productStore.update(productId, payload)
    toast.success('Product updated successfully')
    router.push('/catalog/products')
  } catch (err: any) {
    if (err.error && Array.isArray(err.error)) {
      formRef.value?.setErrors(err.error as ValidationError[])
    } else {
      toast.error(err.message || 'An error occurred while updating product')
    }
  }
}
</script>

<template>
  <AppLayout>
    <div class="max-w-4xl mx-auto">
      <div class="mb-6">
        <h1 class="text-2xl font-semibold">Edit Product</h1>
        <p class="text-muted-foreground">Update the details of your product.</p>
      </div>

      <div v-if="productStore.loading && !productStore.selected" class="flex justify-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>

      <ProductForm
        v-else-if="productStore.selected"
        ref="formRef"
        :initial-data="productStore.selected"
        :loading="productStore.loading"
        @submit="handleSubmit"
      />
    </div>
  </AppLayout>
</template>
