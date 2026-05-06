<script setup lang="ts">
import { ref, onMounted, computed, reactive, watch } from 'vue'
import { useRouter } from 'vue-router'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Plus, Trash2, Package, X } from 'lucide-vue-next'
import type {
  ModifierGroupDetail,
  CreateModifierGroupRequest,
  UpdateModifierGroupRequest,
  CreateModifierOptionRequest,
  UpdateModifierOptionRequest,
} from '@/types/modifier.types'
import type { ValidationError } from '@/types/common.types'
import ProductSelectModal from '@/components/common/product/ProductSelectModal.vue'
import { CancelModal } from '@/components/common/cancel'
import { Toggle } from '@/components/common/form'
import { AmountInput } from '@/components/common/form/input/amount'
import {
  Field,
  FieldContent,
  FieldLabel,
  FieldError,
  FieldDescription,
} from '@/components/ui/field'
import { useProductStore } from '@/stores/product.store'
import { useModifierStore } from '@/stores/modifier.store'
import { toast } from 'vue-sonner'
import { useFormErrors } from '@/composables/common/useFormErrors'

const props = defineProps<{
  initialData?: ModifierGroupDetail | null
}>()

const emit = defineEmits<{
  (e: 'success'): void
}>()

const router = useRouter()
const productStore = useProductStore()
const modifierStore = useModifierStore()
const { setErrors, clearErrors, getErrorMessage, hasError } = useFormErrors()

const isEdit = computed(() => !!props.initialData)

const form = reactive({
  name: '',
  isRequired: false,
  options: [] as (CreateModifierOptionRequest | UpdateModifierOptionRequest)[],
  selectedProductIds: [] as string[],
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
      form.isRequired = newData.isRequired
      form.options = newData.options.map((o) => ({
        id: o.id,
        name: o.name,
        priceAdjustment: o.priceAdjustment,
        cogsAdjustment: o.cogsAdjustment,
      }))
      form.selectedProductIds = newData.productModifiers.map((p) => p.id)
    } else {
      form.name = ''
      form.isRequired = false
      form.options = [{ name: '', priceAdjustment: '0', cogsAdjustment: '0' }]
      form.selectedProductIds = []
    }
    initialForm.value = JSON.stringify(form)
  },
  { immediate: true },
)

const isProductModalOpen = ref(false)

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

const selectedProducts = computed(() => {
  return productStore.products.filter((p) => form.selectedProductIds.includes(p.id))
})

function addOption() {
  form.options.push({ name: '', priceAdjustment: '0', cogsAdjustment: '0' })
}

function removeOption(index: number) {
  if (form.options.length > 1) {
    form.options.splice(index, 1)
  } else {
    toast.error('At least one option is required')
  }
}

function handleProductSelect(ids: string[]) {
  form.selectedProductIds = ids
}

function reset() {
  clearErrors()
  if (props.initialData) {
    form.name = props.initialData.name
    form.isRequired = props.initialData.isRequired
    form.options = props.initialData.options.map((o) => ({
      id: o.id,
      name: o.name,
      priceAdjustment: o.priceAdjustment,
      cogsAdjustment: o.cogsAdjustment,
    }))
    form.selectedProductIds = props.initialData.productModifiers.map((p) => p.id)
  } else {
    form.name = ''
    form.isRequired = false
    form.options = [{ name: '', priceAdjustment: '0', cogsAdjustment: '0' }]
    form.selectedProductIds = []
  }
  initialForm.value = JSON.stringify(form)
}

defineExpose({ reset, setErrors, clearErrors })

async function handleSubmit() {
  clearErrors()

  const payload = {
    name: form.name,
    isRequired: form.isRequired,
    options: form.options.map((o) => ({
      ...o,
      priceAdjustment: o.priceAdjustment.toString(),
      cogsAdjustment: o.cogsAdjustment.toString(),
    })),
    productModifiers: form.selectedProductIds.length > 0 ? form.selectedProductIds : null,
  }

  try {
    if (isEdit.value && props.initialData) {
      await modifierStore.update(props.initialData.id, {
        id: props.initialData.id,
        ...payload,
      } as UpdateModifierGroupRequest)
      toast.success('Modifier group updated successfully')
    } else {
      await modifierStore.create(payload as CreateModifierGroupRequest)
      toast.success('Modifier group created successfully')
    }
    emit('success')
  } catch (err: any) {
    if (err?.error && Array.isArray(err.error)) {
      setErrors(err.error as ValidationError[])
    } else {
      toast.error(err?.message || 'Failed to save modifier group')
    }
  }
}

onMounted(async () => {
  await productStore.ensureDataLoaded()
})
</script>

<template>
  <div class="space-y-6 max-w-4xl mx-auto pb-10">
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-semibold">{{ initialData ? 'Edit' : 'Create' }} Modifier Group</h2>
      <div class="flex gap-2">
        <Button variant="outline" @click="handleCancel">
          <X class="size-4 mr-2" />
          Cancel
        </Button>
        <Button :disabled="modifierStore.loading" @click="handleSubmit">
          {{ modifierStore.loading ? 'Saving...' : 'Save Modifier Group' }}
        </Button>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <!-- General Info -->
      <Card class="md:col-span-2">
        <CardHeader>
          <CardTitle>General Information</CardTitle>
        </CardHeader>
        <CardContent class="space-y-4">
          <Field>
            <FieldLabel>Modifier Group Name</FieldLabel>
            <FieldContent>
              <Input
                id="name"
                v-model="form.name"
                placeholder="e.g. Extra Toppings, Sugar Level"
                required
                :aria-invalid="hasError('Name')"
              />
              <FieldError v-if="hasError('Name')" :errors="[getErrorMessage('Name')]" />
            </FieldContent>
          </Field>

          <Field>
            <FieldLabel>Requirement Status</FieldLabel>
            <FieldDescription
              >Specify whether this modifier is required or optional for
              customers.</FieldDescription
            >
            <FieldContent>
              <Toggle
                v-model="form.isRequired"
                :options="[
                  { label: 'Optional', value: false },
                  { label: 'Required', value: true },
                ]"
              />
              <FieldError v-if="hasError('IsRequired')" :errors="[getErrorMessage('IsRequired')]" />
            </FieldContent>
          </Field>
        </CardContent>
      </Card>

      <!-- Assignment -->
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0">
          <CardTitle>Products</CardTitle>
          <Button variant="outline" size="sm" @click="isProductModalOpen = true">
            <Plus class="size-4 mr-1" />
            Assign
          </Button>
        </CardHeader>
        <CardContent>
          <div v-if="selectedProducts.length === 0" class="text-center py-6 text-muted-foreground">
            <Package class="size-8 mx-auto mb-2 opacity-20" />
            <p class="text-xs">No products assigned</p>
          </div>
          <div v-else class="space-y-2 max-h-60 overflow-y-auto pr-1">
            <div
              v-for="product in selectedProducts"
              :key="product.id"
              class="flex items-center justify-between text-sm p-2 bg-accent/50 rounded"
            >
              <span class="truncate pr-2">{{ product.name }}</span>
              <Button
                variant="ghost"
                size="icon"
                class="size-6 text-destructive h-auto w-auto p-1"
                @click="
                  form.selectedProductIds = form.selectedProductIds.filter(
                    (id) => id !== product.id,
                  )
                "
              >
                <Trash2 class="size-3" />
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Options -->
      <Card class="md:col-span-3">
        <CardHeader class="flex flex-row items-center justify-between space-y-0">
          <CardTitle>Modifier Options</CardTitle>
          <Button variant="outline" size="sm" @click="addOption">
            <Plus class="size-4 mr-1" />
            Add Option
          </Button>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <div
              v-for="(option, index) in form.options"
              :key="index"
              class="grid grid-cols-1 md:grid-cols-12 gap-4 items-start border-b pb-4 last:border-0 last:pb-0"
            >
              <Field class="md:col-span-5">
                <FieldLabel :class="{ 'md:hidden': index > 0 }">Name</FieldLabel>
                <FieldContent>
                  <Input
                    v-model="option.name"
                    placeholder="e.g. Extra Cheese"
                    required
                    :aria-invalid="hasError(`Options[${index}].Name`)"
                  />
                  <FieldError
                    v-if="hasError(`Options[${index}].Name`)"
                    :errors="[getErrorMessage(`Options[${index}].Name`)]"
                  />
                </FieldContent>
              </Field>

              <Field class="md:col-span-3">
                <FieldLabel :class="{ 'md:hidden': index > 0 }">Price</FieldLabel>
                <FieldContent>
                  <AmountInput
                    v-model="option.priceAdjustment"
                    mode="money"
                    :aria-invalid="hasError(`Options[${index}].PriceAdjustment`)"
                  />
                  <FieldError
                    v-if="hasError(`Options[${index}].PriceAdjustment`)"
                    :errors="[getErrorMessage(`Options[${index}].PriceAdjustment`)]"
                  />
                </FieldContent>
              </Field>

              <Field class="md:col-span-3">
                <FieldLabel :class="{ 'md:hidden': index > 0 }">COGS</FieldLabel>
                <FieldContent>
                  <AmountInput
                    v-model="option.cogsAdjustment"
                    mode="money"
                    :aria-invalid="hasError(`Options[${index}].CogsAdjustment`)"
                  />
                  <FieldError
                    v-if="hasError(`Options[${index}].CogsAdjustment`)"
                    :errors="[getErrorMessage(`Options[${index}].CogsAdjustment`)]"
                  />
                </FieldContent>
              </Field>

              <div class="md:col-span-1 flex justify-end pt-2" :class="{ 'md:pt-8': index === 0 }">
                <Button
                  variant="ghost"
                  size="icon"
                  class="text-destructive"
                  :disabled="form.options.length === 1"
                  @click="removeOption(index)"
                >
                  <Trash2 class="size-4" />
                </Button>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>

    <ProductSelectModal
      v-model:open="isProductModalOpen"
      :selected-ids="form.selectedProductIds"
      @select="handleProductSelect"
    />

    <CancelModal v-model:open="showCancelModal" @confirm="confirmCancel" />
  </div>
</template>
