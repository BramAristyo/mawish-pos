<script setup lang="ts">
import { reactive, onMounted, computed, watch, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useCashTransactionStore } from '@/stores/cashTransaction.store'
import { useCoaStore } from '@/stores/coa.store'
import type {
  CashTransaction,
  CreateCashTransactionRequest,
  UpdateCashTransactionRequest,
} from '@/types/cashTransaction.types'
import type { ValidationError } from '@/types/common.types'
import { useFormErrors } from '@/composables/common/useFormErrors'
import { toast } from 'vue-sonner'

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Field, FieldContent, FieldLabel, FieldError } from '@/components/ui/field'
import { AmountInput } from '@/components/common/form/input/amount'
import { Toggle } from '@/components/common/form'
import { CancelModal } from '@/components/common/cancel'

const props = defineProps<{
  initialData?: CashTransaction | null
}>()

const emit = defineEmits<{
  (e: 'success'): void
}>()

const router = useRouter()
const cashTransactionStore = useCashTransactionStore()
const coaStore = useCoaStore()
const { setErrors, clearErrors, getErrorMessage, hasError } = useFormErrors()

const isEdit = computed(() => !!props.initialData)

const form = reactive<CreateCashTransactionRequest>({
  coaId: '',
  type: 'in',
  amount: 0,
  description: '',
  date: new Date().toISOString().split('T')[0] || '',
})

const initialForm = ref<string>('')
const showCancelModal = ref(false)
const isDirty = computed(() => JSON.stringify(form) !== initialForm.value)

const typeOptions = [
  { label: 'Cash In', value: 'in' },
  { label: 'Cash Out', value: 'out' },
]

watch(
  () => form.coaId,
  (newCoaId) => {
    const selectedCoa = coaStore.coas.find((coa) => coa.id === newCoaId)
    if (selectedCoa) {
      form.type = selectedCoa.type as 'in' | 'out'
    }
  },
)

onMounted(async () => {
  await coaStore.ensureDataLoaded()
  if (props.initialData) {
    form.coaId = props.initialData.coaId
    form.type = props.initialData.type
    form.amount = props.initialData.amount
    form.description = props.initialData.description
    form.date = props.initialData.date.substring(0, 10) || ''
  }
  initialForm.value = JSON.stringify(form)
})

function handleCancel() {
  if (isDirty.value) {
    showCancelModal.value = true
  } else {
    router.back()
  }
}

function confirmCancel() {
  router.back()
}

async function handleSubmit() {
  clearErrors()
  try {
    const payload = {
      ...form,
      amount: Number(form.amount),
    }

    if (isEdit.value && props.initialData) {
      await cashTransactionStore.update(
        props.initialData.id,
        payload as UpdateCashTransactionRequest,
      )
      toast.success('Cash transaction updated successfully')
    } else {
      await cashTransactionStore.create(payload)
      toast.success('Cash transaction created successfully')
    }
    emit('success')
  } catch (err: any) {
    if (err?.error && Array.isArray(err.error)) {
      setErrors(err.error as ValidationError[])
    } else {
      toast.error(err?.message || 'Failed to save cash transaction')
    }
  }
}
</script>

<template>
  <div class="max-w-2xl mx-auto py-6">
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">{{ isEdit ? 'Edit' : 'New' }} Cash Transaction</h1>
      <Button variant="outline" @click="router.back()"> Back </Button>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>Transaction Details</CardTitle>
      </CardHeader>
      <CardContent>
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <Field>
            <FieldLabel>Date</FieldLabel>
            <FieldContent>
              <Input v-model="form.date" type="date" required :aria-invalid="hasError('Date')" />
              <FieldError v-if="hasError('Date')" :errors="[getErrorMessage('Date')]" />
            </FieldContent>
          </Field>

          <Field>
            <FieldLabel>Account (COA)</FieldLabel>
            <FieldContent>
              <select
                v-model="form.coaId"
                required
                :aria-invalid="hasError('CoaId')"
                class="flex h-9 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm transition-colors focus:outline-none focus:ring-1 focus:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
              >
                <option value="" disabled>Select account</option>
                <option v-for="coa in coaStore.coas" :key="coa.id" :value="coa.id">
                  {{ coa.name }}
                </option>
              </select>
              <FieldError v-if="hasError('CoaId')" :errors="[getErrorMessage('CoaId')]" />
            </FieldContent>
          </Field>

          <Field>
            <FieldLabel>Type</FieldLabel>
            <FieldContent>
              <Toggle v-model="form.type" :options="typeOptions" disabled />
              <FieldError v-if="hasError('Type')" :errors="[getErrorMessage('Type')]" />
            </FieldContent>
          </Field>

          <Field>
            <FieldLabel>Amount</FieldLabel>
            <FieldContent>
              <AmountInput v-model="form.amount" required :aria-invalid="hasError('Amount')" />
              <FieldError v-if="hasError('Amount')" :errors="[getErrorMessage('Amount')]" />
            </FieldContent>
          </Field>

          <Field>
            <FieldLabel>Description</FieldLabel>
            <FieldContent>
              <Input
                v-model="form.description"
                placeholder="Transaction description"
                required
                :aria-invalid="hasError('Description')"
              />
              <FieldError
                v-if="hasError('Description')"
                :errors="[getErrorMessage('Description')]"
              />
            </FieldContent>
          </Field>

          <div class="flex items-center justify-end gap-2 pt-4">
            <Button type="button" variant="outline" @click="handleCancel"> Cancel </Button>
            <Button type="submit" :disabled="cashTransactionStore.loading">
              {{ cashTransactionStore.loading ? 'Saving...' : 'Save' }}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
    <CancelModal v-model:open="showCancelModal" @confirm="confirmCancel" />
  </div>
</template>
