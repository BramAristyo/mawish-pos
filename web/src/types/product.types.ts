import type { BaseResponse } from './common.types'
import type { Category } from './category.types'
import type { ModifierGroup } from './modifier.types'

export interface Product {
  id: string
  name: string
  description?: string
  price: number
  cogs: number
  categoryId: string
  category: Category
  modifierGroups?: ModifierGroup[]
  imageUrl?: string
  createdAt: string
  updatedAt: string
  deletedAt?: string
}

export interface StoreProductRequest {
  name: string
  description?: string
  price: number
  cogs: number
  categoryId: string
  modifierGroupIds?: string[]
  imageUrl?: string
}

export interface UpdateProductRequest {
  name: string
  description?: string
  price: number
  cogs: number
  categoryId: string
  modifierGroupIds?: string[]
  imageUrl?: string
}

export type PaginatedProductResponse = BaseResponse<Product[]>
