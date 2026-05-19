import { useAuthStore } from '@/stores/auth.stores'
import { createRouter, createWebHistory } from 'vue-router'

declare module 'vue-router' {
  interface RouteMeta {
    requiresAuth?: boolean
    title?: string
  }
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'login',
      component: () => import('@/views/auth/LoginView.vue'),
      meta: { title: 'Login' }
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('@/views/dashboard/DashboardView.vue'),
      meta: { requiresAuth: true, title: 'Dashboard' },
    },
    {
      path: '/categories',
      name: 'categories',
      component: () => import('@/views/category/CategoryView.vue'),
      meta: { requiresAuth: true, title: 'Categories' },
    },
    {
      path: '/catalog/products',
      name: 'products',
      component: () => import('@/views/catalog/product/ProductView.vue'),
      meta: { requiresAuth: true, title: 'Products' },
    },
    {
      path: '/catalog/products/create',
      name: 'product-create',
      component: () => import('@/views/catalog/product/ProductCreateView.vue'),
      meta: { requiresAuth: true, title: 'Create Product' },
    },
    {
      path: '/catalog/products/:id/edit',
      name: 'product-edit',
      component: () => import('@/views/catalog/product/ProductEditView.vue'),
      meta: { requiresAuth: true, title: 'Edit Product' },
    },
    {
      path: '/modifiers',
      name: 'modifiers',
      component: () => import('@/views/modifier/ModifierView.vue'),
      meta: { requiresAuth: true, title: 'Modifier Groups' },
    },
    {
      path: '/modifiers/create',
      name: 'modifier-create',
      component: () => import('@/views/modifier/ModifierCreateView.vue'),
      meta: { requiresAuth: true, title: 'Create Modifier Group' },
    },
    {
      path: '/modifiers/:id/edit',
      name: 'modifier-edit',
      component: () => import('@/views/modifier/ModifierEditView.vue'),
      meta: { requiresAuth: true, title: 'Edit Modifier Group' },
    },
    {
      path: '/taxes',
      name: 'taxes',
      component: () => import('@/views/tax/TaxView.vue'),
      meta: { requiresAuth: true, title: 'Taxes' },
    },
    {
      path: '/discounts',
      name: 'discounts',
      component: () => import('@/views/discount/DiscountView.vue'),
      meta: { requiresAuth: true, title: 'Discounts' },
    },
    {
      path: '/discounts/create',
      name: 'discount-create',
      component: () => import('@/views/discount/DiscountCreateView.vue'),
      meta: { requiresAuth: true, title: 'Create Discount' },
    },
    {
      path: '/discounts/:id/edit',
      name: 'discount-edit',
      component: () => import('@/views/discount/DiscountEditView.vue'),
      meta: { requiresAuth: true, title: 'Edit Discount' },
    },
    {
      path: '/sales-types',
      name: 'sales-types',
      component: () => import('@/views/sales-type/SalesTypeView.vue'),
      meta: { requiresAuth: true, title: 'Sales Types' },
    },
    {
      path: '/employees',
      name: 'employees',
      component: () => import('@/views/employee/EmployeeView.vue'),
      meta: { requiresAuth: true, title: 'Employees' },
    },
    {
      path: '/shift-schedules',
      name: 'shift-schedules',
      component: () => import('@/views/shift-schedule/ShiftScheduleView.vue'),
      meta: { requiresAuth: true, title: 'Shift Schedules' },
    },
    {
      path: '/attendances',
      name: 'attendances',
      component: () => import('@/views/attendance/AttendanceView.vue'),
      meta: { requiresAuth: true, title: 'Attendance' },
    },
    {
      path: '/coa',
      name: 'coa',
      component: () => import('@/views/coa/COAView.vue'),
      meta: { requiresAuth: true, title: 'Chart of Accounts' },
    },
    {
      path: '/payroll',
      name: 'payroll',
      component: () => import('@/views/payroll/PayrollView.vue'),
      meta: { requiresAuth: true, title: 'Payroll' },
    },
    {
      path: '/accounting/cash-transactions',
      name: 'cash-transactions',
      component: () => import('@/views/accounting/cash-transaction/CashTransactionView.vue'),
      meta: { requiresAuth: true, title: 'Cash Transactions' },
    },
    {
      path: '/accounting/cash-transactions/create',
      name: 'cash-transaction-create',
      component: () => import('@/views/accounting/cash-transaction/CashTransactionCreateView.vue'),
      meta: { requiresAuth: true, title: 'Create Cash Transaction' },
    },
    {
      path: '/accounting/cash-transactions/:id/edit',
      name: 'cash-transaction-edit',
      component: () => import('@/views/accounting/cash-transaction/CashTransactionEditView.vue'),
      meta: { requiresAuth: true, title: 'Edit Cash Transaction' },
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

router.beforeEach((to) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return '/'
  } else if (to.path == '/' && authStore.isAuthenticated) {
    return '/dashboard'
  }
})

export default router
