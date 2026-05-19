<script setup lang="ts">
import { reactive, watch, onMounted, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Field, FieldContent, FieldLabel, FieldError } from '@/components/ui/field'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { ScrollArea } from '@/components/ui/scroll-area'
import { AmountInput } from '@/components/common/form/input/amount'
import { Checkbox } from '@/components/ui/checkbox'
import { useCategoryStore } from '@/stores/category.stores'
import { useModifierStore } from '@/stores/modifier.store'
import { useFormErrors } from '@/composables/common/useFormErrors'
import { CancelModal } from '@/components/common/cancel'
import type { Product, StoreProductRequest } from '@/types/product.types'

const props = defineProps<{
  initialData?: Product | null
  loading?: boolean
}>()

const emit = defineEmits<{
  (e: 'submit', data: StoreProductRequest): void
}>()

const router = useRouter()
const categoryStore = useCategoryStore()
const modifierStore = useModifierStore()
const { getErrorMessage, hasError, setErrors, clearErrors } = useFormErrors()

const form = reactive<{
  name: string
  description: string
  price: string
  cogs: string
  categoryId: string
  modifierGroupIds: string[]
}>({
  name: '',
  description: '',
  price: '0',
  cogs: '0',
  categoryId: '',
  modifierGroupIds: [],
})

const initialForm = ref<string>('')
const showCancelModal = ref(false)
const isDirty = computed(() => JSON.stringify(form) !== initialForm.value)

watch(
  () => props.initialData,
  (newData) => {
    if (newData) {
      form.name = newData.name
      form.description = newData.description || ''
      form.price = newData.price.toString()
      form.cogs = newData.cogs.toString()
      form.categoryId = newData.categoryId
      form.modifierGroupIds = newData.modifierGroups?.map((mg) => mg.id) || []
    } else {
      form.name = ''
      form.description = ''
      form.price = '0'
      form.cogs = '0'
      form.categoryId = ''
      form.modifierGroupIds = []
    }
    initialForm.value = JSON.stringify(form)
  },
  { immediate: true },
)

onMounted(async () => {
  await Promise.all([
    categoryStore.categories.length === 0 ? categoryStore.fetchAll() : Promise.resolve(),
    modifierStore.modifiers.length === 0 ? modifierStore.fetchAll() : Promise.resolve(),
  ])
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

function handleModifierToggle(id: string, checked: boolean) {
  if (checked) {
    if (!form.modifierGroupIds.includes(id)) {
      form.modifierGroupIds.push(id)
    }
  } else {
    form.modifierGroupIds = form.modifierGroupIds.filter((mgId) => mgId !== id)
  }
}

const handleSubmit = () => {
  emit('submit', {
    name: form.name,
    description: form.description,
    price: Number(form.price),
    cogs: Number(form.cogs),
    categoryId: form.categoryId,
    modifierGroupIds: form.modifierGroupIds,
  })
}

defineExpose({
  setErrors,
  clearErrors,
})
</script>

<template>
  <form @submit.prevent="handleSubmit" class="space-y-8 pb-10">
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Main Content (Left 2/3) -->
      <div class="lg:col-span-2 space-y-6">
        <Card>
          <CardHeader>
            <CardTitle>General Information</CardTitle>
            <CardDescription>
              Basic details about the product that will be visible to customers.
            </CardDescription>
          </CardHeader>
          <CardContent class="space-y-6">
            <Field>
              <FieldLabel>Product Name</FieldLabel>
              <FieldContent>
                <Input
                  v-model="form.name"
                  placeholder="Enter product name (e.g. Classic Latte)"
                  required
                  :aria-invalid="hasError('Name')"
                />
                <FieldError v-if="hasError('Name')" :errors="[getErrorMessage('Name')]" />
              </FieldContent>
            </Field>

            <Field>
              <FieldLabel>Description</FieldLabel>
              <FieldContent>
                <textarea
                  v-model="form.description"
                  placeholder="Tell us more about this product..."
                  class="flex min-h-[120px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive transition-all"
                  :aria-invalid="hasError('Description')"
                ></textarea>
                <FieldError
                  v-if="hasError('Description')"
                  :errors="[getErrorMessage('Description')]"
                />
              </FieldContent>
            </Field>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Configuration</CardTitle>
            <CardDescription>
              Set category and available customizations for this product.
            </CardDescription>
          </CardHeader>
          <CardContent class="space-y-6">
            <Field>
              <FieldLabel>Category</FieldLabel>
              <FieldContent>
                <select
                  v-model="form.categoryId"
                  class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive"
                  required
                  :aria-invalid="hasError('CategoryId')"
                >
                  <option value="" disabled>Select a category</option>
                  <option
                    v-for="category in categoryStore.categories"
                    :key="category.id"
                    :value="category.id"
                  >
                    {{ category.name }}
                  </option>
                </select>
                <FieldError
                  v-if="hasError('CategoryId')"
                  :errors="[getErrorMessage('CategoryId')]"
                />
              </FieldContent>
            </Field>

            <Field>
              <FieldLabel>Modifier Groups</FieldLabel>
              <div class="rounded-lg border bg-muted/30">
                <ScrollArea class="h-[200px]">
                  <div class="p-4 grid grid-cols-1 sm:grid-cols-2 gap-3">
                    <label
                      v-for="mg in modifierStore.modifiers"
                      :key="mg.id"
                      class="flex items-center space-x-3 p-3 rounded-md border bg-card hover:bg-accent/50 transition-colors cursor-pointer group"
                    >
                      <Checkbox
                        :id="mg.id"
                        :checked="form.modifierGroupIds.includes(mg.id)"
                        @update:checked="(checked: boolean | 'indeterminate') => handleModifierToggle(mg.id, checked === true)"
                      />
                      <span
                        class="text-sm font-medium leading-none flex-1 group-hover:text-accent-foreground"
                      >
                        {{ mg.name }}
                      </span>
                    </label>
                  </div>
                  <div
                    v-if="modifierStore.modifiers.length === 0"
                    class="flex flex-col items-center justify-center h-[160px] text-muted-foreground italic"
                  >
                    <p class="text-sm">No modifier groups available.</p>
                  </div>
                </ScrollArea>
              </div>
              <FieldError
                v-if="hasError('ModifierGroupIds')"
                :errors="[getErrorMessage('ModifierGroupIds')]"
              />
            </Field>
          </CardContent>
        </Card>
      </div>

      <!-- Sidebar (Right 1/3) -->
      <div class="space-y-6">
        <Card>
          <CardHeader>
            <CardTitle>Pricing</CardTitle>
            <CardDescription> Manage your product's financial details. </CardDescription>
          </CardHeader>
          <CardContent class="space-y-6">
            <Field>
              <FieldLabel>Selling Price</FieldLabel>
              <FieldContent>
                <AmountInput
                  v-model="form.price"
                  mode="money"
                  placeholder="0.00"
                  :aria-invalid="hasError('Price')"
                />
                <FieldError v-if="hasError('Price')" :errors="[getErrorMessage('Price')]" />
              </FieldContent>
            </Field>

            <Field>
              <FieldLabel>Cost (COGS)</FieldLabel>
              <FieldContent>
                <AmountInput
                  v-model="form.cogs"
                  mode="money"
                  placeholder="0.00"
                  :aria-invalid="hasError('Cogs')"
                />
                <p class="text-[0.8rem] text-muted-foreground mt-1.5">
                  Cost of Goods Sold (used for profit margin calculation).
                </p>
                <FieldError v-if="hasError('Cogs')" :errors="[getErrorMessage('Cogs')]" />
              </FieldContent>
            </Field>
          </CardContent>
        </Card>
      </div>
    </div>

    <!-- Actions -->
    <div class="flex items-center justify-end gap-3 border-t pt-6">
      <Button variant="ghost" type="button" :disabled="loading" @click="handleCancel">
        Cancel
      </Button>
      <Button type="submit" :disabled="loading" class="min-w-[140px]">
        {{ loading ? 'Saving...' : 'Save Product' }}
      </Button>
    </div>
    <CancelModal v-model:open="showCancelModal" @confirm="confirmCancel" />
  </form>
</template>
