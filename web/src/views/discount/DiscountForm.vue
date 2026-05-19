<script setup lang="ts">
import { ref, onMounted, reactive, watch, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { X } from 'lucide-vue-next'
import type { CreateDiscountRequest, Discount, UpdateDiscountRequest } from '@/types/discount.types'
import type { ValidationError } from '@/types/common.types'
import { CancelModal } from '@/components/common/cancel'
import { AmountInput } from '@/components/common/form/input/amount'
import { Toggle } from '@/components/common/form'
import { Field, FieldContent, FieldLabel, FieldError } from '@/components/ui/field'
import { toast } from 'vue-sonner'
import { useFormErrors } from '@/composables/common/useFormErrors'
import { useDiscountStore } from '@/stores/discount.store'

const props = defineProps<{
  initialData?: Discount | null
}>()

const emit = defineEmits<{
  (e: 'success'): void
}>()

const router = useRouter()
const discountStore = useDiscountStore()
const { setErrors, clearErrors, getErrorMessage, hasError } = useFormErrors()

const formatForInputDate = (value: string | null | undefined): string => {
  if (!value) return ''
  const date = new Date(value)
  if (isNaN(date.getTime())) return ''
  return date.toISOString().split('T')[0] || ''
}

const isEdit = computed(() => !!props.initialData)

const form = reactive({
  name: '',
  type: 'fixed' as 'fixed' | 'percentage',
  value: '0',
  startDate: '',
  endDate: '',
})

const initialForm = ref<string>('')
const showCancelModal = ref(false)
const isDirty = computed(() => JSON.stringify(form) !== initialForm.value)

watch(
  () => props.initialData,
  (newData) => {
    clearErrors()
    if (newData) {
      form.name = newData.name
      form.type = newData.type
      form.value = newData.value
      form.startDate = formatForInputDate(newData.startDate)
      form.endDate = formatForInputDate(newData.endDate)
    } else {
      form.name = ''
      form.type = 'fixed'
      form.value = '0'
      form.startDate = ''
      form.endDate = ''
    }
    initialForm.value = JSON.stringify(form)
  },
  { immediate: true },
)

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

function reset() {
  clearErrors()
  if (props.initialData) {
    form.name = props.initialData.name
    form.type = props.initialData.type
    form.value = props.initialData.value
    form.startDate = formatForInputDate(props.initialData.startDate)
    form.endDate = formatForInputDate(props.initialData.endDate)
  } else {
    form.name = ''
    form.type = 'fixed'
    form.value = '0'
    form.startDate = ''
    form.endDate = ''
  }
  initialForm.value = JSON.stringify(form)
}

defineExpose({ reset, setErrors, clearErrors })

async function handleSubmit() {
  clearErrors()

  const payload: CreateDiscountRequest = {
    name: form.name,
    type: form.type,
    value: form.value.toString(),
    startDate: form.startDate || null,
    endDate: form.endDate || null,
  }

  try {
    if (isEdit.value && props.initialData) {
      await discountStore.update(props.initialData.id, payload as UpdateDiscountRequest)
      toast.success('Discount updated successfully')
    } else {
      await discountStore.create(payload)
      toast.success('Discount created successfully')
    }
    emit('success')
  } catch (err: any) {
    if (err?.error && Array.isArray(err.error)) {
      setErrors(err.error as ValidationError[])
    } else {
      toast.error(err?.message || 'Failed to save discount')
    }
  }
}
</script>

<template>
  <div class="space-y-6 max-w-2xl mx-auto pb-10">
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-semibold">{{ initialData ? 'Edit' : 'Create' }} Discount</h2>
      <div class="flex gap-2">
        <Button variant="outline" @click="handleCancel">
          <X class="size-4 mr-2" />
          Cancel
        </Button>
        <Button :disabled="discountStore.loading" @click="handleSubmit">
          {{ discountStore.loading ? 'Saving...' : 'Save Discount' }}
        </Button>
      </div>
    </div>

    <Card>
      <CardHeader>
        <CardTitle>General Information</CardTitle>
      </CardHeader>
      <CardContent class="space-y-4">
        <Field>
          <FieldLabel>Discount Name</FieldLabel>
          <FieldContent>
            <Input
              id="name"
              v-model="form.name"
              placeholder="e.g. Member Discount, Seasonal Sale"
              required
              :aria-invalid="hasError('Name')"
            />
            <FieldError v-if="hasError('Name')" :errors="[getErrorMessage('Name')]" />
          </FieldContent>
        </Field>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <Field>
            <FieldLabel>Adjustment Type</FieldLabel>
            <FieldContent>
              <Toggle
                v-model="form.type"
                :options="[
                  { label: 'Fixed', value: 'fixed' },
                  { label: 'Percentage', value: 'percentage' },
                ]"
              />
              <FieldError v-if="hasError('Type')" :errors="[getErrorMessage('Type')]" />
            </FieldContent>
          </Field>

          <Field>
            <FieldLabel>Value</FieldLabel>
            <FieldContent>
              <AmountInput
                v-model="form.value"
                :mode="form.type === 'fixed' ? 'money' : 'percentage'"
                :aria-invalid="hasError('Value')"
              />
              <FieldError v-if="hasError('Value')" :errors="[getErrorMessage('Value')]" />
            </FieldContent>
          </Field>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <Field>
            <FieldLabel>Start Date</FieldLabel>
            <FieldContent>
              <Input
                id="startDate"
                v-model="form.startDate"
                type="date"
                :aria-invalid="hasError('StartDate')"
              />
              <FieldError v-if="hasError('StartDate')" :errors="[getErrorMessage('StartDate')]" />
            </FieldContent>
          </Field>
          <Field>
            <FieldLabel>End Date</FieldLabel>
            <FieldContent>
              <Input
                id="endDate"
                v-model="form.endDate"
                type="date"
                :aria-invalid="hasError('EndDate')"
              />
              <FieldError v-if="hasError('EndDate')" :errors="[getErrorMessage('EndDate')]" />
            </FieldContent>
          </Field>
        </div>
      </CardContent>
    </Card>

    <CancelModal v-model:open="showCancelModal" @confirm="confirmCancel" />
  </div>
</template>
