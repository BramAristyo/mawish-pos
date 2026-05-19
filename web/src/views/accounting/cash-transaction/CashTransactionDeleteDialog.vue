<script setup lang="ts">
import { computed } from 'vue'
import { useCashTransactionStore } from '@/stores/cashTransaction.store'
import type { CashTransaction } from '@/types/cashTransaction.types'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { toast } from 'vue-sonner'

const props = defineProps<{
  open: boolean
  transaction: CashTransaction | null
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'success'): void
}>()

const cashTransactionStore = useCashTransactionStore()

const isOpen = computed({
  get: () => props.open,
  set: (value) => emit('update:open', value),
})

async function handleDelete() {
  if (!props.transaction) return
  try {
    await cashTransactionStore.remove(props.transaction.id)
    toast.success('Cash transaction deleted successfully')
    emit('success')
    isOpen.value = false
  } catch (error: any) {
    toast.error(error?.message || 'Failed to delete cash transaction')
  }
}
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogContent class="sm:max-w-[425px]">
      <DialogHeader>
        <DialogTitle>Delete Cash Transaction</DialogTitle>
        <DialogDescription>
          Are you sure you want to delete this transaction: <strong>{{ transaction?.description }}</strong>? This action cannot be undone.
        </DialogDescription>
      </DialogHeader>

      <DialogFooter class="flex gap-2 justify-end">
        <Button type="button" variant="outline" @click="isOpen = false"> Cancel </Button>
        <Button variant="destructive" :disabled="cashTransactionStore.loading" @click="handleDelete">
          {{ cashTransactionStore.loading ? 'Deleting...' : 'Delete' }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
