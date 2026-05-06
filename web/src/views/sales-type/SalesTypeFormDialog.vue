<script setup lang="ts">
import { computed, reactive, watch, ref } from 'vue'
import { useSalesType } from '@/composables/useSalesType'
import type {
  SalesType,
  CreateSalesTypeRequest,
  UpdateSalesTypeRequest,
} from '@/types/salesType.types'
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
import { AmountInput } from '@/components/common/form/input/amount'
import { Toggle } from '@/components/common/form'
import { Field, FieldContent, FieldLabel, FieldError } from '@/components/ui/field'
import { CancelModal } from '@/components/common/cancel'
import { Plus, Trash2 } from 'lucide-vue-next'
import { toast } from 'vue-sonner'
import { useFormErrors } from '@/composables/common/useFormErrors'

const props = defineProps<{
  open: boolean
  salesType?: SalesType | null
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'success'): void
}>()

const { create, update, fetchById, salesType: salesTypeDetail, loading } = useSalesType()
const { setErrors, clearErrors, getErrorMessage, hasError } = useFormErrors()

const isOpen = computed({
  get: () => props.open,
  set: (value) => emit('update:open', value),
})

const isEdit = computed(() => !!props.salesType)

const form = reactive<CreateSalesTypeRequest>({
  name: '',
  charges: [],
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

async function resetForm() {
  clearErrors()
  if (props.salesType) {
    await fetchById(props.salesType.id)
    if (salesTypeDetail.value) {
      form.name = salesTypeDetail.value.name
      form.charges = salesTypeDetail.value.charges.map((c) => ({
        id: c.id,
        name: c.name,
        type: c.type,
        amount: c.amount,
      }))
    }
  } else {
    form.name = ''
    form.charges = []
  }
  initialForm.value = JSON.stringify(form)
}

watch(
  () => props.open,
  async (open) => {
    if (open) {
      await resetForm()
    }
  },
)

function addCharge() {
  form.charges.push({
    name: '',
    type: 'fixed',
    amount: '0',
  })
}

function removeCharge(index: number) {
  form.charges.splice(index, 1)
}

async function handleSubmit() {
  clearErrors()
  try {
    if (isEdit.value && props.salesType) {
      await update(props.salesType.id, form as UpdateSalesTypeRequest)
      toast.success('Sales type updated successfully')
    } else {
      await create(form)
      toast.success('Sales type created successfully')
    }
    emit('success')
    isOpen.value = false
  } catch (err: any) {
    if (err?.error && Array.isArray(err.error)) {
      setErrors(err.error as ValidationError[])
    } else {
      toast.error(err?.message || 'Failed to save sales type')
    }
  }
}
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogContent 
      class="sm:max-w-150 max-h-[90vh] overflow-y-auto"
      @pointer-down-outside="(e) => isDirty && e.preventDefault()"
      @escape-key-down="(e) => isDirty && e.preventDefault()"
    >
      <form @submit.prevent="handleSubmit">
        <DialogHeader>
          <DialogTitle>{{ isEdit ? 'Edit Sales Type' : 'Add Sales Type' }}</DialogTitle>
          <DialogDescription>
            {{
              isEdit
                ? 'Update the details of your sales type here.'
                : 'Create a new sales type for your products.'
            }}
          </DialogDescription>
        </DialogHeader>

        <div class="grid gap-6 py-4">
          <Field>
            <FieldLabel>Name</FieldLabel>
            <FieldContent>
              <Input
                v-model="form.name"
                placeholder="Sales type name (e.g. Dine-In, Delivery)"
                required
                :aria-invalid="hasError('Name')"
              />
              <FieldError v-if="hasError('Name')" :errors="[getErrorMessage('Name')]" />
            </FieldContent>
          </Field>

          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <h3 class="text-sm font-medium">Additional Charges</h3>
              <Button type="button" variant="outline" size="sm" @click="addCharge">
                <Plus class="size-4 mr-2" />
                Add Charge
              </Button>
            </div>

            <div v-if="form.charges.length === 0" class="text-sm text-muted-foreground py-4 text-center border-2 border-dashed rounded-lg">
              No additional charges added.
            </div>

            <div v-else class="space-y-4">
              <div v-for="(charge, index) in form.charges" :key="index" class="p-4 border rounded-lg space-y-4 bg-muted/20 relative">
                <Button 
                  type="button" 
                  variant="ghost" 
                  size="icon" 
                  class="absolute top-2 right-2 text-destructive"
                  @click="removeCharge(index)"
                >
                  <Trash2 class="size-4" />
                </Button>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-4 pt-4">
                  <Field>
                    <FieldLabel>Charge Name</FieldLabel>
                    <FieldContent>
                      <Input
                        v-model="charge.name"
                        placeholder="e.g. Service Charge"
                        required
                        :aria-invalid="hasError(`Charges[${index}].Name`)"
                      />
                      <FieldError v-if="hasError(`Charges[${index}].Name`)" :errors="[getErrorMessage(`Charges[${index}].Name`)]" />
                    </FieldContent>
                  </Field>

                  <Field>
                    <FieldContent>
                      <Toggle
                        v-model="charge.type"
                        label="Type"
                        :options="[
                          { label: 'Fixed', value: 'fixed' },
                          { label: 'Percentage', value: 'percentage' },
                        ]"
                      />
                      <FieldError
                        v-if="hasError(`Charges[${index}].Type`)"
                        :errors="[getErrorMessage(`Charges[${index}].Type`)]"
                      />
                    </FieldContent>
                  </Field>

                  <Field>
                    <FieldLabel>Amount</FieldLabel>
                    <FieldContent>
                      <AmountInput
                        v-model="charge.amount"
                        :mode="charge.type === 'fixed' ? 'money' : 'percentage'"
                        placeholder="0"
                        :aria-invalid="hasError(`Charges[${index}].Amount`)"
                      />
                      <FieldError v-if="hasError(`Charges[${index}].Amount`)" :errors="[getErrorMessage(`Charges[${index}].Amount`)]" />
                    </FieldContent>
                  </Field>
                </div>
              </div>
            </div>
          </div>
        </div>

        <DialogFooter>
          <Button type="button" variant="outline" @click="handleCancel"> Cancel </Button>
          <Button type="submit" :disabled="loading">
            {{ loading ? 'Saving...' : 'Save' }}
          </Button>
        </DialogFooter>
      </form>
      <CancelModal v-model:open="showCancelModal" @confirm="confirmCancel" />
    </DialogContent>
  </Dialog>
</template>
