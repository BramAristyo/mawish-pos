<script setup lang="ts">
import { computed, reactive, watch, ref } from 'vue'
import { useTaxStore } from '@/stores/tax.stores'
import type { Tax, CreateTaxRequest, UpdateTaxRequest } from '@/types/tax.types'
import type { ValidationError } from '@/types/common.types'

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Field, FieldContent, FieldLabel, FieldError } from '@/components/ui/field'
import { CancelModal } from '@/components/common/cancel'
import { toast } from 'vue-sonner'
import { useFormErrors } from '@/composables/common/useFormErrors'

const props = defineProps<{
  open: boolean
  tax?: Tax | null
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'success'): void
}>()

const taxStore = useTaxStore()

const { setErrors, clearErrors, getErrorMessage, hasError } = useFormErrors()

const isOpen = computed({
  get: () => props.open,
  set: (value) => emit('update:open', value),
})

const isEdit = computed(() => !!props.tax)

const form = reactive<CreateTaxRequest>({
  name: '',
  percentage: '',
})

const initialForm = ref<string>('')
const showCancelModal = ref(false)

const isDirty = computed(() => {
  return JSON.stringify(form) !== initialForm.value
})

function handleCancel() {
  if (isDirty.value) {
    showCancelModal.value = true
  } else {
    isOpen.value = false
  }
}

function confirmCancel() {
  isOpen.value = false
}

function resetForm() {
  clearErrors()
  if (props.tax) {
    form.name = props.tax.name
    form.percentage = props.tax.percentage
  } else {
    form.name = ''
    form.percentage = ''
  }
  initialForm.value = JSON.stringify(form)
}

watch(
  () => props.open,
  (open) => {
    if (open) {
      resetForm()
    }
  },
)

watch(
  () => props.tax,
  () => {
    resetForm()
  },
  { immediate: true },
)

async function handleSubmit() {
  clearErrors()
  try {
    if (isEdit.value && props.tax) {
      await taxStore.update(props.tax.id, form as UpdateTaxRequest)
      toast.success('Tax updated successfully')
    } else {
      await taxStore.create(form)
      toast.success('Tax created successfully')
    }
    emit('success')
    isOpen.value = false
  } catch (err: any) {
    if (err?.error && Array.isArray(err.error)) {
      setErrors(err.error as ValidationError[])
    } else {
      toast.error(err?.message || 'Failed to save tax')
    }
  }
}
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogContent 
      class="sm:max-w-106.25"
      @pointer-down-outside="(e) => isDirty && e.preventDefault()"
      @escape-key-down="(e) => isDirty && e.preventDefault()"
    >
      <form @submit.prevent="handleSubmit">
        <DialogHeader>
          <DialogTitle>{{ isEdit ? 'Edit Tax' : 'Add Tax' }}</DialogTitle>
          <DialogDescription>
            {{
              isEdit
                ? 'Update the details of your tax here.'
                : 'Create a new tax for your transactions.'
            }}
          </DialogDescription>
        </DialogHeader>

        <div class="grid gap-4 py-4">
          <Field>
            <FieldLabel>Name</FieldLabel>
            <FieldContent>
              <Input
                v-model="form.name"
                placeholder="Tax name (e.g. VAT)"
                required
                :aria-invalid="hasError('Name')"
              />
              <FieldError v-if="hasError('Name')" :errors="[getErrorMessage('Name')]" />
            </FieldContent>
          </Field>

          <Field>
            <FieldLabel>Percentage (%)</FieldLabel>
            <FieldContent>
              <Input
                v-model="form.percentage"
                type="number"
                step="0.01"
                placeholder="Tax percentage"
                required
                :aria-invalid="hasError('Percentage')"
              />
              <FieldError v-if="hasError('Percentage')" :errors="[getErrorMessage('Percentage')]" />
            </FieldContent>
          </Field>
        </div>

        <DialogFooter>
          <Button type="button" variant="outline" @click="handleCancel"> Cancel </Button>
          <Button type="submit" :disabled="taxStore.loading">
            {{ taxStore.loading ? 'Saving...' : 'Save' }}
          </Button>
        </DialogFooter>
      </form>
      <CancelModal v-model:open="showCancelModal" @confirm="confirmCancel" />
    </DialogContent>
  </Dialog>
</template>
