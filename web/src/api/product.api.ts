import http from '@/lib/http'
import type { BaseFilterRequest, BaseResponse } from '@/types/common.types'
import type {
  Product,
  StoreProductRequest,
  UpdateProductRequest,
} from '@/types/product.types'

export const productApi = {
  paginate: (params: BaseFilterRequest) =>
    http.get<any, BaseResponse<Product[]>>('/products', { params }),

  getAll: () => http.get<any, BaseResponse<Product[]>>('/products/get-all'),

  findById: (id: string) => http.get<any, BaseResponse<Product>>(`/products/${id}`),

  store: (payload: StoreProductRequest) =>
    http.post<any, BaseResponse<Product>>('/products', payload),

  update: (id: string, payload: UpdateProductRequest) =>
    http.put<any, BaseResponse<Product>>(`/products/${id}`, payload),

  delete: (id: string) => http.delete<any, BaseResponse<null>>(`/products/${id}`),

  restore: (id: string) => http.patch<any, BaseResponse<Product>>(`/products/${id}/restore`),
}
