<script setup lang="ts">
import { computed, reactive, watch, ref } from 'vue'
import { useCategoryStore } from '@/stores/category.stores'
import type { Category, CreateCategoryRequest, UpdateCategoryRequest } from '@/types/category.types'
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
  category?: Category | null
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'success'): void
}>()

const categoryStore = useCategoryStore()

const { setErrors, clearErrors, getErrorMessage, hasError } = useFormErrors()

const isOpen = computed({
  get: () => props.open,
  set: (value) => emit('update:open', value),
})

const isEdit = computed(() => !!props.category)

const form = reactive<CreateCategoryRequest>({
  name: '',
  description: '',
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
  if (props.category) {
    form.name = props.category.name
    form.description = props.category.description
  } else {
    form.name = ''
    form.description = ''
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
  () => props.category,
  () => {
    resetForm()
  },
  { immediate: true },
)

async function handleSubmit() {
  clearErrors()
  try {
    if (isEdit.value && props.category) {
      await categoryStore.update(props.category.id, form as UpdateCategoryRequest)
      toast.success('Category updated successfully')
    } else {
      await categoryStore.create(form)
      toast.success('Category created successfully')
    }
    emit('success')
    isOpen.value = false
  } catch (err: any) {
    if (err?.error && Array.isArray(err.error)) {
      setErrors(err.error as ValidationError[])
    } else {
      toast.error(err?.message || 'Failed to save category')
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
          <DialogTitle>{{ isEdit ? 'Edit Category' : 'Add Category' }}</DialogTitle>
          <DialogDescription>
            {{
              isEdit
                ? 'Update the details of your category here.'
                : 'Create a new category for your products.'
            }}
          </DialogDescription>
        </DialogHeader>

        <div class="grid gap-4 py-4">
          <Field>
            <FieldLabel>Name</FieldLabel>
            <FieldContent>
              <Input
                v-model="form.name"
                placeholder="Category name"
                required
                :aria-invalid="hasError('Name')"
              />
              <FieldError v-if="hasError('Name')" :errors="[getErrorMessage('Name')]" />
            </FieldContent>
          </Field>

          <Field>
            <FieldLabel>Description</FieldLabel>
            <FieldContent>
              <Input
                v-model="form.description"
                placeholder="Category description"
                :aria-invalid="hasError('Description')"
              />
              <FieldError
                v-if="hasError('Description')"
                :errors="[getErrorMessage('Description')]"
              />
            </FieldContent>
          </Field>
        </div>

        <DialogFooter>
          <Button type="button" variant="outline" @click="handleCancel"> Cancel </Button>
          <Button type="submit" :disabled="categoryStore.loading">
            {{ categoryStore.loading ? 'Saving...' : 'Save' }}
          </Button>
        </DialogFooter>
      </form>
      <CancelModal v-model:open="showCancelModal" @confirm="confirmCancel" />
    </DialogContent>
  </Dialog>
</template>
