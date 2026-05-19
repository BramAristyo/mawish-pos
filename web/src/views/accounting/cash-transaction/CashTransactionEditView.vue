<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useCashTransactionStore } from '@/stores/cashTransaction.store'
import AppLayout from '@/layouts/AppLayout.vue'
import CashTransactionForm from './CashTransactionForm.vue'
import { Skeleton } from '@/components/ui/skeleton'
import { Button } from '@/components/ui/button'

const router = useRouter()
const route = useRoute()
const cashTransactionStore = useCashTransactionStore()

const id = route.params.id as string

onMounted(async () => {
  if (id) {
    await cashTransactionStore.fetchById(id)
  }
})

function handleSuccess() {
  router.push({ name: 'cash-transactions' })
}
</script>

<template>
  <AppLayout>
    <div v-if="cashTransactionStore.loading && !cashTransactionStore.selected" class="max-w-2xl mx-auto py-6 space-y-4">
      <Skeleton class="h-10 w-48" />
      <Skeleton class="h-[400px] w-full" />
    </div>
    <CashTransactionForm
      v-else-if="cashTransactionStore.selected"
      :initial-data="cashTransactionStore.selected"
      @success="handleSuccess"
    />
    <div v-else class="text-center py-20">
      <p class="text-muted-foreground">Transaction not found.</p>
      <Button variant="outline" class="mt-4" @click="router.push({ name: 'cash-transactions' })">
        Back to list
      </Button>
    </div>
  </AppLayout>
</template>
