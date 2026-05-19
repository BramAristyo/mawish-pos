<script setup lang="ts">
import { computed, reactive, watch, ref } from 'vue'
import { useCoaStore } from '@/stores/coa.store'
import type { Coa, CreateCoaRequest, UpdateCoaRequest } from '@/types/coa.types'
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
import { Toggle } from '@/components/common/form'
import { CancelModal } from '@/components/common/cancel'
import { toast } from 'vue-sonner'
import { useFormErrors } from '@/composables/common/useFormErrors'

const props = defineProps<{
  open: boolean
  coa?: Coa | null
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'success'): void
}>()

const coaStore = useCoaStore()

const { setErrors, clearErrors, getErrorMessage, hasError } = useFormErrors()

const isOpen = computed({
  get: () => props.open,
  set: (value) => emit('update:open', value),
})

const isEdit = computed(() => !!props.coa)

const form = reactive<CreateCoaRequest>({
  name: '',
  type: 'in',
  isOperational: false,
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
  if (props.coa) {
    form.name = props.coa.name
    form.type = props.coa.type
    form.isOperational = props.coa.IsOperational
  } else {
    form.name = ''
    form.type = 'in'
    form.isOperational = false
  }
  initialForm.value = JSON.stringify(form)
}

const typeOptions = [
  { label: 'Income', value: 'in' },
  { label: 'Expenses', value: 'out' },
]

const operationalOptions = [
  { label: 'Operational', value: true },
  { label: 'General', value: false },
]

watch(
  () => props.open,
  (open) => {
    if (open) {
      resetForm()
    }
  },
)

watch(
  () => props.coa,
  () => {
    resetForm()
  },
  { immediate: true },
)

async function handleSubmit() {
  clearErrors()
  try {
    if (isEdit.value && props.coa) {
      await coaStore.update(props.coa.id, form as UpdateCoaRequest)
      toast.success('Account updated successfully')
    } else {
      await coaStore.create(form)
      toast.success('Account created successfully')
    }
    emit('success')
    isOpen.value = false
  } catch (err: any) {
    if (err?.error && Array.isArray(err.error)) {
      setErrors(err.error as ValidationError[])
    } else {
      toast.error(err?.message || 'Failed to save account')
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
          <DialogTitle>{{ isEdit ? 'Edit Account' : 'Add Account' }}</DialogTitle>
          <DialogDescription>
            {{
              isEdit
                ? 'Update the details of your account here.'
                : 'Create a new account for your finances.'
            }}
          </DialogDescription>
        </DialogHeader>

        <div class="grid gap-4 py-4">
          <Field>
            <FieldLabel>Name</FieldLabel>
            <FieldContent>
              <Input
                v-model="form.name"
                placeholder="Account name"
                required
                :aria-invalid="hasError('Name')"
              />
              <FieldError v-if="hasError('Name')" :errors="[getErrorMessage('Name')]" />
            </FieldContent>
          </Field>

          <Field>
            <FieldLabel>Type</FieldLabel>
            <FieldContent>
              <Toggle
                v-model="form.type"
                :options="typeOptions"
              />
              <FieldError v-if="hasError('Type')" :errors="[getErrorMessage('Type')]" />
            </FieldContent>
          </Field>

          <Field>
            <FieldLabel>Account Classification</FieldLabel>
            <FieldContent>
              <Toggle
                v-model="form.isOperational"
                :options="operationalOptions"
              />
              <FieldError v-if="hasError('IsOperational')" :errors="[getErrorMessage('IsOperational')]" />
            </FieldContent>
          </Field>
        </div>

        <DialogFooter>
          <Button type="button" variant="outline" @click="handleCancel"> Cancel </Button>
          <Button type="submit" :disabled="coaStore.loading">
            {{ coaStore.loading ? 'Saving...' : 'Save' }}
          </Button>
        </DialogFooter>
      </form>
      <CancelModal v-model:open="showCancelModal" @confirm="confirmCancel" />
    </DialogContent>
  </Dialog>
</template>
